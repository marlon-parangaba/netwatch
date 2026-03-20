package services

import (
	"context"

	"github.com/ertel/netwatch-backend/internal/config"
	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	snmpengine "github.com/ertel/netwatch-backend/internal/snmp"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PollerService orquestra o scheduler de polling e o receptor de traps
type PollerService struct {
	scheduler  *snmpengine.PollerScheduler
	trapRcv    *snmpengine.TrapReceiver
	deviceRepo repository.DeviceRepository
	metricRepo repository.MetricRepository
	cfg        *config.Config
	log        *zap.Logger
	resultChan chan snmpengine.PollResult
	onTrap     func(snmpengine.TrapEvent)
}

// NewPollerService cria o PollerService
func NewPollerService(
	deviceRepo repository.DeviceRepository,
	metricRepo repository.MetricRepository,
	cfg *config.Config,
	log *zap.Logger,
	onTrap func(snmpengine.TrapEvent),
) *PollerService {
	resultChan := make(chan snmpengine.PollResult, 500)

	scheduler := snmpengine.NewPollerScheduler(
		metricRepo,
		deviceRepo,
		cfg.SNMP.Timeout,
		cfg.SNMP.Retries,
		resultChan,
		log,
	)

	trapRcv := snmpengine.NewTrapReceiver(
		cfg.SNMP.TrapPort,
		cfg.SNMP.Community,
		deviceRepo,
		metricRepo,
		log,
		onTrap,
	)

	return &PollerService{
		scheduler:  scheduler,
		trapRcv:    trapRcv,
		deviceRepo: deviceRepo,
		metricRepo: metricRepo,
		cfg:        cfg,
		log:        log,
		resultChan: resultChan,
		onTrap:     onTrap,
	}
}

// Start inicializa o serviço: carrega devices do banco e inicia polling + traps
func (s *PollerService) Start(ctx context.Context) error {
	// Carregar todos os devices habilitados
	devices, _, err := s.deviceRepo.FindAll(1, 10000, repository.DeviceFilters{})
	if err != nil {
		return err
	}

	enabled := true
	enabledDevices, _, err := s.deviceRepo.FindAll(1, 10000, repository.DeviceFilters{Enabled: &enabled})
	if err != nil {
		return err
	}
	_ = devices

	for i := range enabledDevices {
		s.scheduler.AddDevice(&enabledDevices[i])
	}

	s.log.Info("PollerService iniciado",
		zap.Int("devices", len(enabledDevices)),
	)

	// Iniciar receptor de traps em goroutine separada
	go func() {
		if err := s.trapRcv.Start(); err != nil {
			s.log.Error("Erro no receptor de traps", zap.Error(err))
		}
	}()

	// Consumir resultados em background (para WebSocket / alertas futuros)
	go s.consumeResults(ctx)

	return nil
}

// Stop encerra o polling e o receptor de traps
func (s *PollerService) Stop() {
	s.scheduler.StopAll()
	s.trapRcv.Stop()
	s.log.Info("PollerService encerrado")
}

// AddDevice registra um novo dispositivo no scheduler em tempo real
func (s *PollerService) AddDevice(device *models.Device) {
	s.scheduler.AddDevice(device)
}

// RemoveDevice remove um dispositivo do scheduler
func (s *PollerService) RemoveDevice(id uuid.UUID) {
	s.scheduler.RemoveDevice(id)
}

// ResultChan expõe o canal de resultados para outros consumidores (ex: WebSocket)
func (s *PollerService) ResultChan() <-chan snmpengine.PollResult {
	return s.resultChan
}

// consumeResults processa resultados em background (extensível para alertas)
func (s *PollerService) consumeResults(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case result, ok := <-s.resultChan:
			if !ok {
				return
			}
			if !result.Success {
				s.log.Debug("Poll falhou",
					zap.String("device_id", result.DeviceID.String()),
					zap.Error(result.Error),
				)
			}
		}
	}
}
