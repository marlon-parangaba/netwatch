package handlers

import (
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/ertel/netwatch-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MetricHandler struct {
	metricRepo repository.MetricRepository
	log        *zap.Logger
}

func NewMetricHandler(metricRepo repository.MetricRepository, log *zap.Logger) *MetricHandler {
	return &MetricHandler{metricRepo: metricRepo, log: log}
}

// GetHistory godoc
// GET /api/metrics/:deviceId/history
func (h *MetricHandler) GetHistory(c *fiber.Ctx) error {
	deviceID, err := uuid.Parse(c.Params("deviceId"))
	if err != nil {
		return utils.BadRequest(c, "deviceId inválido")
	}

	metricType := models.MetricType(c.Query("type"))

	// Padrão: últimas 24h
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	limit := c.QueryInt("limit", 1000)

	metrics, err := h.metricRepo.GetHistory(deviceID, metricType, from, to, limit)
	if err != nil {
		h.log.Error("erro ao buscar histórico de métricas", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.OK(c, metrics)
}

// GetSummary godoc
// GET /api/metrics/:deviceId/summary
func (h *MetricHandler) GetSummary(c *fiber.Ctx) error {
	deviceID, err := uuid.Parse(c.Params("deviceId"))
	if err != nil {
		return utils.BadRequest(c, "deviceId inválido")
	}

	metricType := models.MetricType(c.Query("type"))
	if metricType == "" {
		return utils.BadRequest(c, "parâmetro 'type' é obrigatório")
	}

	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	summary, err := h.metricRepo.GetSummary(deviceID, metricType, from, to)
	if err != nil {
		h.log.Error("erro ao calcular resumo de métricas", zap.Error(err))
		return utils.InternalError(c, "")
	}

	return utils.OK(c, summary)
}
