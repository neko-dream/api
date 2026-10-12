package auth_usecase

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/auth"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	tokens "github.com/neko-dream/api/internal/infrastructure/auth/session"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/internal/test/integration"
)

type organizationOAuthAuthentication struct {
	service.AuthenticationService
	usr *user.User
}

func (a *organizationOAuthAuthentication) Authenticate(context.Context, string, string) (*user.User, error) {
	return a.usr, nil
}

func TestOAuthAutomaticallyUsesPersistedOrganizationIntegration(t *testing.T) {
	conn := integration.NewDatabase(t)
	manager := db.NewDBManager(conn)
	cfg := &config.Config{Env: config.LOCAL, TokenSecret: "org-test-secret", ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
	encryptor, err := crypto.NewEncryptor(cfg)
	if err != nil {
		t.Fatal(err)
	}
	users := repository.NewUserRepository(manager, nil, encryptor)
	sessions := repository.NewSessionRepository(manager)
	orgs := repository.NewOrganizationRepository(manager)
	members := repository.NewOrganizationUserRepository(manager)
	states := repository.NewAuthStateRepository(manager)
	tm := tokens.NewSessionTokenManager(cfg, manager, sessions, users, members, orgs)
	for _, provider := range []string{"google", "line"} {
		for _, scenario := range []string{"member", "ordinary", "wrong requested organization", "code without membership", "removed membership"} {
			t.Run(provider+"/"+scenario, func(t *testing.T) {
				ctx := context.Background()
				usr := user.NewUser(shared.NewUUID[user.User](), nil, nil, "oauth-"+shared.NewUUID[any]().String(), shared.AuthProviderName(provider), nil)
				if err := users.Create(ctx, usr); err != nil {
					t.Fatal(err)
				}
				oid := shared.NewUUID[organization.Organization]()
				org := organization.NewOrganization(oid, organization.OrganizationTypeGovernment, "City "+oid.String(), "city-"+oid.String(), nil, usr.UserID())
				if err := orgs.Create(ctx, org); err != nil {
					t.Fatal(err)
				}
				assigned := scenario == "member" || scenario == "wrong requested organization" || scenario == "removed membership"
				if assigned {
					membership := organization.NewOrganizationUser(shared.NewUUID[organization.OrganizationUser](), oid, usr.UserID(), organization.OrganizationUserRoleAdmin)
					if err := members.Create(ctx, *membership); err != nil {
						t.Fatal(err)
					}
				}
				var requested *shared.UUID[any]
				if scenario == "wrong requested organization" {
					id := shared.NewUUID[any]()
					requested = &id
				}
				if scenario == "code without membership" {
					id := shared.UUID[any](oid.UUID())
					requested = &id
				}
				stateValue := shared.NewUUID[any]().String()
				state := auth.NewState(stateValue, provider, "https://localhost:3000/home", time.Now().Add(time.Minute), nil, requested)
				if err := states.Create(ctx, state); err != nil {
					t.Fatal(err)
				}
				callback := NewAuthCallback(manager, cfg, &organizationOAuthAuthentication{usr: &usr}, sessions, service.NewSessionService(sessions), tm, states, orgs, members)
				output, err := callback.Execute(ctx, CallbackInput{Provider: provider, Code: "test-code", State: stateValue})
				if strings.Contains(scenario, "requested") || scenario == "code without membership" {
					if !errors.Is(err, messages.OrganizationPermissionDenied) {
						t.Fatalf("untrusted organization accepted: %v", err)
					}
					active, err := sessions.FindByUserID(ctx, usr.UserID())
					if err != nil || len(active) != 0 {
						t.Fatal("rejected OAuth created session", err)
					}
					if _, err := states.Get(ctx, stateValue); err != nil {
						t.Fatal("failed OAuth did not rollback state", err)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				claim, err := tm.Parse(ctx, output.Token)
				if err != nil {
					t.Fatal(err)
				}
				if assigned {
					if claim.OrganizationID == nil || *claim.OrganizationID != oid.String() || claim.OrganizationRole == nil || *claim.OrganizationRole != organization.RoleToName(organization.OrganizationUserRoleAdmin) {
						t.Fatalf("organization missing: %#v", claim)
					}
				} else if claim.OrganizationID != nil {
					t.Fatal("ordinary account got organization access")
				}
				if scenario == "removed membership" {
					if _, err := conn.Exec("DELETE FROM organization_users WHERE user_id=$1", usr.UserID().UUID()); err != nil {
						t.Fatal(err)
					}
					if _, err := tm.Parse(ctx, output.Token); err == nil {
						t.Fatal("removed membership retained organization session")
					}
				}
			})
		}
	}
}

func TestSingleMembershipConstraintConcurrencyAndMigrationIntegration(t *testing.T) {
	conn := integration.NewDatabase(t)
	manager := db.NewDBManager(conn)
	members := repository.NewOrganizationUserRepository(manager)
	uid := shared.NewUUID[user.User]()
	var wg sync.WaitGroup
	results := make(chan error, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			membership := organization.NewOrganizationUser(shared.NewUUID[organization.OrganizationUser](), shared.NewUUID[organization.Organization](), uid, organization.OrganizationUserRoleMember)
			results <- members.Create(context.Background(), *membership)
		}()
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		} else if !errors.Is(err, messages.OrganizationMembershipConflict) {
			t.Fatalf("unexpected error %T: %#v", err, err)
		}
	}
	if success != 1 {
		t.Fatalf("parallel creates succeeded %d times", success)
	}
	rows, err := members.FindByUserID(context.Background(), uid)
	if err != nil || len(rows) != 1 {
		t.Fatal("constraint did not persist one membership", err)
	}
	// Simulate pre-migration legacy duplicates only inside this isolated schema.
	if _, err := conn.Exec("ALTER TABLE organization_users DROP CONSTRAINT organization_users_one_per_user"); err != nil {
		t.Fatal(err)
	}
	duplicate := organization.NewOrganizationUser(shared.NewUUID[organization.OrganizationUser](), shared.NewUUID[organization.Organization](), uid, organization.OrganizationUserRoleMember)
	if err := members.Create(context.Background(), *duplicate); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../../migrations/000056_single_organization_membership.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(string(migration)); err == nil || !strings.Contains(err.Error(), "resolve explicitly") {
		t.Fatal("duplicate migration did not fail safely", err)
	}
	rows, err = members.FindByUserID(context.Background(), uid)
	if err != nil || len(rows) != 2 {
		t.Fatal("migration silently deleted legacy memberships", err)
	}
}
