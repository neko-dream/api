package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/pkg/utils"
)

func init() {
	utils.LoadEnv()
}

func main() {
	ctx := context.Background()
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// トランザクション開始
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// 1. 一時カラムの追加
	if _, err := tx.ExecContext(ctx, `
		ALTER TABLE user_demographics
		ADD COLUMN gender_encrypted VARCHAR(255),
		ADD COLUMN year_of_birth_encrypted VARCHAR(255),
		ADD COLUMN city_encrypted VARCHAR(255),
		ADD COLUMN prefecture_encrypted VARCHAR(255)
	`); err != nil {
		log.Fatal("Failed to add encrypted columns:", err)
	}

	encryptor := crypto.NewCBCEncryptor([]byte(os.Getenv("ENCRYPTION_SECRET")))

	// 2. 既存データの暗号化と移行
	rows, err := tx.QueryContext(ctx, `
		SELECT user_id, gender, year_of_birth, city, prefecture
		FROM user_demographics
		WHERE gender IS NOT NULL
		   OR year_of_birth IS NOT NULL
		   OR city IS NOT NULL
		   OR prefecture IS NOT NULL
	`)
	if err != nil {
		log.Fatal("Failed to query existing data:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var gender sql.NullInt16
		var yearOfBirth sql.NullInt32
		var city, prefecture sql.NullString

		if err := rows.Scan(&id, &gender, &yearOfBirth, &city, &prefecture); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// gender暗号化
		if gender.Valid {
			if encryptedGender, err := encryptor.EncryptString(ctx, fmt.Sprintf("%d", gender.Int16)); err == nil {
				_, err = tx.ExecContext(ctx, `
					UPDATE user_demographics
					SET gender_encrypted = $1
					WHERE user_id = $2
				`, encryptedGender, id)
				if err != nil {
					log.Printf("Error updating gender for id %s: %v", id, err)
				}
			}
		}

		// year_of_birth暗号化
		if yearOfBirth.Valid {
			if encryptedYear, err := encryptor.EncryptString(ctx, fmt.Sprintf("%d", yearOfBirth.Int32)); err == nil {
				_, err = tx.ExecContext(ctx, `
					UPDATE user_demographics
					SET year_of_birth_encrypted = $1
					WHERE user_id = $2
				`, encryptedYear, id)
				if err != nil {
					log.Printf("Error updating year_of_birth for id %s: %v", id, err)
				}
			}
		}

		// city暗号化
		if city.Valid {
			if encryptedCity, err := encryptor.EncryptString(ctx, city.String); err == nil {
				_, err = tx.ExecContext(ctx, `
					UPDATE user_demographics
					SET city_encrypted = $1
					WHERE user_id = $2
				`, encryptedCity, id)
				if err != nil {
					log.Printf("Error updating city for id %s: %v", id, err)
				}
			}
		}

		// prefecture暗号化
		if prefecture.Valid {
			if encryptedPrefecture, err := encryptor.EncryptString(ctx, prefecture.String); err == nil {
				_, err = tx.ExecContext(ctx, `
					UPDATE user_demographics
					SET prefecture_encrypted = $1
					WHERE user_id = $2
				`, encryptedPrefecture, id)
				if err != nil {
					log.Printf("Error updating prefecture for id %s: %v", id, err)
				}
			}
		}
	}

	// 3. 古いカラムの削除と新しいカラムのリネーム
	if _, err := tx.ExecContext(ctx, `
		ALTER TABLE user_demographics
		DROP COLUMN gender,
		DROP COLUMN year_of_birth,
		DROP COLUMN city,
		DROP COLUMN prefecture
	`); err != nil {
		log.Fatal("Failed to drop old columns:", err)
	}

	if _, err := tx.ExecContext(ctx, `
		ALTER TABLE user_demographics
		RENAME COLUMN gender_encrypted TO gender;
		ALTER TABLE user_demographics
		RENAME COLUMN year_of_birth_encrypted TO year_of_birth;
		ALTER TABLE user_demographics
		RENAME COLUMN city_encrypted TO city;
		ALTER TABLE user_demographics
		RENAME COLUMN prefecture_encrypted TO prefecture
	`); err != nil {
		log.Fatal("Failed to rename columns:", err)
	}

	// トランザクションのコミット
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	log.Println("Migration completed successfully")
}
