package models

import (
	"time"

	"github.com/google/uuid"
)

type MetricType string

const (
	MetricCPU         MetricType = "cpu"
	MetricMemory      MetricType = "memory"
	MetricDisk        MetricType = "disk"
	MetricIfInOctets  MetricType = "if_in_octets"
	MetricIfOutOctets MetricType = "if_out_octets"
	MetricIfInErrors  MetricType = "if_in_errors"
	MetricIfOutErrors MetricType = "if_out_errors"
	MetricPingRTT     MetricType = "ping_rtt"
	MetricTemperature MetricType = "temperature"
	MetricUptime      MetricType = "uptime"
	MetricCustom      MetricType = "custom"
)

// Metric armazena uma única leitura de métrica de um dispositivo.
// Particionado por mês no PostgreSQL para retenção de 1 ano.
type Metric struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	DeviceID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	Type       MetricType `gorm:"type:varchar(50);not null;index" json:"type"`
	InterfaceIndex int    `gorm:"default:0" json:"interface_index"` // para métricas de interface
	OID        string     `gorm:"type:varchar(255)" json:"oid"`
	Value      float64    `gorm:"not null" json:"value"`
	Unit       string     `gorm:"type:varchar(20)" json:"unit"`
	CollectedAt time.Time `gorm:"not null;index" json:"collected_at"`
}

// InterfaceInfo armazena informações sobre as interfaces de um dispositivo
type InterfaceInfo struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	IfIndex     int       `gorm:"not null" json:"if_index"`
	IfDescr     string    `gorm:"type:varchar(255)" json:"if_descr"`
	IfType      int       `json:"if_type"`
	IfSpeed     int64     `json:"if_speed"` // bps
	IfAdminStatus int     `json:"if_admin_status"` // 1=up, 2=down
	IfOperStatus  int     `json:"if_oper_status"`  // 1=up, 2=down
	IfAlias     string    `gorm:"type:varchar(255)" json:"if_alias"`
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MetricQuery são os parâmetros para consultar métricas históricas
type MetricQuery struct {
	DeviceID  string     `query:"device_id"`
	Type      MetricType `query:"type"`
	IfIndex   int        `query:"if_index"`
	From      time.Time  `query:"from"`
	To        time.Time  `query:"to"`
	Limit     int        `query:"limit"`
}

// MetricSummary é um resumo agregado de métricas
type MetricSummary struct {
	DeviceID uuid.UUID  `json:"device_id"`
	Type     MetricType `json:"type"`
	Min      float64    `json:"min"`
	Max      float64    `json:"max"`
	Avg      float64    `json:"avg"`
	Last     float64    `json:"last"`
	LastAt   time.Time  `json:"last_at"`
}
