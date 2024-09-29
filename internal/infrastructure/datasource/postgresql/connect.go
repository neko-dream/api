package postgresql

import (
	"database/sql"

	"github.com/neko-dream/server/internal/infrastructure/config"
)

func Connect(config *config.Config) *sql.DB {
	pgx, err := sql.Open(
		"pgx",
		config.DatabaseURL,
	)
	if err != nil {
		panic(err)
	}
	return pgx
}
