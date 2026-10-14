package admin_usecase

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

type invitationTx struct{}

func (invitationTx) ExecTx(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

type invitationUsers struct {
	user.UserRepository
	existing, created *user.User
	err               error
}

func (r *invitationUsers) FindBySubject(context.Context, user.UserSubject) (*user.User, error) {
	return r.existing, r.err
}
func (r *invitationUsers) FindByID(context.Context, shared.UUID[user.User]) (*user.User, error) {
	return r.created, r.err
}
func (r *invitationUsers) Create(_ context.Context, u user.User) error { r.created = &u; return r.err }

type invitationTokens struct {
	password_auth.SetupTokenRepository
	value password_auth.SetupToken
	err   error
	saves int
}

func (*invitationTokens) LockEmailSubject(context.Context, string) error         { return nil }
func (*invitationTokens) LockUser(context.Context, shared.UUID[user.User]) error { return nil }
func (r *invitationTokens) FindByUserID(context.Context, shared.UUID[user.User]) (password_auth.SetupToken, error) {
	return r.value, r.err
}
func (r *invitationTokens) Save(_ context.Context, token password_auth.SetupToken) error {
	r.saves++
	r.value = token
	return r.err
}

type invitationPasswords struct {
	password_auth.PasswordAuthRepository
	err error
}

func (r invitationPasswords) GetPasswordAuthByUserID(context.Context, shared.UUID[user.User]) (password_auth.PasswordAuth, error) {
	return password_auth.PasswordAuth{}, r.err
}

type invitationSender struct {
	calls int
	data  map[string]any
	err   error
	cfg   *config.Config
}

func (s *invitationSender) Send(_ context.Context, to string, tmpl email_template.EmailTemplateType, data map[string]any) error {
	s.calls++
	s.data = data
	if s.err != nil {
		return s.err
	}
	_, err := email.RenderMessage(s.cfg, to, tmpl, data)
	return err
}

func invitationFixture() (*manageUserInvitations, context.Context) {
	cfg := &config.Config{HASH_PEPPER: "test-pepper", WEBSITE_URL: "https://example.com", EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ"}
	m := &manageUserInvitations{users: &invitationUsers{}, passwords: invitationPasswords{err: sql.ErrNoRows}, tokens: &invitationTokens{}, sender: &invitationSender{cfg: cfg}, cfg: cfg, tx: invitationTx{}}
	ctx := clock.SetNow(admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: "operator-test"}), time.Now().UTC())
	return m, ctx
}

func TestUserInvitationsRequireOperatorBeforeDependencies(t *testing.T) {
	m := &manageUserInvitations{}
	if _, err := m.Create(context.Background(), InviteUserInput{}); !errors.Is(err, messages.ForbiddenError) {
		t.Fatal(err)
	}
	if _, err := m.Resend(context.Background(), shared.NewUUID[user.User]()); !errors.Is(err, messages.ForbiddenError) {
		t.Fatal(err)
	}
}

func TestUserInvitationValidationAndFailures(t *testing.T) {
	for _, input := range []InviteUserInput{
		{"", "user@example.com"}, {" ", "user@example.com"}, {strings.Repeat("名", 51), "user@example.com"},
		{"表示名", "bad"}, {"表示名", "Name <user@example.com>"}, {"表示名", "u@example.com\r\nBcc: x@example.com"}, {"表示名", strings.Repeat("a", 245) + "@example.com"},
		{strings.Repeat("😀", 51), "user@example.com"}, {"\u0085\u3000", "user@example.com"},
		{"表示名", `.user@example.com`}, {"表示名", `user.@example.com`}, {"表示名", `u..ser@example.com`},
		{"表示名", `user@.example.com`}, {"表示名", `user@example..com`}, {"表示名", `"user"@example.com`},
		{"表示名", `user@localhost`}, {"表示名", `user@例え.test`}, {"表示名", strings.Repeat("a", 243) + "@example.com"},
	} {
		m, ctx := invitationFixture()
		if _, err := m.Create(ctx, input); !errors.Is(err, messages.BadRequestError) {
			t.Fatalf("invalid input accepted: %v", err)
		}
		if m.sender.(*invitationSender).calls != 0 || m.users.(*invitationUsers).created != nil {
			t.Fatal("invalid input caused side effects")
		}
	}
	for _, website := range []string{"", "http://example.com", "https://evil@example.com", "https://example.com?next=evil", "https://example.com#fragment", "https://example.com/path"} {
		m, ctx := invitationFixture()
		m.cfg.WEBSITE_URL = website
		if _, err := m.Create(ctx, InviteUserInput{"表示名", "user@example.com"}); err == nil {
			t.Fatal("unsafe website accepted")
		}
		if m.users.(*invitationUsers).created != nil {
			t.Fatal("invalid configuration created user")
		}
	}
	m, ctx := invitationFixture()
	m.users.(*invitationUsers).existing = &user.User{}
	if _, err := m.Create(ctx, InviteUserInput{"表示名", "user@example.com"}); !errors.Is(err, messages.ManagedUserEmailExists) {
		t.Fatal(err)
	}
	dbErr := errors.New("database unavailable")
	m.users.(*invitationUsers).existing = nil
	m.users.(*invitationUsers).err = dbErr
	if _, err := m.Create(ctx, InviteUserInput{"表示名", "user@example.com"}); !errors.Is(err, dbErr) {
		t.Fatal("database failure treated as missing user", err)
	}
	if m.sender.(*invitationSender).calls != 0 {
		t.Fatal("failed creation sent mail")
	}
}

func TestUserInvitationUnicodeAndEmailBoundaries(t *testing.T) {
	for _, input := range []InviteUserInput{
		{strings.Repeat("名", 50), "user@example.com"}, {strings.Repeat("😀", 50), "user@example.com"},
		{"\u0085表示名\u0085", "\u0085user@example.com\u0085"},
		{"表示名", "user+tag@example.com"}, {"表示名", "user.name@sub.example.com"}, {"表示名", "USER@EXAMPLE.COM"},
		{"表示名", strings.Repeat("a", 242) + "@example.com"},
	} {
		m, ctx := invitationFixture()
		if _, err := m.Create(ctx, input); err != nil {
			t.Fatal("valid Unicode or email boundary rejected", err)
		}
	}
}

func TestUserInvitationCreationAndDeliveryFailureAudit(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	m, ctx := invitationFixture()
	out, err := m.Create(ctx, InviteUserInput{" 表示名 ", " user@example.com "})
	if err != nil || !out.InvitationSent {
		t.Fatalf("%#v %v", out, err)
	}
	u := m.users.(*invitationUsers).created
	if !u.Verify() || u.Provider() != "PASSWORD" || u.IsEmailVerified() || *u.DisplayName() != "表示名" || len(*u.DisplayID()) > 25 {
		t.Fatal("incorrect initial account state")
	}
	link, _ := url.Parse(m.sender.(*invitationSender).data["SetupURL"].(string))
	digest, valid := password_auth.SetupTokenHash(link.Fragment)
	pending := m.tokens.(*invitationTokens).value
	if link.Path != "/auth/password/setup" || link.RawQuery != "" || !valid || !bytes.Equal(digest, pending.TokenHash) || !pending.ExpiresAt.Equal(clock.Now(ctx).Add(24*time.Hour)) {
		t.Fatal("incorrect setup link")
	}
	if _, err := m.Resend(ctx, out.UserID); !errors.Is(err, messages.PasswordSetupRateLimit) {
		t.Fatal("cooldown missing", err)
	}
	later := clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
	m.sender.(*invitationSender).err = errors.New("private transport error: user@example.com " + link.Fragment)
	out, err = m.Resend(later, out.UserID)
	if err != nil || out.InvitationSent {
		t.Fatal("delivery failure must leave a retryable account", err)
	}
	if m.tokens.(*invitationTokens).saves != 2 || bytes.Equal(pending.TokenHash, m.tokens.(*invitationTokens).value.TokenHash) {
		t.Fatal("resend did not replace link")
	}
	for _, secret := range []string{"user@example.com", link.Fragment, "private transport error"} {
		if strings.Contains(logs.String(), secret) {
			t.Fatal("audit leaked email or token")
		}
	}
	if !strings.Contains(logs.String(), "operator-test") || !strings.Contains(logs.String(), u.UserID().String()) {
		t.Fatal("audit actor missing")
	}
}
