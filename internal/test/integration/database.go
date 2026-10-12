// Package integration provides isolated databases for tests that need commits,
// concurrent connections, or schema changes rather than a rollback-only test.
package integration

import (
	"database/sql"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
)

func NewDatabase(t *testing.T) *sql.DB {
	t.Helper()
	// Retain the existing CI variable for compatibility.
	raw := os.Getenv("MAIL_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("set MAIL_TEST_DATABASE_URL to run database integration tests")
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	admin, err := sql.Open("pgx", raw)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Close() })
	schema := "integration_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.Exec(`CREATE SCHEMA "` + schema + `"`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(`DROP SCHEMA "` + schema + `" CASCADE`); err != nil {
			t.Error(err)
		}
	})
	q := u.Query()
	q.Set("search_path", schema+",public")
	u.RawQuery = q.Encode()
	if err := db.NewMigrator(&config.Config{DatabaseURL: u.String()}).Up(); err != nil {
		t.Fatal(err)
	}
	conn, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}
