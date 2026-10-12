package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	userquery "github.com/neko-dream/api/internal/infrastructure/persistence/query/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/migrations"
)

func permissionTestDB(t *testing.T) (*sql.DB, *db.DBManager) {
	t.Helper()
	url := os.Getenv("ADMIN_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("ADMIN_TEST_DATABASE_URL required")
	}
	c, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c, db.NewDBManager(c)
}

func TestSessionCreationPermissionDatabase(t *testing.T) {
	_, manager := permissionTestDB(t)
	if err := manager.TestTx(context.Background(), func(ctx context.Context) error {
		id := shared.NewUUID[user.User]()
		q := manager.GetQueries(ctx)
		if err := q.CreateUser(ctx, model.CreateUserParams{UserID: id.UUID(), CreatedAt: time.Now(), Email: sql.NullString{String: "invalid-ciphertext", Valid: true}}); err != nil {
			return err
		}
		if err := q.CreateUserAuth(ctx, model.CreateUserAuthParams{UserAuthID: shared.NewUUID[user.User]().UUID(), UserID: id.UUID(), Provider: "GOOGLE", Subject: id.String(), CreatedAt: time.Now()}); err != nil {
			return err
		}
		r := repository.NewSessionCreationPermissionRepository(manager)
		current := userquery.NewCurrentQuery(manager)
		initial, err := r.Get(ctx, id)
		if err != nil {
			return err
		}
		if initial {
			return errors.New("new user granted by default")
		}
		for _, value := range []bool{false, true, false} {
			if err := r.Set(ctx, id, value); err != nil {
				return err
			}
			got, err := r.Get(ctx, id)
			if err != nil {
				return err
			}
			if got != value {
				return errors.New("grant/revoke mismatch")
			}
			me, err := current.FindByID(ctx, id)
			if err != nil {
				return err
			}
			if me == nil || me.CanCreateTalkSession != value {
				return errors.New("current-user permission stale or private email dependency")
			}
		}
		if _, err := r.Get(ctx, shared.NewUUID[user.User]()); !errors.Is(err, sql.ErrNoRows) {
			return errors.New("missing user not distinguished")
		}
		me, err := current.FindByID(ctx, shared.NewUUID[user.User]())
		if err != nil {
			return err
		}
		if me != nil {
			return errors.New("fabricated missing current user")
		}
		guestID := shared.NewUUID[user.User]()
		if err := q.CreateUser(ctx, model.CreateUserParams{UserID: guestID.UUID(), CreatedAt: time.Now()}); err != nil {
			return err
		}
		if err := q.CreateUserAuth(ctx, model.CreateUserAuthParams{UserAuthID: shared.NewUUID[user.User]().UUID(), UserID: guestID.UUID(), Provider: "GUEST", Subject: guestID.String(), CreatedAt: time.Now()}); err != nil {
			return err
		}
		if err := r.Set(ctx, guestID, true); !errors.Is(err, sql.ErrNoRows) {
			return errors.New("storage accepted guest grant")
		}
		if err := r.Set(ctx, guestID, false); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestSessionCreationPermissionMigration(t *testing.T) {
	c, _ := permissionTestDB(t)
	tx, err := c.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	down, err := migrations.MigrationFiles.ReadFile("000058_user_session_creation_permission.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(string(down)); err != nil {
		t.Fatal(err)
	}
	ids := make([]shared.UUID[user.User], 5)
	for i := range ids {
		ids[i] = shared.NewUUID[user.User]()
		var withdrawn interface{}
		if i == 2 {
			withdrawn = time.Now()
		}
		if _, err := tx.Exec("INSERT INTO users(user_id,withdrawal_date) VALUES($1,$2)", ids[i].UUID(), withdrawn); err != nil {
			t.Fatal(err)
		}
		provider := "GOOGLE"
		if i == 3 {
			provider = "GUEST"
		}
		if _, err := tx.Exec("INSERT INTO user_auths(user_auth_id,user_id,provider,subject,is_verified) VALUES($1,$2,$3,$4,false)", shared.NewUUID[user.User]().UUID(), ids[i].UUID(), provider, ids[i].String()); err != nil {
			t.Fatal(err)
		}
		if i != 1 && i != 4 {
			if _, err := tx.Exec("INSERT INTO organization_users(organization_user_id,user_id,organization_id) VALUES($1,$2,$3)", shared.NewUUID[user.User]().UUID(), ids[i].UUID(), shared.NewUUID[user.User]().UUID()); err != nil {
				t.Fatal(err)
			}
		}
	}
	up, err := migrations.MigrationFiles.ReadFile("000058_user_session_creation_permission.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(string(up)); err != nil {
		t.Fatal(err)
	}
	for i, id := range ids {
		var value bool
		if err := tx.QueryRow("SELECT can_create_talk_session FROM users WHERE user_id=$1", id.UUID()).Scan(&value); err != nil {
			t.Fatal(err)
		}
		if value != (i == 0) {
			t.Fatalf("migration case %d value=%v", i, value)
		}
	}
	// Joining after migration must not implicitly grant creation access.
	if _, err := tx.Exec("INSERT INTO organization_users(organization_user_id,user_id,organization_id) VALUES($1,$2,$3)", shared.NewUUID[user.User]().UUID(), ids[4].UUID(), shared.NewUUID[user.User]().UUID()); err != nil {
		t.Fatal(err)
	}
	var value bool
	if err := tx.QueryRow("SELECT can_create_talk_session FROM users WHERE user_id=$1", ids[4].UUID()).Scan(&value); err != nil || value {
		t.Fatal("new membership granted permission", err)
	}
}
