package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DeviceStatus string
type DeviceType string
type SNMPVersion string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusWarning DeviceStatus = "warning"
	DeviceStatusUnknown DeviceStatus = "unknown"
)

const (
	DeviceTypeMikrotik  DeviceType = "mikrotik"
	DeviceTypeCisco     DeviceType = "cisco"
	DeviceTypeHuawei    DeviceType = "huawei"
	DeviceTypeJuniper   DeviceType = "juniper"
	DeviceTypeUbiquiti  DeviceType = "ubiquiti"
	DeviceTypeGeneric   DeviceType = "generic"
)

const (
	SNMPv1  SNMPVersion = "v1"
	SNMPv2c SNMPVersion = "v2c"
	SNMPv3  SNMPVersion = "v3"
)

type Device struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Hostname    string         `gorm:"type:varchar(255)" json:"hostname"`
	IPAddress   string         `gorm:"type:varchar(45);not null;index" json:"ip_address"`
	Type        DeviceType     `gorm:"type:varchar(50);default:'generic'" json:"type"`
	Status      DeviceStatus   `gorm:"type:varchar(20);default:'unknown'" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"type:varchar(255)" json:"location"`

	// SNMP Config
	SNMPVersion   SNMPVersion `gorm:"type:varchar(5);default:'v2c'" json:"snmp_version"`
	SNMPCommunity string      `gorm:"type:varchar(100);default:'public'" json:"snmp_community,omitempty"`
	SNMPPort      int         `gorm:"default:161" json:"snmp_port"`

	// SNMPv3 fields
	SNMPv3Username   string `gorm:"type:varchar(100)" json:"snmpv3_username,omitempty"`
	SNMPv3AuthProto  string `gorm:"type:varchar(10)" json:"snmpv3_auth_proto,omitempty"`
	SNMPv3AuthKey    string `gorm:"type:varchar(255)" json:"-"`
	SNMPv3PrivProto  string `gorm:"type:varchar(10)" json:"snmpv3_priv_proto,omitempty"`
	SNMPv3PrivKey    string `gorm:"type:varchar(255)" json:"-"`

	// Polling
	PollingInterval int  `gorm:"default:60" json:"polling_interval"` // segundos
	Enabled         bool `gorm:"default:true" json:"enabled"`

	// Metadata descoberta via SNMP
	SysName    string `gorm:"type:varchar(255)" json:"sys_name"`
	SysDescr   string `gorm:"type:text" json:"sys_descr"`
	SysOID     string `gorm:"type:varchar(255)" json:"sys_oid"`
	SysContact string `gorm:"type:varchar(255)" json:"sys_contact"`

	// Tags customizadas (JSON)
	Tags datatypes.JSON `gorm:"type:jsonb" json:"tags"`

	// Timestamps de status
	LastSeenAt    *time.Time `json:"last_seen_at"`
	LastPolledAt  *time.Time `json:"last_polled_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Device) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// CreateDeviceRequest DTO de entrada para criar dispositivo
type CreateDeviceRequest struct {
	Name          string      `json:"name" validate:"required,min=1,max=255"`
	IPAddress     string      `json:"ip_address" validate:"required,ip"`
	Hostname      string      `json:"hostname" validate:"omitempty,max=255"`
	Type          DeviceType  `json:"type" validate:"omitempty,oneof=mikrotik cisco huawei juniper ubiquiti generic"`
	Description   string      `json:"description"`
	Location      string      `json:"location"`
	SNMPVersion   SNMPVersion `json:"snmp_version" validate:"omitempty,oneof=v1 v2c v3"`
	SNMPCommunity string      `json:"snmp_community"`
	SNMPPort      int         `json:"snmp_port" validate:"omitempty,min=1,max=65535"`
	PollingInterval int       `json:"polling_interval" validate:"omitempty,min=10,max=3600"`
}

// UpdateDeviceRequest DTO de entrada para atualizar dispositivo
type UpdateDeviceRequest struct {
	Name            string      `json:"name" validate:"omitempty,min=1,max=255"`
	Hostname        string      `json:"hostname"`
	Type            DeviceType  `json:"type" validate:"omitempty,oneof=mikrotik cisco huawei juniper ubiquiti generic"`
	Description     string      `json:"description"`
	Location        string      `json:"location"`
	SNMPVersion     SNMPVersion `json:"snmp_version" validate:"omitempty,oneof=v1 v2c v3"`
	SNMPCommunity   string      `json:"snmp_community"`
	SNMPPort        int         `json:"snmp_port" validate:"omitempty,min=1,max=65535"`
	PollingInterval int         `json:"polling_interval" validate:"omitempty,min=10,max=3600"`
	Enabled         *bool       `json:"enabled"`
}

// DiscoverRequest DTO para iniciar discovery de rede
type DiscoverRequest struct {
	Network       string      `json:"network" validate:"required,cidr"`
	SNMPVersion   SNMPVersion `json:"snmp_version" validate:"omitempty,oneof=v1 v2c v3"`
	SNMPCommunity string      `json:"snmp_community"`
	SNMPPort      int         `json:"snmp_port" validate:"omitempty,min=1,max=65535"`
	Timeout       int         `json:"timeout_seconds" validate:"omitempty,min=1,max=60"`
}
