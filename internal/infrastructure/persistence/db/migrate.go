package db

import (
	"database/sql"
	"log"

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

func (m *Migrator) Up() {
	d, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		panic(err)
	}

	pgx, err := sql.Open("pgx", m.config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	mi, err := migrate.NewWithInstance(
		"iofs",
		d,
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := mi.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println(err)
		}
	}
}

func (m *Migrator) Down() {
	d, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		panic(err)
	}

	pgx, err := sql.Open("pgx", m.config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	mi, err := migrate.NewWithInstance(
		"iofs",
		d,
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := mi.Down(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println(err)
		}
	}
}
