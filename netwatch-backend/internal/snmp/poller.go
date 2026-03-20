package snmp

import (
	"context"
	"sync"
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PollResult é o resultado de uma coleta de um dispositivo
type PollResult struct {
	DeviceID    uuid.UUID
	Metrics     []models.Metric
	Interfaces  []InterfaceData
	SysInfo     *SysInfo
	Success     bool
	Error       error
	CollectedAt time.Time
}

// Poller gerencia o ciclo de polling de um dispositivo
type Poller struct {
	device     *models.Device
	metricRepo repository.MetricRepository
	deviceRepo repository.DeviceRepository
	timeout    time.Duration
	retries    int
	log        *zap.Logger
	resultChan chan<- PollResult
}

// NewPoller cria um Poller para um dispositivo
func NewPoller(
	device *models.Device,
	metricRepo repository.MetricRepository,
	deviceRepo repository.DeviceRepository,
	timeout time.Duration,
	retries int,
	resultChan chan<- PollResult,
	log *zap.Logger,
) *Poller {
	return &Poller{
		device:     device,
		metricRepo: metricRepo,
		deviceRepo: deviceRepo,
		timeout:    timeout,
		retries:    retries,
		log:        log,
		resultChan: resultChan,
	}
}

// Poll executa uma única coleta SNMP do dispositivo
func (p *Poller) Poll(ctx context.Context) {
	result := PollResult{
		DeviceID:    p.device.ID,
		CollectedAt: time.Now().UTC(),
	}

	client, err := NewClient(p.device, p.timeout, p.retries)
	if err != nil {
		result.Error = err
		p.handleFailure(result)
		return
	}

	if err := client.Connect(); err != nil {
		result.Error = err
		p.handleFailure(result)
		return
	}
	defer client.Close()

	// Coletar info do sistema
	sysInfo, err := client.GetSysInfo()
	if err != nil {
		result.Error = err
		p.handleFailure(result)
		return
	}
	result.SysInfo = sysInfo

	// Coletar interfaces
	ifaces, err := client.GetInterfaces()
	if err != nil {
		p.log.Warn("Erro ao coletar interfaces", zap.String("device", p.device.IPAddress), zap.Error(err))
	} else {
		result.Interfaces = ifaces
	}

	// Montar métricas de interfaces
	now := result.CollectedAt
	for _, iface := range ifaces {
		result.Metrics = append(result.Metrics,
			models.Metric{
				DeviceID:       p.device.ID,
				Type:           models.MetricIfInOctets,
				InterfaceIndex: iface.Index,
				Value:          float64(iface.InOctets),
				Unit:           "bytes",
				CollectedAt:    now,
			},
			models.Metric{
				DeviceID:       p.device.ID,
				Type:           models.MetricIfOutOctets,
				InterfaceIndex: iface.Index,
				Value:          float64(iface.OutOctets),
				Unit:           "bytes",
				CollectedAt:    now,
			},
			models.Metric{
				DeviceID:       p.device.ID,
				Type:           models.MetricIfInErrors,
				InterfaceIndex: iface.Index,
				Value:          float64(iface.InErrors),
				Unit:           "errors",
				CollectedAt:    now,
			},
			models.Metric{
				DeviceID:       p.device.ID,
				Type:           models.MetricIfOutErrors,
				InterfaceIndex: iface.Index,
				Value:          float64(iface.OutErrors),
				Unit:           "errors",
				CollectedAt:    now,
			},
		)
	}

	// Coletar métricas específicas por vendor
	p.collectVendorMetrics(client, sysInfo.Vendor, &result, now)

	// Uptime
	result.Metrics = append(result.Metrics, models.Metric{
		DeviceID:    p.device.ID,
		Type:        models.MetricUptime,
		OID:         OIDSysUptime,
		Value:       sysInfo.Uptime / 100, // centissegundos → segundos
		Unit:        "seconds",
		CollectedAt: now,
	})

	result.Success = true
	p.handleSuccess(result)
}

func (p *Poller) collectVendorMetrics(client *Client, vendor string, result *PollResult, now time.Time) {
	oids, ok := MetricOIDsByVendor[vendor]
	if !ok {
		oids = MetricOIDsByVendor["generic"]
	}

	for metricName, oid := range oids {
		val, err := client.GetSingle(oid)
		if err != nil {
			p.log.Debug("OID não disponível",
				zap.String("device", p.device.IPAddress),
				zap.String("metric", metricName),
				zap.String("oid", oid),
			)
			continue
		}

		metricType, unit := metricNameToType(metricName, vendor)
		result.Metrics = append(result.Metrics, models.Metric{
			DeviceID:    p.device.ID,
			Type:        metricType,
			OID:         oid,
			Value:       val,
			Unit:        unit,
			CollectedAt: now,
		})
	}
}

func (p *Poller) handleSuccess(result PollResult) {
	// Salvar métricas no banco
	if len(result.Metrics) > 0 {
		if err := p.metricRepo.SaveBatch(result.Metrics); err != nil {
			p.log.Error("Erro ao salvar métricas", zap.Error(err), zap.String("device", p.device.IPAddress))
		}
	}

	// Atualizar status do device para online
	_ = p.deviceRepo.UpdateStatus(p.device.ID, models.DeviceStatusOnline)
	_ = p.deviceRepo.UpdateLastSeen(p.device.ID)

	p.log.Debug("Poll concluído com sucesso",
		zap.String("device", p.device.IPAddress),
		zap.Int("metrics", len(result.Metrics)),
	)

	if p.resultChan != nil {
		select {
		case p.resultChan <- result:
		default:
		}
	}
}

func (p *Poller) handleFailure(result PollResult) {
	p.log.Warn("Falha no poll do dispositivo",
		zap.String("device", p.device.IPAddress),
		zap.Error(result.Error),
	)

	_ = p.deviceRepo.UpdateStatus(p.device.ID, models.DeviceStatusOffline)

	if p.resultChan != nil {
		select {
		case p.resultChan <- result:
		default:
		}
	}
}

// PollerScheduler gerencia os timers de polling de todos os dispositivos
type PollerScheduler struct {
	metricRepo repository.MetricRepository
	deviceRepo repository.DeviceRepository
	timeout    time.Duration
	retries    int
	log        *zap.Logger
	resultChan chan PollResult
	timers     map[uuid.UUID]*time.Ticker
	cancel     map[uuid.UUID]context.CancelFunc
	mu         sync.Mutex
}

func NewPollerScheduler(
	metricRepo repository.MetricRepository,
	deviceRepo repository.DeviceRepository,
	timeout time.Duration,
	retries int,
	resultChan chan PollResult,
	log *zap.Logger,
) *PollerScheduler {
	return &PollerScheduler{
		metricRepo: metricRepo,
		deviceRepo: deviceRepo,
		timeout:    timeout,
		retries:    retries,
		log:        log,
		resultChan: resultChan,
		timers:     make(map[uuid.UUID]*time.Ticker),
		cancel:     make(map[uuid.UUID]context.CancelFunc),
	}
}

// AddDevice registra um dispositivo para polling
func (s *PollerScheduler) AddDevice(device *models.Device) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.timers[device.ID]; exists {
		s.removeDeviceLocked(device.ID)
	}

	if !device.Enabled {
		return
	}

	interval := time.Duration(device.PollingInterval) * time.Second
	ticker := time.NewTicker(interval)
	ctx, cancel := context.WithCancel(context.Background())

	s.timers[device.ID] = ticker
	s.cancel[device.ID] = cancel

	deviceCopy := *device

	go func() {
		// Poll imediato na adição
		p := NewPoller(&deviceCopy, s.metricRepo, s.deviceRepo, s.timeout, s.retries, s.resultChan, s.log)
		p.Poll(ctx)

		for {
			select {
			case <-ticker.C:
				p.Poll(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	s.log.Info("Dispositivo adicionado ao scheduler",
		zap.String("device", device.IPAddress),
		zap.Duration("interval", interval),
	)
}

// RemoveDevice para o polling de um dispositivo
func (s *PollerScheduler) RemoveDevice(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeDeviceLocked(id)
}

func (s *PollerScheduler) removeDeviceLocked(id uuid.UUID) {
	if t, ok := s.timers[id]; ok {
		t.Stop()
		delete(s.timers, id)
	}
	if c, ok := s.cancel[id]; ok {
		c()
		delete(s.cancel, id)
	}
}

// StopAll encerra todos os pollers
func (s *PollerScheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id := range s.timers {
		s.removeDeviceLocked(id)
	}
	s.log.Info("Todos os pollers encerrados")
}

// metricNameToType converte nome de métrica interna para MetricType e unidade
func metricNameToType(name, vendor string) (models.MetricType, string) {
	switch name {
	case "cpu", "cpu_user":
		return models.MetricCPU, "%"
	case "memory_used", "mem_used":
		return models.MetricMemory, "bytes"
	case "memory_total", "mem_total":
		return models.MetricMemory, "bytes"
	case "temperature":
		return models.MetricTemperature, "°C"
	case "cpu_idle":
		return models.MetricCPU, "%"
	default:
		return models.MetricCustom, ""
	}
}
