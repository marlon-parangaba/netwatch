package handlers

import (
	"errors"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/ertel/netwatch-backend/internal/services"
	"github.com/ertel/netwatch-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DeviceHandler struct {
	deviceSvc services.DeviceService
	log       *zap.Logger
}

func NewDeviceHandler(deviceSvc services.DeviceService, log *zap.Logger) *DeviceHandler {
	return &DeviceHandler{deviceSvc: deviceSvc, log: log}
}

// List godoc
// GET /api/devices
func (h *DeviceHandler) List(c *fiber.Ctx) error {
	page, limit := ParsePage(c)

	filters := repository.DeviceFilters{
		Status: c.Query("status"),
		Type:   c.Query("type"),
		Search: c.Query("search"),
	}

	if enabledStr := c.Query("enabled"); enabledStr != "" {
		enabled := enabledStr == "true"
		filters.Enabled = &enabled
	}

	devices, total, err := h.deviceSvc.List(page, limit, filters)
	if err != nil {
		h.log.Error("erro ao listar dispositivos", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.OKPaginated(c, devices, page, limit, total)
}

// Get godoc
// GET /api/devices/:id
func (h *DeviceHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "ID inválido")
	}

	device, err := h.deviceSvc.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			return utils.NotFound(c, "dispositivo não encontrado")
		}
		h.log.Error("erro ao buscar dispositivo", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.OK(c, device)
}

// Create godoc
// POST /api/devices
func (h *DeviceHandler) Create(c *fiber.Ctx) error {
	var req models.CreateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "corpo da requisição inválido")
	}

	device, err := h.deviceSvc.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrDeviceDuplicate) {
			return utils.Conflict(c, "já existe um dispositivo com este endereço IP")
		}
		h.log.Error("erro ao criar dispositivo", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.Created(c, device)
}

// Update godoc
// PUT /api/devices/:id
func (h *DeviceHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "ID inválido")
	}

	var req models.UpdateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "corpo da requisição inválido")
	}

	device, err := h.deviceSvc.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			return utils.NotFound(c, "dispositivo não encontrado")
		}
		h.log.Error("erro ao atualizar dispositivo", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.OK(c, device)
}

// Delete godoc
// DELETE /api/devices/:id
func (h *DeviceHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "ID inválido")
	}

	if err := h.deviceSvc.Delete(id); err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			return utils.NotFound(c, "dispositivo não encontrado")
		}
		h.log.Error("erro ao deletar dispositivo", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.NoContent(c)
}
