package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/migrations"
)

type Migrator struct {
	config *config.Config
}

func NewMigrator(config *config.Config) *Migrator {
	return &Migrator{config: config}
}

func (m *Migrator) Up() error {
	d, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	pgx, err := sql.Open("pgx", m.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer pgx.Close()

	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	mi, err := migrate.NewWithInstance(
		"iofs",
		d,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := mi.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Down() error {
	d, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	pgx, err := sql.Open("pgx", m.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer pgx.Close()

	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	mi, err := migrate.NewWithInstance(
		"iofs",
		d,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := mi.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}
