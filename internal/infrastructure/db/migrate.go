package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/neko-dream/server/internal/infrastructure/config"
)

type Migrator struct {
	config *config.Config
}

func NewMigrator(config *config.Config) *Migrator {
	return &Migrator{config: config}
}

func (m *Migrator) Up() {
	pgx, err := sql.Open("pgx", m.config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	mi, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
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
	pgx, err := sql.Open("pgx", m.config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	pgx.Exec("SET session_replication_role = replica;")
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	mi, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err := mi.Drop(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println(err)
		}
	}
}
