package database

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:embed ../../migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executa todas as migrations pendentes
func RunMigrations(db *gorm.DB, log *zap.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("falha ao obter sql.DB para migrations: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("falha ao criar driver de migrations: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("falha ao carregar migrations: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("falha ao inicializar migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("falha ao executar migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	log.Info("Migrations executadas com sucesso",
		zap.Uint("versao", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}
