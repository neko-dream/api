package user_query

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	cryptoinfra "github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"os"
	"testing"
	"time"
)

func TestManagementQuery(t *testing.T) {
	url := os.Getenv("ADMIN_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("ADMIN_TEST_DATABASE_URL required")
	}
	conn, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	manager := db.NewDBManager(conn)
	encryptor, err := cryptoinfra.NewGCMEncryptor([]byte("local-dev-key16!"))
	if err != nil {
		t.Fatal(err)
	}
	q := &managementQuery{db: manager, encryptor: encryptor}
	err = manager.TestTx(context.Background(), func(ctx context.Context) error {
		queries := manager.GetQueries(ctx)
		id := shared.NewUUID[user.User]()
		second := shared.NewUUID[user.User]()
		prefix := id.String()
		cipher, err := encryptor.EncryptString(ctx, "city@example.com")
		if err != nil {
			return err
		}
		for index, uid := range []shared.UUID[user.User]{id, second} {
			name := prefix + "-second"
			if index == 0 {
				name = prefix + "%_"
			}
			if err := queries.CreateUser(ctx, model.CreateUserParams{UserID: uid.UUID(), CreatedAt: time.Now().Add(time.Duration(index) * time.Second), DisplayName: sql.NullString{String: name, Valid: true}, DisplayID: sql.NullString{String: prefix + string(rune('a'+index)), Valid: true}, Email: sql.NullString{String: cipher, Valid: true}, EmailVerified: true}); err != nil {
				return err
			}
		}
		for _, provider := range []string{"GOOGLE", "LINE", "GOOGLE"} {
			if err := queries.CreateUserAuth(ctx, model.CreateUserAuthParams{UserAuthID: shared.NewUUID[user.User]().UUID(), UserID: id.UUID(), Provider: provider, Subject: prefix + provider, CreatedAt: time.Now()}); err != nil {
				return err
			}
		}
		out, err := q.List(ctx, userquery.ManagementListInput{Search: prefix, Limit: 1})
		if err := queries.CreatePasswordAuth(ctx, model.CreatePasswordAuthParams{
			PasswordAuthID: shared.NewUUID[user.User]().UUID(), UserID: id.UUID(),
			PasswordHash: "test-only-unused-password-hash", LastChanged: time.Now(),
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}); err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if out.TotalCount != 2 || len(out.Users) != 1 || out.Users[0].UserID != second {
			t.Fatal("pagination or stable ordering incorrect")
		}
		out, err = q.List(ctx, userquery.ManagementListInput{Search: prefix, Offset: 1, Limit: 1})
		if err != nil {
			return err
		}
		if len(out.Users) != 1 || out.Users[0].UserID != id || len(out.Users[0].Providers) != 3 {
			t.Fatal("providers duplicated users")
		}
		if out.Users[0].Email == nil || *out.Users[0].Email != "city@example.com" {
			t.Fatal("encrypted email not decrypted")
		}
		out, err = q.List(ctx, userquery.ManagementListInput{Search: prefix + "%_", Limit: 20})
		if err != nil {
			return err
		}
		if out.TotalCount != 1 {
			t.Fatal("search wildcard was not literal")
		}
		out, err = q.List(ctx, userquery.ManagementListInput{Search: prefix, Offset: 50, Limit: 20})
		if err != nil {
			return err
		}
		if len(out.Users) != 0 || out.TotalCount != 2 {
			t.Fatal("empty page lost total")
		}
		missing, err := q.FindByID(ctx, shared.NewUUID[user.User]())
		if err != nil {
			return err
		}
		if missing != nil {
			t.Fatal("missing user was not nil")
		}
		found, err := q.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if len(found.Providers) != 3 {
			t.Fatal("detail auth methods incorrect")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
