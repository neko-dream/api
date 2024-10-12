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
		log.Println("Panic in pgx")
		panic(err)
	}
	log.Println("DatabaseURL:", m.config.DatabaseURL)
	log.Println("pgx:", pgx)
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		log.Println("Panic in postgres", err)
		panic(err)
	}
	mi, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Println("Panic in migrate")
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
		log.Println("Panic in pgx")
		panic(err)
	}
	if _, err := pgx.Exec("SET session_replication_role = replica;"); err != nil {
		log.Println("Error setting session replication role:", err)
	}
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		log.Println("Panic in postgres")
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
