package database

import (
	"fmt"
	"time"

	"github.com/ertel/netwatch-backend/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB é a instância global do banco de dados
var DB *gorm.DB

// Connect inicializa a conexão com o PostgreSQL
func Connect(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	gormLogger := logger.Default.LogMode(logger.Silent)
	if !cfg.IsProduction() {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao conectar (ping): %w", err)
	}

	DB = db
	log.Info("Conectado ao PostgreSQL com sucesso",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("db", cfg.Database.Name),
	)

	return db, nil
}

// Close encerra a conexão com o banco
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping verifica se o banco está acessível
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
