package snmp

import (
	"fmt"
	"net"
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/gosnmp/gosnmp"
	"go.uber.org/zap"
)

// TrapEvent representa um SNMP Trap recebido e processado
type TrapEvent struct {
	SourceIP    string
	Community   string
	Version     string
	Enterprise  string
	GenericTrap int
	SpecificTrap int
	Variables   []TrapVariable
	ReceivedAt  time.Time
}

// TrapVariable é um binding OID→valor de um Trap
type TrapVariable struct {
	OID   string
	Value string
	Type  string
}

// TrapReceiver escuta traps UDP e os processa
type TrapReceiver struct {
	port       int
	community  string
	deviceRepo repository.DeviceRepository
	metricRepo repository.MetricRepository
	log        *zap.Logger
	listener   *gosnmp.TrapListener
	onTrap     func(TrapEvent)
}

// NewTrapReceiver cria um receptor de SNMP Traps
func NewTrapReceiver(
	port int,
	community string,
	deviceRepo repository.DeviceRepository,
	metricRepo repository.MetricRepository,
	log *zap.Logger,
	onTrap func(TrapEvent),
) *TrapReceiver {
	return &TrapReceiver{
		port:       port,
		community:  community,
		deviceRepo: deviceRepo,
		metricRepo: metricRepo,
		log:        log,
		onTrap:     onTrap,
	}
}

// Start inicia o listener UDP para traps (bloqueante — usar em goroutine)
func (r *TrapReceiver) Start() error {
	r.listener = gosnmp.NewTrapListener()
	r.listener.OnNewTrap = r.handleTrap
	r.listener.Params = gosnmp.Default

	addr := fmt.Sprintf("0.0.0.0:%d", r.port)
	r.log.Info("Receptor de SNMP Traps iniciado", zap.String("addr", addr))

	return r.listener.Listen(addr)
}

// Stop encerra o listener
func (r *TrapReceiver) Stop() {
	if r.listener != nil {
		r.listener.Close()
		r.log.Info("Receptor de SNMP Traps encerrado")
	}
}

// handleTrap é chamado pelo gosnmp para cada trap recebido
func (r *TrapReceiver) handleTrap(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	sourceIP := addr.IP.String()

	event := TrapEvent{
		SourceIP:   sourceIP,
		ReceivedAt: time.Now().UTC(),
	}

	switch packet.Version {
	case gosnmp.Version1:
		event.Version = "v1"
		event.Community = packet.Community
		if trap, ok := packet.PDU.(gosnmp.Trap); ok {
			event.Enterprise = trap.Enterprise
			event.GenericTrap = trap.GenericTrap
			event.SpecificTrap = trap.SpecificTrap
		}
	case gosnmp.Version2c:
		event.Version = "v2c"
		event.Community = packet.Community
	case gosnmp.Version3:
		event.Version = "v3"
	}

	for _, v := range packet.Variables {
		event.Variables = append(event.Variables, TrapVariable{
			OID:   v.Name,
			Value: toString(v.Value),
			Type:  v.Type.String(),
		})
	}

	r.log.Info("Trap recebido",
		zap.String("source", sourceIP),
		zap.String("version", event.Version),
		zap.Int("variables", len(event.Variables)),
	)

	// Processar o trap: atualizar status do device, gerar alertas, etc.
	r.processTrap(event)

	// Notificar callback externo (ex: WebSocket)
	if r.onTrap != nil {
		r.onTrap(event)
	}
}

// processTrap trata o trap e atualiza o banco de dados conforme necessário
func (r *TrapReceiver) processTrap(event TrapEvent) {
	// Tentar identificar o dispositivo pelo IP de origem
	device, err := r.deviceRepo.FindByIP(event.SourceIP)
	if err != nil {
		r.log.Debug("Trap recebido de IP desconhecido", zap.String("ip", event.SourceIP))
		return
	}

	// Verificar tipo de trap e tomar ação correspondente
	r.classifyAndAct(device, event)
}

// classifyAndAct identifica o tipo de trap e age de acordo
func (r *TrapReceiver) classifyAndAct(device *models.Device, event TrapEvent) {
	// coldStart (0) / warmStart (1) → dispositivo reiniciou
	if event.GenericTrap == 0 || event.GenericTrap == 1 {
		r.log.Info("Dispositivo reiniciou (coldStart/warmStart)",
			zap.String("device", device.IPAddress),
			zap.Int("trap", event.GenericTrap),
		)
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusOnline)
		return
	}

	// linkDown (2) → interface caiu
	if event.GenericTrap == 2 {
		r.log.Warn("linkDown recebido", zap.String("device", device.IPAddress))
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusWarning)
		return
	}

	// linkUp (3) → interface voltou
	if event.GenericTrap == 3 {
		r.log.Info("linkUp recebido", zap.String("device", device.IPAddress))
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusOnline)
		return
	}

	// Traps v2c padrão: verificar OID de sysUptime e snmpTrapOID
	for _, v := range event.Variables {
		switch v.OID {
		case "1.3.6.1.6.3.1.1.4.1.0": // snmpTrapOID
			r.handleSNMPTrapOID(device, v.Value, event)
		}
	}
}

// handleSNMPTrapOID processa o OID do trap para v2c/v3
func (r *TrapReceiver) handleSNMPTrapOID(device *models.Device, trapOID string, event TrapEvent) {
	r.log.Info("Processando trap OID",
		zap.String("device", device.IPAddress),
		zap.String("trap_oid", trapOID),
	)

	// Mapeamento básico de OIDs de trap conhecidos
	switch trapOID {
	case "1.3.6.1.6.3.1.1.5.1": // coldStart
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusOnline)
	case "1.3.6.1.6.3.1.1.5.2": // warmStart
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusOnline)
	case "1.3.6.1.6.3.1.1.5.3": // linkDown
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusWarning)
	case "1.3.6.1.6.3.1.1.5.4": // linkUp
		_ = r.deviceRepo.UpdateStatus(device.ID, models.DeviceStatusOnline)
	case "1.3.6.1.6.3.1.1.5.5": // authenticationFailure
		r.log.Warn("Falha de autenticação SNMP reportada pelo dispositivo",
			zap.String("device", device.IPAddress),
		)
	}
}
