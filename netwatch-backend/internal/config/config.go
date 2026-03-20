package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	SNMP     SNMPConfig     `mapstructure:"snmp"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	CORS     CORSConfig     `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Env          string        `mapstructure:"env"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=America/Sao_Paulo",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

type JWTConfig struct {
	Secret                  string `mapstructure:"secret"`
	ExpirationHours         int    `mapstructure:"expiration_hours"`
	RefreshExpirationHours  int    `mapstructure:"refresh_expiration_hours"`
}

type SNMPConfig struct {
	TrapPort  int           `mapstructure:"trap_port"`
	Community string        `mapstructure:"community"`
	Timeout   time.Duration `mapstructure:"timeout"`
	Retries   int           `mapstructure:"retries"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type RateLimitConfig struct {
	Requests      int `mapstructure:"requests"`
	WindowSeconds int `mapstructure:"window_seconds"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Variáveis de ambiente sobrescrevem o config.yaml
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bindings explícitos para as vars de ambiente mais comuns
	_ = v.BindEnv("database.host", "POSTGRES_HOST")
	_ = v.BindEnv("database.port", "POSTGRES_PORT")
	_ = v.BindEnv("database.name", "POSTGRES_DB")
	_ = v.BindEnv("database.user", "POSTGRES_USER")
	_ = v.BindEnv("database.password", "POSTGRES_PASSWORD")
	_ = v.BindEnv("database.sslmode", "POSTGRES_SSLMODE")
	_ = v.BindEnv("server.port", "API_PORT")
	_ = v.BindEnv("server.env", "API_ENV")
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("logging.level", "LOG_LEVEL")
	_ = v.BindEnv("logging.file", "LOG_FILE")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erro ao ler config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao desserializar config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}
