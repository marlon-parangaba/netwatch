package services

import (
	"errors"
	"fmt"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrDeviceNotFound  = errors.New("dispositivo não encontrado")
	ErrDeviceDuplicate = errors.New("já existe um dispositivo com este IP")
)

type DeviceService interface {
	Create(req models.CreateDeviceRequest) (*models.Device, error)
	GetByID(id uuid.UUID) (*models.Device, error)
	List(page, limit int, filters repository.DeviceFilters) ([]models.Device, int64, error)
	Update(id uuid.UUID, req models.UpdateDeviceRequest) (*models.Device, error)
	Delete(id uuid.UUID) error
}

type deviceService struct {
	deviceRepo repository.DeviceRepository
	log        *zap.Logger
}

func NewDeviceService(deviceRepo repository.DeviceRepository, log *zap.Logger) DeviceService {
	return &deviceService{deviceRepo: deviceRepo, log: log}
}

func (s *deviceService) Create(req models.CreateDeviceRequest) (*models.Device, error) {
	// Verificar se já existe dispositivo com este IP
	existing, err := s.deviceRepo.FindByIP(req.IPAddress)
	if err == nil && existing != nil {
		return nil, ErrDeviceDuplicate
	}

	snmpVersion := req.SNMPVersion
	if snmpVersion == "" {
		snmpVersion = models.SNMPv2c
	}

	snmpPort := req.SNMPPort
	if snmpPort == 0 {
		snmpPort = 161
	}

	pollingInterval := req.PollingInterval
	if pollingInterval == 0 {
		pollingInterval = 60
	}

	community := req.SNMPCommunity
	if community == "" {
		community = "public"
	}

	deviceType := req.Type
	if deviceType == "" {
		deviceType = models.DeviceTypeGeneric
	}

	device := &models.Device{
		Name:            req.Name,
		IPAddress:       req.IPAddress,
		Hostname:        req.Hostname,
		Type:            deviceType,
		Description:     req.Description,
		Location:        req.Location,
		SNMPVersion:     snmpVersion,
		SNMPCommunity:   community,
		SNMPPort:        snmpPort,
		PollingInterval: pollingInterval,
		Status:          models.DeviceStatusUnknown,
		Enabled:         true,
	}

	if err := s.deviceRepo.Create(device); err != nil {
		return nil, fmt.Errorf("erro ao criar dispositivo: %w", err)
	}

	s.log.Info("Dispositivo criado", zap.String("id", device.ID.String()), zap.String("ip", device.IPAddress))
	return device, nil
}

func (s *deviceService) GetByID(id uuid.UUID) (*models.Device, error) {
	device, err := s.deviceRepo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrDeviceNotFound
	}
	return device, err
}

func (s *deviceService) List(page, limit int, filters repository.DeviceFilters) ([]models.Device, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	return s.deviceRepo.FindAll(page, limit, filters)
}

func (s *deviceService) Update(id uuid.UUID, req models.UpdateDeviceRequest) (*models.Device, error) {
	device, err := s.deviceRepo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Hostname != "" {
		device.Hostname = req.Hostname
	}
	if req.Type != "" {
		device.Type = req.Type
	}
	if req.Description != "" {
		device.Description = req.Description
	}
	if req.Location != "" {
		device.Location = req.Location
	}
	if req.SNMPVersion != "" {
		device.SNMPVersion = req.SNMPVersion
	}
	if req.SNMPCommunity != "" {
		device.SNMPCommunity = req.SNMPCommunity
	}
	if req.SNMPPort > 0 {
		device.SNMPPort = req.SNMPPort
	}
	if req.PollingInterval > 0 {
		device.PollingInterval = req.PollingInterval
	}
	if req.Enabled != nil {
		device.Enabled = *req.Enabled
	}

	if err := s.deviceRepo.Update(device); err != nil {
		return nil, fmt.Errorf("erro ao atualizar dispositivo: %w", err)
	}

	return device, nil
}

func (s *deviceService) Delete(id uuid.UUID) error {
	if _, err := s.deviceRepo.FindByID(id); errors.Is(err, repository.ErrNotFound) {
		return ErrDeviceNotFound
	}
	return s.deviceRepo.Delete(id)
}
