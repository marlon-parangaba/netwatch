package repository

import (
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetricRepository interface {
	Save(metric *models.Metric) error
	SaveBatch(metrics []models.Metric) error
	GetHistory(deviceID uuid.UUID, metricType models.MetricType, from, to time.Time, limit int) ([]models.Metric, error)
	GetLatest(deviceID uuid.UUID, metricType models.MetricType) (*models.Metric, error)
	GetSummary(deviceID uuid.UUID, metricType models.MetricType, from, to time.Time) (*models.MetricSummary, error)
	DeleteOlderThan(cutoff time.Time) (int64, error)
}

type metricRepository struct {
	db *gorm.DB
}

func NewMetricRepository(db *gorm.DB) MetricRepository {
	return &metricRepository{db: db}
}

func (r *metricRepository) Save(metric *models.Metric) error {
	return r.db.Create(metric).Error
}

func (r *metricRepository) SaveBatch(metrics []models.Metric) error {
	return r.db.CreateInBatches(metrics, 500).Error
}

func (r *metricRepository) GetHistory(deviceID uuid.UUID, metricType models.MetricType, from, to time.Time, limit int) ([]models.Metric, error) {
	var metrics []models.Metric

	if limit <= 0 {
		limit = 1000
	}

	query := r.db.Where("device_id = ? AND collected_at BETWEEN ? AND ?", deviceID, from, to).
		Order("collected_at DESC").
		Limit(limit)

	if metricType != "" {
		query = query.Where("type = ?", metricType)
	}

	return metrics, query.Find(&metrics).Error
}

func (r *metricRepository) GetLatest(deviceID uuid.UUID, metricType models.MetricType) (*models.Metric, error) {
	var metric models.Metric
	err := r.db.Where("device_id = ? AND type = ?", deviceID, metricType).
		Order("collected_at DESC").
		First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *metricRepository) GetSummary(deviceID uuid.UUID, metricType models.MetricType, from, to time.Time) (*models.MetricSummary, error) {
	var summary models.MetricSummary

	type result struct {
		Min  float64
		Max  float64
		Avg  float64
		Last float64
		LastAt time.Time
	}

	var res result
	err := r.db.Raw(`
		SELECT
			MIN(value) as min,
			MAX(value) as max,
			AVG(value) as avg,
			(SELECT value FROM metrics WHERE device_id = ? AND type = ? AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) as last,
			(SELECT collected_at FROM metrics WHERE device_id = ? AND type = ? AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) as last_at
		FROM metrics
		WHERE device_id = ? AND type = ? AND collected_at BETWEEN ? AND ?
	`, deviceID, metricType, from, to,
		deviceID, metricType, from, to,
		deviceID, metricType, from, to).
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	summary = models.MetricSummary{
		DeviceID: deviceID,
		Type:     metricType,
		Min:      res.Min,
		Max:      res.Max,
		Avg:      res.Avg,
		Last:     res.Last,
		LastAt:   res.LastAt,
	}

	return &summary, nil
}

func (r *metricRepository) DeleteOlderThan(cutoff time.Time) (int64, error) {
	result := r.db.Where("collected_at < ?", cutoff).Delete(&models.Metric{})
	return result.RowsAffected, result.Error
}
