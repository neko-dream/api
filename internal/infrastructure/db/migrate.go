package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Migration() {
	pgx, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"catsdream",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println("a ", err)
		}
	}
}

func Down() {
	log.Println("a ", os.Getenv("DATABASE_URL"))
	pgx, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("a ", err)
		panic(err)
	}
	pgx.Exec("SET session_replication_role = replica;")
	driver, err := postgres.WithInstance(pgx, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"catsdream",
		driver,
	)
	if err := m.Drop(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println("a ", err)
		}
	}
}
