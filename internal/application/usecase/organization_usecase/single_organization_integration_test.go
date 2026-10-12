package organization_usecase_test

import (
	"context"
	"errors"
	"testing"

	organization_usecase "github.com/neko-dream/api/internal/application/usecase/organization_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	orgservice "github.com/neko-dream/api/internal/domain/service/organization"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/internal/test/integration"
	"github.com/neko-dream/api/pkg/hash"
)

type invitationTestSender struct {
	cfg   *config.Config
	calls int
}

func (s *invitationTestSender) Send(_ context.Context, to string, tmpl email_template.EmailTemplateType, data map[string]any) error {
	s.calls++
	_, err := email.RenderMessage(s.cfg, to, tmpl, data)
	return err
}

func TestSingleOrganizationInvitationAndCreationRollbackIntegration(t *testing.T) {
	conn := integration.NewDatabase(t)
	manager := db.NewDBManager(conn)
	cfg := &config.Config{Env: config.LOCAL, HASH_PEPPER: "org-invite-pepper", HASH_ITERATIONS: 4, EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://localhost:3000", ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
	encryptor, err := crypto.NewEncryptor(cfg)
	if err != nil {
		t.Fatal(err)
	}
	users := repository.NewUserRepository(manager, nil, encryptor)
	orgs := repository.NewOrganizationRepository(manager)
	members := repository.NewOrganizationUserRepository(manager)
	passwords := repository.NewPasswordAuthRepository(manager)
	passwordManager := service.NewPasswordAuthManager(passwords, cfg)
	ctx := context.Background()
	actor := user.NewUser(shared.NewUUID[user.User](), nil, nil, "operator-"+shared.NewUUID[any]().String(), "PASSWORD", nil)
	if err := users.Create(ctx, actor); err != nil {
		t.Fatal(err)
	}
	orgA := organization.NewOrganization(shared.NewUUID[organization.Organization](), organization.OrganizationTypeGovernment, "City A", "city-a", nil, actor.UserID())
	orgB := organization.NewOrganization(shared.NewUUID[organization.Organization](), organization.OrganizationTypeGovernment, "City B", "city-b", nil, actor.UserID())
	for _, org := range []*organization.Organization{orgA, orgB} {
		if err := orgs.Create(ctx, org); err != nil {
			t.Fatal(err)
		}
	}
	if err := members.Create(ctx, *organization.NewOrganizationUser(shared.NewUUID[organization.OrganizationUser](), orgA.OrganizationID, actor.UserID(), organization.OrganizationUserRoleSuperAdmin)); err != nil {
		t.Fatal(err)
	}
	address := "provisioned@example.com"
	subject, err := hash.HMACHash(address, cfg.HASH_PEPPER)
	if err != nil {
		t.Fatal(err)
	}
	usr := user.NewUser(shared.NewUUID[user.User](), nil, nil, subject, "PASSWORD", nil)
	usr.ChangeEmail(address)
	if err := users.Create(ctx, usr); err != nil {
		t.Fatal(err)
	}
	if err := passwordManager.RegisterPassword(ctx, usr.UserID(), "original-password", true); err != nil {
		t.Fatal(err)
	}
	if err := members.Create(ctx, *organization.NewOrganizationUser(shared.NewUUID[organization.OrganizationUser](), orgA.OrganizationID, usr.UserID(), organization.OrganizationUserRoleAdmin)); err != nil {
		t.Fatal(err)
	}
	sender := &invitationTestSender{cfg: cfg}
	memberManager := orgservice.NewOrganizationMemberManager(orgs, members, cfg, users, nil, nil, passwordManager, sender, manager)
	// Cross-organization invitations must not send mail or change credentials/role.
	input := orgservice.InviteUserParams{UserID: actor.UserID(), OrganizationID: orgB.OrganizationID, Email: address, Role: organization.OrganizationUserRoleMember}
	if result, err := memberManager.InviteUser(ctx, input); result != nil || !errors.Is(err, messages.OrganizationMembershipConflict) {
		t.Fatalf("cross-organization invite: %#v %v", result, err)
	}
	if sender.calls != 0 {
		t.Fatal("rejected cross-organization invite sent email")
	}
	// The member manager used by display-ID invitations preserves the conflict.
	if err := memberManager.AddUser(ctx, orgservice.InviteUserParams{UserID: usr.UserID(), OrganizationID: orgB.OrganizationID, Role: organization.OrganizationUserRoleMember}); !errors.Is(err, messages.OrganizationMembershipConflict) {
		t.Fatalf("cross-organization AddUser: %v", err)
	}
	rows, err := members.FindByUserID(ctx, usr.UserID())
	if err != nil || len(rows) != 1 || rows[0].OrganizationID != orgA.OrganizationID || rows[0].Role != organization.OrganizationUserRoleAdmin {
		t.Fatal("invite changed original membership", err)
	}
	if ok, err := passwordManager.VerifyPassword(ctx, usr.UserID(), "original-password"); err != nil || !ok {
		t.Fatal("invite changed original password", err)
	}
	// Same-organization re-invites retain the existing role and password.
	input.OrganizationID = orgA.OrganizationID
	result, err := memberManager.InviteUser(ctx, input)
	if err != nil || result.Role != organization.OrganizationUserRoleAdmin || sender.calls != 1 {
		t.Fatal("same-organization re-invite failed", err)
	}
	if ok, err := passwordManager.VerifyPassword(ctx, usr.UserID(), "original-password"); err != nil || !ok {
		t.Fatal("re-invite changed password", err)
	}
	// Creating another organization for an already-affiliated owner rolls back
	// the organization row as well as the failed membership; no orphan remains.
	orgService := orgservice.NewOrganizationService(orgs, members, memberManager, nil, cfg)
	create := organization_usecase.NewCreateOrganizationInteractor(orgService, manager)
	if result, err := create.Execute(ctx, organization_usecase.CreateOrganizationInput{UserID: actor.UserID(), Name: "Rejected City", Code: "rejected-city", Type: int(organization.OrganizationTypeGovernment)}); result != nil || !errors.Is(err, messages.OrganizationMembershipConflict) {
		t.Fatal("second owner organization accepted", err)
	}
	var count int
	if err := conn.QueryRow("SELECT COUNT(*) FROM organizations WHERE code='rejected-city'").Scan(&count); err != nil || count != 0 {
		t.Fatal("failed creation left orphan organization", err)
	}
}
