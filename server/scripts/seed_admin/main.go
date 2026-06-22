// admin_users 許可リストにローカル/開発用の管理者を登録するシードツール。
// ログインは Google 認証なので、ここでは email(Googleアカウント) と表示名を登録するだけ。
//
//	go run ./scripts/seed_admin -db <url> -email <google-email> -name <name>
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	dbURL := flag.String("db", env("DATABASE_URL", "postgres://kotohiro:kotohiro@localhost:5432/kotohiro_adm2"), "database url")
	email := flag.String("email", "admin@example.com", "admin google account email")
	name := flag.String("name", "管理者", "display name")
	flag.Parse()

	db, err := sql.Open("pgx", *dbURL)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(
		`INSERT INTO admin_users (admin_user_id, email, display_name, role)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (email) DO UPDATE
		 SET display_name = EXCLUDED.display_name,
		     updated_at = CURRENT_TIMESTAMP`,
		uuid.New(), *email, *name, 10,
	); err != nil {
		log.Fatalf("insert: %v", err)
	}

	fmt.Printf("seeded admin user (allowlist): %s\n", *email)
}
