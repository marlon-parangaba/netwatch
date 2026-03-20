package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AlertSeverity string
type AlertStatus string
type AlertCondition string
type NotificationChannel string

const (
	SeverityCritical AlertSeverity = "critical"
	SeverityHigh     AlertSeverity = "high"
	SeverityMedium   AlertSeverity = "medium"
	SeverityLow      AlertSeverity = "low"
	SeverityInfo     AlertSeverity = "info"
)

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusAcked    AlertStatus = "acknowledged"
)

const (
	ConditionGt  AlertCondition = "gt"  // greater than
	ConditionGte AlertCondition = "gte" // greater than or equal
	ConditionLt  AlertCondition = "lt"  // less than
	ConditionLte AlertCondition = "lte" // less than or equal
	ConditionEq  AlertCondition = "eq"  // equal
	ConditionNeq AlertCondition = "neq" // not equal
	ConditionDown AlertCondition = "down" // device down
)

const (
	ChannelEmail    NotificationChannel = "email"
	ChannelTelegram NotificationChannel = "telegram"
	ChannelSMS      NotificationChannel = "sms"
	ChannelWebhook  NotificationChannel = "webhook"
)

// AlertRule define uma regra de alerta
type AlertRule struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	DeviceID    *uuid.UUID     `gorm:"type:uuid;index" json:"device_id"` // nil = aplica a todos
	MetricType  MetricType     `gorm:"type:varchar(50);not null" json:"metric_type"`
	Condition   AlertCondition `gorm:"type:varchar(10);not null" json:"condition"`
	Threshold   float64        `json:"threshold"`
	Severity    AlertSeverity  `gorm:"type:varchar(20);not null;default:'medium'" json:"severity"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	// Canais de notificação em JSON: [{"channel":"telegram","config":{...}}]
	Notifications datatypes.JSON `gorm:"type:jsonb" json:"notifications"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// AlertEvent é um evento de alerta disparado
// Retenção: 6 meses com rolling window
type AlertEvent struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RuleID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"rule_id"`
	DeviceID    uuid.UUID     `gorm:"type:uuid;not null;index" json:"device_id"`
	Status      AlertStatus   `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	Severity    AlertSeverity `gorm:"type:varchar(20);not null" json:"severity"`
	Value       float64       `json:"value"` // valor que disparou o alerta
	Message     string        `gorm:"type:text" json:"message"`
	TriggeredAt time.Time     `gorm:"not null;index" json:"triggered_at"`
	ResolvedAt  *time.Time    `json:"resolved_at"`
	AckedAt     *time.Time    `json:"acked_at"`
	AckedByID   *uuid.UUID    `gorm:"type:uuid" json:"acked_by_id"`
	Duration    *int64        `json:"duration_seconds"` // calculado ao resolver

	// Relações
	Rule   AlertRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
	Device Device    `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// DowntimeEvent registra períodos de indisponibilidade de dispositivos.
// Retenção: vitalícia (usado para cálculo de MTBF, MTTR e duração de bateria).
type DowntimeEvent struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DeviceID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	StartedAt  time.Time  `gorm:"not null;index" json:"started_at"`
	EndedAt    *time.Time `json:"ended_at"`
	DurationSec *int64    `json:"duration_seconds"` // calculado ao restaurar
	Cause      string     `gorm:"type:varchar(255)" json:"cause"` // e.g., "ping_down", "snmp_timeout"
	Notes      string     `gorm:"type:text" json:"notes"`

	Device Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// MTBF e MTTR calculados sob demanda
type DeviceReliabilityStats struct {
	DeviceID     uuid.UUID `json:"device_id"`
	TotalDowntime int64    `json:"total_downtime_seconds"`
	DowntimeCount int      `json:"downtime_count"`
	MTBF         float64   `json:"mtbf_hours"` // Mean Time Between Failures
	MTTR         float64   `json:"mttr_hours"` // Mean Time To Recover
	Availability float64   `json:"availability_percent"`
	Since        time.Time `json:"since"`
}

// CreateAlertRuleRequest DTO de entrada para criar regra de alerta
type CreateAlertRuleRequest struct {
	Name          string         `json:"name" validate:"required,min=1,max=255"`
	Description   string         `json:"description"`
	DeviceID      *uuid.UUID     `json:"device_id"`
	MetricType    MetricType     `json:"metric_type" validate:"required"`
	Condition     AlertCondition `json:"condition" validate:"required,oneof=gt gte lt lte eq neq down"`
	Threshold     float64        `json:"threshold"`
	Severity      AlertSeverity  `json:"severity" validate:"required,oneof=critical high medium low info"`
	Notifications datatypes.JSON `json:"notifications"`
}
