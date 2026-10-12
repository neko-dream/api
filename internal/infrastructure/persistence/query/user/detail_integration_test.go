package user_query_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	user_query "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	profile "github.com/neko-dream/api/internal/infrastructure/persistence/query/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/internal/test/integration"
)

func TestUserProfileOptionalDemographicsIntegration(t *testing.T) {
	conn := integration.NewDatabase(t)
	manager := db.NewDBManager(conn)
	encryptor, err := crypto.NewEncryptor(&config.Config{ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"})
	if err != nil {
		t.Fatal(err)
	}
	users := repository.NewUserRepository(manager, nil, encryptor)
	handler := profile.NewDetailHandler(manager, encryptor)
	ctx := context.Background()
	for _, scenario := range []string{"absent", "empty", "encrypted", "corrupt"} {
		t.Run(scenario, func(t *testing.T) {
			id := shared.NewUUID[user.User]()
			displayID, name := "profile_"+uuid.NewString(), "プロフィールテスト"
			usr := user.NewUser(id, &displayID, &name, uuid.NewString(), "PASSWORD", nil)
			usr.ChangeEmail("profile@example.com")
			if err := users.Create(ctx, usr); err != nil {
				t.Fatal(err)
			}
			if scenario != "absent" {
				var city any
				if scenario == "encrypted" {
					city, err = encryptor.EncryptString(ctx, "検証市")
					if err != nil {
						t.Fatal(err)
					}
				} else if scenario == "corrupt" {
					city = "invalid-ciphertext"
				}
				if _, err := conn.ExecContext(ctx, `INSERT INTO user_demographics (user_demographics_id,user_id,city) VALUES ($1,$2,$3)`, uuid.New(), id.UUID(), city); err != nil {
					t.Fatal(err)
				}
			}
			result, err := handler.Execute(ctx, user_query.DetailInput{UserID: id})
			if scenario == "corrupt" {
				if err == nil {
					t.Fatal("corrupt demographics must not be treated as absent")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result.User.DisplayName != name || result.User.Email == nil || *result.User.Email != "profile@example.com" {
				t.Fatal("profile/email missing")
			}
			if scenario == "absent" && result.User.UserDemographic != nil {
				t.Fatal("absent demographics must remain nil")
			}
			if scenario != "absent" && result.User.UserDemographic == nil {
				t.Fatal("existing demographics missing")
			}
			if scenario == "encrypted" && (result.User.City == nil || *result.User.City != "検証市") {
				t.Fatal("demographics not decrypted")
			}
		})
	}
	if _, err := handler.Execute(ctx, user_query.DetailInput{UserID: shared.NewUUID[user.User]()}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("missing user: %v", err)
	}
}
