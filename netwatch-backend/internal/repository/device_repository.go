package repository

import (
	"errors"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device *models.Device) error
	FindByID(id uuid.UUID) (*models.Device, error)
	FindByIP(ip string) (*models.Device, error)
	FindAll(page, limit int, filters DeviceFilters) ([]models.Device, int64, error)
	Update(device *models.Device) error
	UpdateStatus(id uuid.UUID, status models.DeviceStatus) error
	UpdateLastSeen(id uuid.UUID) error
	Delete(id uuid.UUID) error
	BulkCreate(devices []models.Device) error
}

type DeviceFilters struct {
	Status  string
	Type    string
	Enabled *bool
	Search  string
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(device *models.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepository) FindByID(id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &device, err
}

func (r *deviceRepository) FindByIP(ip string) (*models.Device, error) {
	var device models.Device
	err := r.db.Where("ip_address = ? AND deleted_at IS NULL", ip).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &device, err
}

func (r *deviceRepository) FindAll(page, limit int, filters DeviceFilters) ([]models.Device, int64, error) {
	var devices []models.Device
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Device{}).Where("deleted_at IS NULL")

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.Enabled != nil {
		query = query.Where("enabled = ?", *filters.Enabled)
	}
	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR ip_address ILIKE ? OR hostname ILIKE ?", search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&devices).Error

	return devices, total, err
}

func (r *deviceRepository) Update(device *models.Device) error {
	return r.db.Save(device).Error
}

func (r *deviceRepository) UpdateStatus(id uuid.UUID, status models.DeviceStatus) error {
	return r.db.Model(&models.Device{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"last_seen_at": gorm.Expr("CASE WHEN ? = 'online' THEN NOW() ELSE last_seen_at END", status),
		}).Error
}

func (r *deviceRepository) UpdateLastSeen(id uuid.UUID) error {
	return r.db.Model(&models.Device{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_seen_at":  gorm.Expr("NOW()"),
			"last_polled_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *deviceRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Device{}).Error
}

func (r *deviceRepository) BulkCreate(devices []models.Device) error {
	return r.db.CreateInBatches(devices, 100).Error
}
