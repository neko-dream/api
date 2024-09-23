package postgresql

import (
	"database/sql"
	"os"
)

func Connect() *sql.DB {
	pgx, err := sql.Open(
		"pgx",
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		panic(err)
	}
	return pgx
}
