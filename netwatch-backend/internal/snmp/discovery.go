package snmp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"go.uber.org/zap"
)

// DiscoveryResult representa um dispositivo encontrado durante a varredura
type DiscoveryResult struct {
	IPAddress string
	SysInfo   *SysInfo
	Vendor    string
	Type      models.DeviceType
}

// DiscoveryConfig parâmetros para a varredura de rede
type DiscoveryConfig struct {
	Network       string            // CIDR, ex: "192.168.1.0/24"
	SNMPVersion   models.SNMPVersion
	Community     string
	Port          int
	Timeout       time.Duration
	Retries       int
	Concurrency   int               // goroutines paralelas (default: 50)
	ProgressChan  chan<- DiscoveryProgress
}

// DiscoveryProgress é enviado no canal de progresso a cada host testado
type DiscoveryProgress struct {
	IP        string
	Found     bool
	Error     string
	Processed int
	Total     int
}

// Discover varre uma rede CIDR e retorna dispositivos respondentes ao SNMP
func Discover(ctx context.Context, cfg DiscoveryConfig, log *zap.Logger) ([]DiscoveryResult, error) {
	if cfg.Port == 0 {
		cfg.Port = 161
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 3 * time.Second
	}
	if cfg.Retries == 0 {
		cfg.Retries = 1
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 50
	}
	if cfg.Community == "" {
		cfg.Community = "public"
	}
	if cfg.SNMPVersion == "" {
		cfg.SNMPVersion = models.SNMPv2c
	}

	hosts, err := expandCIDR(cfg.Network)
	if err != nil {
		return nil, fmt.Errorf("CIDR inválido %q: %w", cfg.Network, err)
	}

	log.Info("Iniciando discovery de rede",
		zap.String("network", cfg.Network),
		zap.Int("hosts", len(hosts)),
		zap.Int("concurrency", cfg.Concurrency),
	)

	type job struct {
		ip string
	}

	jobs := make(chan job, len(hosts))
	resultsCh := make(chan DiscoveryResult, len(hosts))

	var wg sync.WaitGroup
	processed := 0
	var mu sync.Mutex

	// Worker pool
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				result, err := probeHost(j.ip, cfg)

				mu.Lock()
				processed++
				current := processed
				mu.Unlock()

				if cfg.ProgressChan != nil {
					prog := DiscoveryProgress{
						IP:        j.ip,
						Found:     err == nil,
						Processed: current,
						Total:     len(hosts),
					}
					if err != nil {
						prog.Error = err.Error()
					}
					select {
					case cfg.ProgressChan <- prog:
					default:
					}
				}

				if err == nil {
					resultsCh <- *result
					log.Info("Dispositivo encontrado",
						zap.String("ip", j.ip),
						zap.String("name", result.SysInfo.Name),
						zap.String("vendor", result.Vendor),
					)
				}
			}
		}()
	}

	// Enfileirar jobs
	for _, ip := range hosts {
		jobs <- job{ip: ip}
	}
	close(jobs)

	// Aguardar todos os workers e fechar resultados
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var results []DiscoveryResult
	for r := range resultsCh {
		results = append(results, r)
	}

	log.Info("Discovery concluído",
		zap.String("network", cfg.Network),
		zap.Int("found", len(results)),
		zap.Int("scanned", len(hosts)),
	)

	return results, nil
}

// probeHost testa um único host via SNMP
func probeHost(ip string, cfg DiscoveryConfig) (*DiscoveryResult, error) {
	// Cria um Device temporário para usar o client
	device := &models.Device{
		IPAddress:     ip,
		SNMPVersion:   cfg.SNMPVersion,
		SNMPCommunity: cfg.Community,
		SNMPPort:      cfg.Port,
	}

	client, err := NewClient(device, cfg.Timeout, cfg.Retries)
	if err != nil {
		return nil, err
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}
	defer client.Close()

	sysInfo, err := client.GetSysInfo()
	if err != nil {
		return nil, err
	}

	return &DiscoveryResult{
		IPAddress: ip,
		SysInfo:   sysInfo,
		Vendor:    sysInfo.Vendor,
		Type:      vendorToDeviceType(sysInfo.Vendor),
	}, nil
}

// expandCIDR retorna todos os IPs utilizáveis de um bloco CIDR
func expandCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		// Pular endereço de rede e broadcast
		if ip[len(ip)-1] == 0 || ip[len(ip)-1] == 255 {
			continue
		}
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func vendorToDeviceType(vendor string) models.DeviceType {
	switch vendor {
	case "mikrotik":
		return models.DeviceTypeMikrotik
	case "cisco":
		return models.DeviceTypeCisco
	case "huawei":
		return models.DeviceTypeHuawei
	case "juniper":
		return models.DeviceTypeJuniper
	case "ubiquiti":
		return models.DeviceTypeUbiquiti
	default:
		return models.DeviceTypeGeneric
	}
}

// ToDevice converte um DiscoveryResult em um Device pronto para salvar
func (r *DiscoveryResult) ToDevice(community string, version models.SNMPVersion, port int) models.Device {
	name := r.SysInfo.Name
	if name == "" {
		name = r.IPAddress
	}

	return models.Device{
		Name:          name,
		IPAddress:     r.IPAddress,
		Hostname:      r.SysInfo.Name,
		Type:          r.Type,
		SNMPVersion:   version,
		SNMPCommunity: community,
		SNMPPort:      port,
		Status:        models.DeviceStatusOnline,
		SysName:       r.SysInfo.Name,
		SysDescr:      r.SysInfo.Descr,
		SysOID:        r.SysInfo.OID,
		SysContact:    r.SysInfo.Contact,
		Enabled:       true,
	}
}
