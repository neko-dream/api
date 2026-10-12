// Creates a fixture only in the dedicated loopback local-development database.
// No .env, AWS credentials or production database configuration are loaded.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/pkg/hash"
)

func main() {
	if os.Getenv("ENV") != "local" {
		log.Fatal("set ENV=local; this command is only for local test fixtures")
	}
	if len(os.Args) != 3 && len(os.Args) != 4 {
		log.Fatal("usage: ENV=local go run ./cmd/local-email-account EMAIL LOCAL_TEST_PASSWORD [NEW_ORGANIZATION_CODE]")
	}
	address, password := os.Args[1], os.Args[2]
	orgCode := ""
	if len(os.Args) == 4 {
		orgCode = os.Args[3]
		if err := organization.ValidateOrganizationCode(orgCode); err != nil {
			log.Fatal(err)
		}
	}
	parsed, err := mail.ParseAddress(address)
	if err != nil || parsed.Address != address || strings.ContainsAny(address, "\r\n") || len(password) < 12 {
		log.Fatal("provide a bare email address and a local-only password of at least 12 characters")
	}
	// These values deliberately match scripts/dev-local.sh; arbitrary DB URLs
	// are never accepted so invoking this command cannot mutate a remote DB.
	cfg := &config.Config{Env: config.LOCAL, HASH_PEPPER: "local-development-only-password-pepper", HASH_ITERATIONS: 10, ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
	connection, err := sql.Open("pgx", "postgres://kotohiro:kotohiro@127.0.0.1:5433/kotohiro?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	manager := db.NewDBManager(connection)
	encryptor, err := crypto.NewEncryptor(cfg)
	if err != nil {
		log.Fatal(err)
	}
	users := repository.NewUserRepository(manager, nil, encryptor)
	organizations := repository.NewOrganizationRepository(manager)
	members := repository.NewOrganizationUserRepository(manager)
	passwords := service.NewPasswordAuthManager(repository.NewPasswordAuthRepository(manager), cfg)
	subject, err := hash.HMACHash(address, cfg.HASH_PEPPER)
	if err != nil {
		log.Fatal(err)
	}
	err = manager.ExecTx(context.Background(), func(ctx context.Context) error {
		existing, err := users.FindBySubject(ctx, user.UserSubject(subject))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if existing != nil {
			return fmt.Errorf("account already exists; credentials were not changed")
		}
		id := shared.NewUUID[user.User]()
		displayID := "localmail_" + strings.ReplaceAll(id.String(), "-", "")[:12]
		name := "ローカルメール認証テスト"
		usr := user.NewUser(id, &displayID, &name, subject, "PASSWORD", nil)
		usr.ChangeEmail(address)
		if err := users.Create(ctx, usr); err != nil {
			return err
		}
		if err := passwords.RegisterPassword(ctx, id, password, false); err != nil {
			return err
		}
		if orgCode == "" {
			return nil
		}
		existingOrg, err := organizations.FindByCode(ctx, orgCode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if existingOrg != nil {
			return fmt.Errorf("organization already exists; no existing data was changed")
		}
		org := organization.NewOrganization(shared.NewUUID[organization.Organization](), organization.OrganizationTypeGovernment, "ローカル検証市", orgCode, nil, id)
		if err := organizations.Create(ctx, org); err != nil {
			return err
		}
		return members.Create(ctx, *organization.NewOrganizationUser(shared.NewUUID[organization.OrganizationUser](), org.OrganizationID, id, organization.OrganizationUserRoleOwner))
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created local-only fixture: %s\nOpen https://localhost:3000/organization/login\n", address)
	if orgCode != "" {
		fmt.Printf("Organization: ローカル検証市 (%s), role: owner\n", orgCode)
	}
}
