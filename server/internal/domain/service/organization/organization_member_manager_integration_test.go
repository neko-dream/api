package organization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	mock_consent "github.com/neko-dream/api/internal/domain/mock/consent"
	"github.com/neko-dream/api/internal/domain/model/consent"
	org "github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/pkg/hash"
	"go.uber.org/mock/gomock"
)

type recordingUserRepository struct {
	user.UserRepository
	createdID shared.UUID[user.User]
}

func (r *recordingUserRepository) Create(ctx context.Context, u user.User) error {
	r.createdID = u.UserID()
	return r.UserRepository.Create(ctx, u)
}

// Use an isolated schema, not application tables. No AWS or external mail is
// involved. CI sets MAIL_TEST_DATABASE_URL after its migration step.
func TestInvitationTransactionsIntegration(t *testing.T) {
	databaseURL := os.Getenv("MAIL_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set MAIL_TEST_DATABASE_URL to run PostgreSQL transaction tests")
	}
	ctx := context.Background()
	admin, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := "mail_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA "`+schema+`"`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.ExecContext(ctx, `DROP SCHEMA "`+schema+`" CASCADE`); err != nil {
			t.Errorf("clean up isolated test schema: %v", err)
		}
	}()
	dsn, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	query := dsn.Query()
	query.Set("search_path", schema+",public")
	dsn.RawQuery = query.Encode()
	if err := db.NewMigrator(&config.Config{DatabaseURL: dsn.String()}).Up(); err != nil {
		t.Fatal(err)
	}
	connection, err := sql.Open("pgx", dsn.String())
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	manager := db.NewDBManager(connection)
	for _, existing := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing=%t", existing), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cfg := &config.Config{HASH_PEPPER: "integration-pepper", HASH_ITERATIONS: 4, EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://example.com", ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
			encryptor, err := crypto.NewEncryptor(cfg)
			if err != nil {
				t.Fatal(err)
			}
			users := &recordingUserRepository{UserRepository: repository.NewUserRepository(manager, nil, encryptor)}
			orgs := repository.NewOrganizationRepository(manager)
			members := repository.NewOrganizationUserRepository(manager)
			passwordRepo := repository.NewPasswordAuthRepository(manager)
			password := service.NewPasswordAuthManager(passwordRepo, cfg)
			actor := user.NewUser(shared.NewUUID[user.User](), nil, nil, "mail-test-actor-"+uuid.NewString(), "PASSWORD", nil)
			if err := users.Create(ctx, actor); err != nil {
				t.Fatal(err)
			}
			organization := org.NewOrganization(shared.NewUUID[org.Organization](), org.OrganizationTypeGovernment, "メールテスト市役所", "MAIL"+strings.ReplaceAll(uuid.NewString(), "-", ""), nil, actor.UserID())
			if err := orgs.Create(ctx, organization); err != nil {
				t.Fatal(err)
			}
			if err := members.Create(ctx, *org.NewOrganizationUser(shared.NewUUID[org.OrganizationUser](), organization.OrganizationID, actor.UserID(), org.OrganizationUserRoleSuperAdmin)); err != nil {
				t.Fatal(err)
			}
			address := "invite-" + uuid.NewString() + "@example.com"
			subject, err := hash.HMACHash(address, cfg.HASH_PEPPER)
			if err != nil {
				t.Fatal(err)
			}
			if existing {
				usr := user.NewUser(shared.NewUUID[user.User](), nil, nil, subject, "PASSWORD", nil)
				usr.ChangeEmail(address)
				if err := users.Create(ctx, usr); err != nil {
					t.Fatal(err)
				}
				if err := password.RegisterPassword(ctx, usr.UserID(), "original-password", true); err != nil {
					t.Fatal(err)
				}
			}
			policy := mock_consent.NewMockPolicyRepository(ctrl)
			consentService := mock_consent.NewMockConsentService(ctrl)
			if !existing {
				policy.EXPECT().FetchLatestPolicy(gomock.Any()).Return(&consent.Policy{Version: "test"}, nil).Times(2)
				consentService.EXPECT().RecordConsent(gomock.Any(), gomock.Any(), "test", "", "").Return(nil, nil).Times(2)
			}
			sendErr := errors.New("simulated email failure")
			sender := &testEmailSender{cfg: cfg, err: sendErr}
			svc := NewOrganizationMemberManager(orgs, members, cfg, users, policy, consentService, password, sender, manager)
			input := InviteUserParams{UserID: actor.UserID(), OrganizationID: organization.OrganizationID, Email: address, Role: org.OrganizationUserRoleMember}
			if result, err := svc.InviteUser(ctx, input); !errors.Is(err, sendErr) || result != nil {
				t.Fatalf("expected send failure: %#v %v", result, err)
			}
			createdID := users.createdID
			if _, err := members.FindByOrganizationIDAndUserID(ctx, organization.OrganizationID, createdID); !errors.Is(err, sql.ErrNoRows) {
				t.Fatalf("failed invitation left membership behind: %v", err)
			}
			if existing {
				ok, err := password.VerifyPassword(ctx, createdID, "original-password")
				if err != nil || !ok {
					t.Fatal("existing password changed after failed invitation")
				}
			} else {
				if _, err := users.FindBySubject(ctx, user.UserSubject(subject)); !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("failed invitation left user behind: %v", err)
				}
				if _, err := passwordRepo.GetPasswordAuthByUserID(ctx, createdID); !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("failed invitation left password behind: %v", err)
				}
			}
			sender.err = nil
			result, err := svc.InviteUser(ctx, input)
			if err != nil || result == nil {
				t.Fatalf("retry failed: %v", err)
			}
			stored, err := members.FindByOrganizationIDAndUserID(ctx, organization.OrganizationID, result.UserID)
			if err != nil || stored.Role != org.OrganizationUserRoleMember {
				t.Fatalf("retry did not persist membership: %v", err)
			}
			expectedPassword := "original-password"
			if !existing {
				expectedPassword = sender.data["Password"].(string)
			}
			ok, err := password.VerifyPassword(ctx, result.UserID, expectedPassword)
			if err != nil || !ok {
				t.Fatal("password and invitation did not agree")
			}
			before, err := passwordRepo.GetPasswordAuthByUserID(ctx, result.UserID)
			if err != nil {
				t.Fatal(err)
			}
			input.Role = org.OrganizationUserRoleAdmin
			reinvited, err := svc.InviteUser(ctx, input)
			if err != nil || reinvited.OrganizationUserID != result.OrganizationUserID || reinvited.Role != org.OrganizationUserRoleMember {
				t.Fatalf("re-invitation changed membership: %v", err)
			}
			after, err := passwordRepo.GetPasswordAuthByUserID(ctx, result.UserID)
			if err != nil || before.PasswordHash != after.PasswordHash || before.RequiredPasswordChange != after.RequiredPasswordChange {
				t.Fatal("re-invitation changed password or change requirement")
			}
			var count int
			if err := connection.QueryRowContext(ctx, "SELECT count(*) FROM organization_users WHERE organization_id=$1 AND user_id=$2", organization.OrganizationID.UUID(), result.UserID.UUID()).Scan(&count); err != nil || count != 1 {
				t.Fatalf("duplicate membership after retry: %d %v", count, err)
			}
		})
	}
}
