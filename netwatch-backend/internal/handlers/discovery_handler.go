package handlers

import (
	"context"
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/ertel/netwatch-backend/internal/services"
	snmpengine "github.com/ertel/netwatch-backend/internal/snmp"
	"github.com/ertel/netwatch-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DiscoveryHandler struct {
	deviceSvc  services.DeviceService
	deviceRepo repository.DeviceRepository
	pollerSvc  *services.PollerService
	log        *zap.Logger
}

func NewDiscoveryHandler(
	deviceSvc services.DeviceService,
	deviceRepo repository.DeviceRepository,
	pollerSvc *services.PollerService,
	log *zap.Logger,
) *DiscoveryHandler {
	return &DiscoveryHandler{
		deviceSvc:  deviceSvc,
		deviceRepo: deviceRepo,
		pollerSvc:  pollerSvc,
		log:        log,
	}
}

// Discover godoc
// POST /api/devices/discover
func (h *DiscoveryHandler) Discover(c *fiber.Ctx) error {
	var req models.DiscoverRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "corpo da requisição inválido")
	}

	if req.Network == "" {
		return utils.BadRequest(c, "campo 'network' é obrigatório (ex: 192.168.1.0/24)")
	}

	if req.SNMPCommunity == "" {
		req.SNMPCommunity = "public"
	}
	if req.SNMPVersion == "" {
		req.SNMPVersion = models.SNMPv2c
	}
	if req.SNMPPort == 0 {
		req.SNMPPort = 161
	}

	timeout := time.Duration(req.Timeout) * time.Second
	if timeout == 0 {
		timeout = 3 * time.Second
	}

	h.log.Info("Iniciando discovery", zap.String("network", req.Network))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cfg := snmpengine.DiscoveryConfig{
		Network:     req.Network,
		SNMPVersion: req.SNMPVersion,
		Community:   req.SNMPCommunity,
		Port:        req.SNMPPort,
		Timeout:     timeout,
		Retries:     1,
		Concurrency: 100,
	}

	results, err := snmpengine.Discover(ctx, cfg, h.log)
	if err != nil {
		h.log.Error("Erro no discovery", zap.Error(err))
		return utils.InternalError(c, "erro durante a varredura de rede")
	}

	// Montar resposta com quais são novos e quais já existem
	type DiscoveryItemResponse struct {
		IPAddress string             `json:"ip_address"`
		SysName   string             `json:"sys_name"`
		SysDescr  string             `json:"sys_descr"`
		Vendor    string             `json:"vendor"`
		Type      models.DeviceType  `json:"type"`
		IsNew     bool               `json:"is_new"`
	}

	var response []DiscoveryItemResponse
	for _, r := range results {
		_, err := h.deviceRepo.FindByIP(r.IPAddress)
		isNew := err != nil // ErrNotFound → é novo

		response = append(response, DiscoveryItemResponse{
			IPAddress: r.IPAddress,
			SysName:   r.SysInfo.Name,
			SysDescr:  r.SysInfo.Descr,
			Vendor:    r.Vendor,
			Type:      r.Type,
			IsNew:     isNew,
		})
	}

	return utils.OK(c, fiber.Map{
		"network":  req.Network,
		"found":    len(results),
		"devices":  response,
	})
}

// ImportDiscovered godoc
// POST /api/devices/discover/import
// Importa dispositivos descobertos que ainda não estão no banco
func (h *DiscoveryHandler) ImportDiscovered(c *fiber.Ctx) error {
	var body struct {
		IPs           []string           `json:"ips"`
		SNMPVersion   models.SNMPVersion `json:"snmp_version"`
		SNMPCommunity string             `json:"snmp_community"`
		SNMPPort      int                `json:"snmp_port"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequest(c, "corpo inválido")
	}

	if len(body.IPs) == 0 {
		return utils.BadRequest(c, "lista de IPs não pode ser vazia")
	}

	if body.SNMPCommunity == "" {
		body.SNMPCommunity = "public"
	}
	if body.SNMPVersion == "" {
		body.SNMPVersion = models.SNMPv2c
	}
	if body.SNMPPort == 0 {
		body.SNMPPort = 161
	}

	var created []models.Device
	var skipped []string

	for _, ip := range body.IPs {
		// Verificar se já existe
		if _, err := h.deviceRepo.FindByIP(ip); err == nil {
			skipped = append(skipped, ip)
			continue
		}

		device, err := h.deviceSvc.Create(models.CreateDeviceRequest{
			Name:          ip,
			IPAddress:     ip,
			SNMPVersion:   body.SNMPVersion,
			SNMPCommunity: body.SNMPCommunity,
			SNMPPort:      body.SNMPPort,
		})
		if err != nil {
			h.log.Warn("Falha ao importar device", zap.String("ip", ip), zap.Error(err))
			continue
		}

		created = append(created, *device)

		// Adicionar ao scheduler de polling imediatamente
		if h.pollerSvc != nil {
			h.pollerSvc.AddDevice(device)
		}
	}

	return utils.Created(c, fiber.Map{
		"created": created,
		"skipped": skipped,
	})
}
