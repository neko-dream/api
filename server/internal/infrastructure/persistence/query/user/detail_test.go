package user_query

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	cryptoinfra "github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
)

func TestDetailHandlerOptionalDemographics(t *testing.T) {
	databaseURL := os.Getenv("ADMIN_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("ADMIN_TEST_DATABASE_URL required")
	}
	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	manager := db.NewDBManager(conn)
	encryptor, err := cryptoinfra.NewGCMEncryptor([]byte("local-dev-key16!"))
	if err != nil {
		t.Fatal(err)
	}
	query := NewDetailHandler(manager, encryptor)

	for _, scenario := range []string{"missing", "empty", "encrypted", "corrupt demographics", "corrupt email"} {
		t.Run(scenario, func(t *testing.T) {
			err := manager.TestTx(context.Background(), func(ctx context.Context) error {
				queries := manager.GetQueries(ctx)
				id := shared.NewUUID[user.User]()
				cipher, err := encryptor.EncryptString(ctx, "profile@example.com")
				if err != nil {
					return err
				}
				if scenario == "corrupt email" {
					cipher = "invalid-ciphertext"
				}
				if err := queries.CreateUser(ctx, model.CreateUserParams{
					UserID: id.UUID(), CreatedAt: time.Now(),
					DisplayID:   sql.NullString{String: id.String(), Valid: true},
					DisplayName: sql.NullString{String: "プロフィール確認", Valid: true},
					Email:       sql.NullString{String: cipher, Valid: true}, EmailVerified: true,
				}); err != nil {
					return err
				}
				if err := queries.CreateUserAuth(ctx, model.CreateUserAuthParams{
					UserAuthID: shared.NewUUID[user.User]().UUID(), UserID: id.UUID(),
					Provider: "PASSWORD", Subject: id.String(), CreatedAt: time.Now(),
				}); err != nil {
					return err
				}
				if scenario != "missing" && scenario != "corrupt email" {
					demographics := model.UpdateOrCreateUserDemographicParams{
						UserDemographicsID: shared.NewUUID[user.User]().UUID(), UserID: id.UUID(),
					}
					if scenario == "encrypted" {
						city, err := encryptor.EncryptString(ctx, "鯖江市")
						if err != nil {
							return err
						}
						demographics.City = sql.NullString{String: city, Valid: true}
					}
					if scenario == "corrupt demographics" {
						demographics.City = sql.NullString{String: "invalid-ciphertext", Valid: true}
					}
					if err := queries.UpdateOrCreateUserDemographic(ctx, demographics); err != nil {
						return err
					}
				}
				out, err := query.Execute(ctx, userquery.DetailInput{UserID: id})
				if scenario == "corrupt demographics" {
					if err == nil {
						t.Error("corrupt demographics must not be treated as absent")
					}
					return nil
				}
				if err != nil {
					return err
				}
				if out.User.DisplayName != "プロフィール確認" || out.User.DisplayID != id.String() {
					t.Error("profile fields were lost")
				}
				if scenario == "corrupt email" {
					if out.User.Email != nil || out.User.EmailVerified {
						t.Error("corrupt email must not be exposed or marked verified")
					}
				} else if out.User.Email == nil || *out.User.Email != "profile@example.com" || !out.User.EmailVerified {
					t.Error("verified email was not decrypted")
				}
				if scenario == "missing" || scenario == "corrupt email" {
					if out.User.UserDemographic != nil {
						t.Error("absent demographics must remain nil")
					}
					if _, err := queries.GetUserDemographicByUserID(ctx, id.UUID()); !errors.Is(err, sql.ErrNoRows) {
						t.Error("reading the profile must not create demographics")
					}
				} else if out.User.UserDemographic == nil {
					t.Error("existing demographics were lost")
				} else if scenario == "encrypted" && (out.User.City == nil || *out.User.City != "鯖江市") {
					t.Error("existing demographics were not decrypted")
				} else if scenario == "empty" && (out.User.City != nil || out.User.Gender != nil || out.User.DateOfBirth != nil) {
					t.Error("empty demographics were filled with invented data")
				}
				if _, err := query.Execute(ctx, userquery.DetailInput{UserID: shared.NewUUID[user.User]()}); !errors.Is(err, sql.ErrNoRows) {
					t.Error("missing user must remain an error")
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
