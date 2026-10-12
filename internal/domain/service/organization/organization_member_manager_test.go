package organization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/neko-dream/api/internal/domain/messages"
	mock_consent "github.com/neko-dream/api/internal/domain/mock/consent"
	mock_user "github.com/neko-dream/api/internal/domain/mock/user"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/consent"
	mock_org "github.com/neko-dream/api/internal/domain/model/mock/organization"
	org "github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"go.uber.org/mock/gomock"
)

type testTransaction struct {
	calls               int
	result              error
	beginErr, commitErr error
}

func (tx *testTransaction) ExecTx(ctx context.Context, fn func(context.Context) error) error {
	tx.calls++
	if tx.beginErr != nil {
		return tx.beginErr
	}
	tx.result = fn(ctx)
	if tx.result != nil {
		return tx.result
	}
	return tx.commitErr
}

type testPasswordManager struct {
	registered, updated int
	password            string
	requiredChange      bool
	err                 error
}

func (p *testPasswordManager) RegisterPassword(_ context.Context, _ shared.UUID[user.User], value string, required bool) error {
	p.registered++
	p.password = value
	p.requiredChange = required
	return p.err
}
func (p *testPasswordManager) UpdatePassword(context.Context, shared.UUID[user.User], string) error {
	p.updated++
	return p.err
}
func (p *testPasswordManager) VerifyPassword(context.Context, shared.UUID[user.User], string) (bool, error) {
	return false, nil
}

var _ password_auth.PasswordAuthManager = (*testPasswordManager)(nil)

type testEmailSender struct {
	cfg   *config.Config
	calls int
	data  map[string]any
	to    string
	err   error
}

func (s *testEmailSender) Send(_ context.Context, to string, tmpl email_template.EmailTemplateType, data map[string]any) error {
	s.calls++
	s.data = data
	s.to = to
	if _, err := email.RenderMessage(s.cfg, to, tmpl, data); err != nil {
		return err
	}
	return s.err
}

type invitationFixture struct {
	manager      *organizationMemberManager
	users        *mock_user.MockUserRepository
	orgs         *mock_org.MockOrganizationRepository
	members      *mock_org.MockOrganizationUserRepository
	policy       *mock_consent.MockPolicyRepository
	consent      *mock_consent.MockConsentService
	password     *testPasswordManager
	email        *testEmailSender
	tx           *testTransaction
	input        InviteUserParams
	organization *org.Organization
}

func newInvitationFixture(t *testing.T) *invitationFixture {
	t.Helper()
	ctrl := gomock.NewController(t)
	f := &invitationFixture{users: mock_user.NewMockUserRepository(ctrl), orgs: mock_org.NewMockOrganizationRepository(ctrl), members: mock_org.NewMockOrganizationUserRepository(ctrl), policy: mock_consent.NewMockPolicyRepository(ctrl), consent: mock_consent.NewMockConsentService(ctrl), password: &testPasswordManager{}, tx: &testTransaction{}}
	cfg := &config.Config{HASH_PEPPER: "test-pepper", EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://example.com"}
	f.email = &testEmailSender{cfg: cfg}
	f.input = InviteUserParams{UserID: shared.NewUUID[user.User](), OrganizationID: shared.NewUUID[org.Organization](), Email: "invitee@example.com", Role: org.OrganizationUserRoleMember}
	f.organization = &org.Organization{OrganizationID: f.input.OrganizationID, Name: "テスト市役所"}
	f.manager = &organizationMemberManager{organizationRepo: f.orgs, organizationUserRepo: f.members, cfg: cfg, userRep: f.users, policyRep: f.policy, consentService: f.consent, passwordAuthManager: f.password, emailSender: f.email, txRunner: f.tx}
	return f
}

func (f *invitationFixture) allowInvitation() {
	f.members.EXPECT().FindByUserID(gomock.Any(), f.input.UserID).Return([]*org.OrganizationUser{{Role: org.OrganizationUserRoleSuperAdmin}}, nil)
	f.orgs.EXPECT().FindByID(gomock.Any(), f.input.OrganizationID).Return(f.organization, nil)
}

func TestAddUserPreservesMembershipConflict(t *testing.T) {
	for _, tc := range []struct {
		name      string
		createErr error
		wantErr   error
	}{
		{"membership conflict", fmt.Errorf("wrapped: %w", messages.OrganizationMembershipConflict), messages.OrganizationMembershipConflict},
		{"database failure", errors.New("database unavailable"), messages.OrganizationInternalServerError},
		{"success", nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newInvitationFixture(t)
			f.members.EXPECT().FindByOrganizationIDAndUserID(gomock.Any(), f.input.OrganizationID, f.input.UserID).Return(nil, sql.ErrNoRows)
			f.orgs.EXPECT().FindByID(gomock.Any(), f.input.OrganizationID).Return(f.organization, nil)
			f.members.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tc.createErr)
			if err := f.manager.AddUser(context.Background(), f.input); !errors.Is(err, tc.wantErr) {
				t.Fatalf("error %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestAddUserRejectsInvalidRoleBeforeRepositoryAccess(t *testing.T) {
	f := newInvitationFixture(t)
	f.input.Role = 11
	if err := f.manager.AddUser(context.Background(), f.input); !errors.Is(err, messages.OrganizationPermissionDenied) {
		t.Fatalf("invalid role accepted: %v", err)
	}
}

func TestIsSuperAdminRejectsInvalidRole(t *testing.T) {
	f := newInvitationFixture(t)
	f.members.EXPECT().FindByUserID(gomock.Any(), f.input.UserID).Return([]*org.OrganizationUser{{Role: 0}}, nil)
	if ok, err := f.manager.IsSuperAdmin(context.Background(), f.input.UserID); err != nil || ok {
		t.Fatalf("invalid role granted operator access: %v, %v", ok, err)
	}
}

func TestInviteExistingUser(t *testing.T) {
	for _, alreadyMember := range []bool{false, true} {
		for _, sendFailure := range []bool{false, true} {
			t.Run(fmt.Sprintf("alreadyMember=%t/sendFailure=%t", alreadyMember, sendFailure), func(t *testing.T) {
				f := newInvitationFixture(t)
				f.allowInvitation()
				usr := user.NewUser(shared.NewUUID[user.User](), nil, nil, "existing", "PASSWORD", nil)
				f.users.EXPECT().FindBySubject(gomock.Any(), gomock.Any()).Return(&usr, nil)
				membership := org.NewOrganizationUser(shared.NewUUID[org.OrganizationUser](), f.input.OrganizationID, usr.UserID(), org.OrganizationUserRoleMember)
				if alreadyMember {
					f.input.Role = org.OrganizationUserRoleAdmin // Re-invitation must not promote an existing member.
					f.members.EXPECT().FindByOrganizationIDAndUserID(gomock.Any(), f.input.OrganizationID, usr.UserID()).Return(membership, nil)
				} else {
					f.members.EXPECT().FindByOrganizationIDAndUserID(gomock.Any(), f.input.OrganizationID, usr.UserID()).Return(nil, sql.ErrNoRows)
					f.members.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, m org.OrganizationUser) error {
						if m.UserID != usr.UserID() || m.OrganizationID != f.input.OrganizationID || m.Role != f.input.Role {
							t.Fatal("incorrect membership")
						}
						return nil
					})
				}
				sendErr := errors.New("delivery unavailable")
				if sendFailure {
					f.email.err = sendErr
				}
				result, err := f.manager.InviteUser(context.Background(), f.input)
				if sendFailure {
					if !errors.Is(err, sendErr) || result != nil || !errors.Is(f.tx.result, sendErr) {
						t.Fatalf("send failure did not fail the transaction: %v", err)
					}
				} else if err != nil || result == nil || result.UserID != usr.UserID() || result.Role != org.OrganizationUserRoleMember {
					t.Fatalf("incorrect invitation result: %#v %v", result, err)
				}
				if f.password.registered != 0 || f.password.updated != 0 {
					t.Fatal("existing user's password changed")
				}
				if f.email.calls != 1 || f.email.data["ExistingAccount"] != true || f.email.data["Password"] != "" {
					t.Fatal("existing account email includes replacement credentials")
				}
			})
		}
	}
}

func TestInviteNewUser(t *testing.T) {
	for _, sendFailure := range []bool{false, true} {
		t.Run(fmt.Sprintf("sendFailure=%t", sendFailure), func(t *testing.T) {
			f := newInvitationFixture(t)
			f.allowInvitation()
			f.users.EXPECT().FindBySubject(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
			f.policy.EXPECT().FetchLatestPolicy(gomock.Any()).Return(&consent.Policy{Version: "1.0"}, nil)
			f.consent.EXPECT().RecordConsent(gomock.Any(), gomock.Any(), "1.0", "", "").Return(nil, nil)
			var createdID shared.UUID[user.User]
			f.users.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, u user.User) error { createdID = u.UserID(); return nil })
			f.members.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, m org.OrganizationUser) error {
				if m.UserID != createdID || m.OrganizationID != f.input.OrganizationID {
					t.Fatal("incorrect new membership")
				}
				return nil
			})
			sendErr := errors.New("delivery failed")
			if sendFailure {
				f.email.err = sendErr
			}
			result, err := f.manager.InviteUser(context.Background(), f.input)
			if sendFailure {
				if !errors.Is(err, sendErr) || result != nil || !errors.Is(f.tx.result, sendErr) {
					t.Fatalf("send failure not propagated: %v", err)
				}
			} else if err != nil || result == nil || result.UserID != createdID {
				t.Fatalf("incorrect result: %#v %v", result, err)
			}
			if f.password.registered != 1 || f.password.updated != 0 || !f.password.requiredChange || len(f.password.password) != 16 {
				t.Fatal("initial password must be registered with a change requirement")
			}
			if f.email.data["Password"] != f.password.password || f.email.data["ExistingAccount"] != false {
				t.Fatal("initial password mail differs from registered credentials")
			}
		})
	}
}

func TestInviteRejectsInvalidInput(t *testing.T) {
	for _, address := range []string{"", "invalid", "a@example.com,b@example.com", "a@example.com\r\nBcc:b@example.com", "Name <a@example.com>"} {
		t.Run(address, func(t *testing.T) {
			f := newInvitationFixture(t)
			f.input.Email = address
			if _, err := f.manager.InviteUser(context.Background(), f.input); err == nil {
				t.Fatal("invalid email accepted")
			}
			if f.tx.calls != 0 || f.email.calls != 0 {
				t.Fatal("invalid input caused side effects")
			}
		})
	}
	for _, role := range []org.OrganizationUserRole{0, 11, 99} {
		t.Run(fmt.Sprintf("role=%d", role), func(t *testing.T) {
			f := newInvitationFixture(t)
			f.input.Role = role
			if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, messages.OrganizationPermissionDenied) {
				t.Fatalf("invalid role: %v", err)
			}
		})
	}
}

func TestInvitePermissionAndRepositoryFailures(t *testing.T) {
	t.Run("non-super-admin", func(t *testing.T) {
		f := newInvitationFixture(t)
		f.members.EXPECT().FindByUserID(gomock.Any(), f.input.UserID).Return([]*org.OrganizationUser{{Role: org.OrganizationUserRoleAdmin}}, nil)
		if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, messages.OrganizationForbidden) {
			t.Fatal(err)
		}
		if f.tx.calls != 0 {
			t.Fatal("unauthorized transaction started")
		}
	})
	t.Run("missing organization", func(t *testing.T) {
		f := newInvitationFixture(t)
		f.members.EXPECT().FindByUserID(gomock.Any(), f.input.UserID).Return([]*org.OrganizationUser{{Role: org.OrganizationUserRoleSuperAdmin}}, nil)
		f.orgs.EXPECT().FindByID(gomock.Any(), f.input.OrganizationID).Return(nil, nil)
		if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, messages.OrganizationForbidden) {
			t.Fatal(err)
		}
	})
	t.Run("membership write failure", func(t *testing.T) {
		f := newInvitationFixture(t)
		f.allowInvitation()
		usr := user.NewUser(shared.NewUUID[user.User](), nil, nil, "existing", "PASSWORD", nil)
		f.users.EXPECT().FindBySubject(gomock.Any(), gomock.Any()).Return(&usr, nil)
		f.members.EXPECT().FindByOrganizationIDAndUserID(gomock.Any(), f.input.OrganizationID, usr.UserID()).Return(nil, sql.ErrNoRows)
		repoErr := errors.New("membership write failed")
		f.members.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)
		if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, repoErr) {
			t.Fatal(err)
		}
		if f.email.calls != 0 {
			t.Fatal("email sent despite write failure")
		}
	})
	t.Run("lookup failure", func(t *testing.T) {
		f := newInvitationFixture(t)
		f.allowInvitation()
		repoErr := errors.New("lookup failed")
		f.users.EXPECT().FindBySubject(gomock.Any(), gomock.Any()).Return(nil, repoErr)
		if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, repoErr) {
			t.Fatal(err)
		}
		if f.email.calls != 0 || f.password.registered != 0 {
			t.Fatal("lookup failure caused side effects")
		}
	})
}

func TestInviteNewUserFailures(t *testing.T) {
	for _, stage := range []string{"policy", "consent", "user", "password", "membership"} {
		t.Run(stage, func(t *testing.T) {
			f := newInvitationFixture(t)
			f.allowInvitation()
			f.users.EXPECT().FindBySubject(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
			failure := errors.New(stage + " failure")
			policyErr := error(nil)
			if stage == "policy" {
				policyErr = failure
			}
			f.policy.EXPECT().FetchLatestPolicy(gomock.Any()).Return(&consent.Policy{Version: "1.0"}, policyErr)
			if stage != "policy" {
				consentErr := error(nil)
				if stage == "consent" {
					consentErr = failure
				}
				f.consent.EXPECT().RecordConsent(gomock.Any(), gomock.Any(), "1.0", "", "").Return(nil, consentErr)
				if stage != "consent" {
					userErr := error(nil)
					if stage == "user" {
						userErr = failure
					}
					f.users.EXPECT().Create(gomock.Any(), gomock.Any()).Return(userErr)
					if stage == "password" {
						f.password.err = failure
					}
					if stage == "membership" {
						f.members.EXPECT().Create(gomock.Any(), gomock.Any()).Return(failure)
					}
				}
			}
			result, err := f.manager.InviteUser(context.Background(), f.input)
			if !errors.Is(err, failure) || result != nil || !errors.Is(f.tx.result, failure) {
				t.Fatalf("failure lost: %#v %v", result, err)
			}
			if f.email.calls != 0 {
				t.Fatal("email sent despite failed preparation")
			}
		})
	}
}

func TestInviteActorLookupFailure(t *testing.T) {
	f := newInvitationFixture(t)
	failure := errors.New("actor lookup failed")
	f.members.EXPECT().FindByUserID(gomock.Any(), f.input.UserID).Return(nil, failure)
	if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if f.tx.calls != 0 || f.email.calls != 0 {
		t.Fatal("failed authorization caused side effects")
	}
}

func TestInviteTransactionBeginFailure(t *testing.T) {
	f := newInvitationFixture(t)
	failure := errors.New("cannot start transaction")
	f.tx.beginErr = failure
	f.members.EXPECT().FindByUserID(gomock.Any(), f.input.UserID).Return([]*org.OrganizationUser{{Role: org.OrganizationUserRoleSuperAdmin}}, nil)
	if _, err := f.manager.InviteUser(context.Background(), f.input); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if f.email.calls != 0 {
		t.Fatal("email sent without a transaction")
	}
}
