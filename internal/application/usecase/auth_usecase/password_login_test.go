package auth_usecase

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

var loginTestErr = errors.New("simulated login failure")

type loginTestTx struct{ err error }

func (tx *loginTestTx) ExecTx(ctx context.Context, fn func(context.Context) error) error {
	if tx.err != nil {
		return tx.err
	}
	return fn(ctx)
}

type loginTestUsers struct {
	user.UserRepository
	usr                *user.User
	findErr, updateErr error
	updates            int
}

func (r *loginTestUsers) FindBySubject(context.Context, user.UserSubject) (*user.User, error) {
	return r.usr, r.findErr
}
func (r *loginTestUsers) FindByDisplayID(context.Context, string) (*user.User, error) {
	return r.usr, r.findErr
}
func (r *loginTestUsers) FindByID(context.Context, shared.UUID[user.User]) (*user.User, error) {
	return r.usr, r.findErr
}
func (r *loginTestUsers) Update(context.Context, user.User) error { r.updates++; return r.updateErr }

type loginTestPasswords struct {
	password_auth.PasswordAuthManager
	password_auth.PasswordAuthRepository
	valid             bool
	verifyErr, getErr error
	auth              password_auth.PasswordAuth
}

func (r *loginTestPasswords) VerifyPassword(context.Context, shared.UUID[user.User], string) (bool, error) {
	return r.valid, r.verifyErr
}
func (r *loginTestPasswords) GetPasswordAuthByUserID(context.Context, shared.UUID[user.User]) (password_auth.PasswordAuth, error) {
	return r.auth, r.getErr
}

type loginTestChallenges struct {
	c                   password_auth.LoginChallenge
	saves, consumes     int
	consumed            bool
	saveErr, consumeErr error
}

func (r *loginTestChallenges) Save(_ context.Context, c password_auth.LoginChallenge, _ time.Time) error {
	r.saves++
	if r.saveErr != nil {
		return r.saveErr
	}
	r.c = c
	r.consumed = false
	return nil
}
func (r *loginTestChallenges) Consume(_ context.Context, token, browser []byte, now time.Time) (password_auth.LoginChallenge, error) {
	r.consumes++
	if r.consumeErr != nil {
		return password_auth.LoginChallenge{}, r.consumeErr
	}
	if r.consumed || r.c.PINAttempts >= 5 || !now.Before(r.c.ExpiresAt) || !bytes.Equal(token, r.c.TokenHash) || !bytes.Equal(browser, r.c.BrowserHash) {
		return password_auth.LoginChallenge{}, sql.ErrNoRows
	}
	r.consumed = true
	return r.c, nil
}

func (r *loginTestChallenges) FindForPIN(_ context.Context, browser []byte) (password_auth.LoginChallenge, error) {
	if !bytes.Equal(browser, r.c.BrowserHash) {
		return password_auth.LoginChallenge{}, sql.ErrNoRows
	}
	c := r.c
	c.Consumed = r.consumed
	return c, nil
}
func (r *loginTestChallenges) RecordPINFailure(_ context.Context, browser []byte) (int, error) {
	if !bytes.Equal(browser, r.c.BrowserHash) || r.c.PINAttempts >= 5 {
		return 0, sql.ErrNoRows
	}
	r.c.PINAttempts++
	return r.c.PINAttempts, nil
}

type loginTestSender struct {
	cfg   *config.Config
	calls int
	data  map[string]any
	err   error
}

func (s *loginTestSender) Send(_ context.Context, to string, tmpl email_template.EmailTemplateType, data map[string]any) error {
	s.calls++
	s.data = data
	if _, err := email.RenderMessage(s.cfg, to, tmpl, data); err != nil {
		return err
	}
	return s.err
}

type loginTestSessions struct {
	session.SessionRepository
	session.SessionService
	created, deactivated     int
	createErr, deactivateErr error
}

type loginTestOrganizations struct {
	organization.OrganizationRepository
	org *organization.Organization
}

func (r *loginTestOrganizations) FindByID(context.Context, shared.UUID[organization.Organization]) (*organization.Organization, error) {
	return r.org, nil
}

type loginTestMemberships struct {
	organization.OrganizationUserRepository
	memberships []*organization.OrganizationUser
	err         error
}

func (r *loginTestMemberships) FindByUserID(context.Context, shared.UUID[user.User]) ([]*organization.OrganizationUser, error) {
	return r.memberships, r.err
}

func (s *loginTestSessions) Create(_ context.Context, sess session.Session) (*session.Session, error) {
	s.created++
	return &sess, s.createErr
}
func (s *loginTestSessions) DeactivateUserSessions(context.Context, shared.UUID[user.User]) error {
	s.deactivated++
	return s.deactivateErr
}
func (s *loginTestSessions) DeactivateAllByUserID(ctx context.Context, id shared.UUID[user.User]) error {
	return s.DeactivateUserSessions(ctx, id)
}

type loginTestTokens struct {
	session.TokenManager
	calls int
	err   error
}

func (t *loginTestTokens) Generate(context.Context, user.User, shared.UUID[session.Session]) (string, error) {
	t.calls++
	return "signed-session", t.err
}

type loginTestHarness struct {
	p          *passwordLoginInteractor
	ctx        context.Context
	users      *loginTestUsers
	passwords  *loginTestPasswords
	challenges *loginTestChallenges
	sender     *loginTestSender
	sessions   *loginTestSessions
	tokens     *loginTestTokens
	tx         *loginTestTx
}

func newLoginHarness() *loginTestHarness {
	cfg := &config.Config{Env: config.LOCAL, TokenSecret: "test-pin-secret", HASH_PEPPER: "test-pepper", EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://localhost:3000"}
	u := user.NewUser(shared.NewUUID[user.User](), nil, nil, "subject", "PASSWORD", nil)
	u.ChangeEmail("city@example.com")
	h := &loginTestHarness{ctx: clock.SetNow(context.Background(), time.Date(2026, 10, 8, 0, 0, 0, 0, time.UTC)), users: &loginTestUsers{usr: &u}, passwords: &loginTestPasswords{valid: true, auth: password_auth.PasswordAuth{PasswordHash: "hash", Salt: "salt"}}, challenges: &loginTestChallenges{}, sender: &loginTestSender{cfg: cfg}, sessions: &loginTestSessions{}, tokens: &loginTestTokens{}, tx: &loginTestTx{}}
	h.p = &passwordLoginInteractor{UserRepository: h.users, SessionService: h.sessions, SessionRepository: h.sessions, TokenManager: h.tokens, txRunner: h.tx, Config: cfg, PasswordAuthManager: h.passwords, passwordRepository: h.passwords, challenges: h.challenges, emailSender: h.sender}
	h.p.organizationUserRepository = &loginTestMemberships{}
	h.p.organizationRepository = &loginTestOrganizations{}
	return h
}
func (h *loginTestHarness) begin(t *testing.T) (*PasswordLoginOutput, string) {
	t.Helper()
	out, err := h.p.Execute(h.ctx, PasswordLoginInput{IDorEmail: "city@example.com", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(h.sender.data["VerificationURL"].(string))
	if err != nil {
		t.Fatal(err)
	}
	return out, u.Fragment
}

func TestPasswordLoginRequiresFreshEmailEveryTime(t *testing.T) {
	h := newLoginHarness()
	h.users.usr.SetEmailVerified(true)
	out, token := h.begin(t)
	if h.sessions.created != 0 || h.sessions.deactivated != 0 || h.tokens.calls != 0 {
		t.Fatal("password alone created or revoked a session")
	}
	if !validLoginSecret(token) || !validLoginSecret(out.BrowserSecret) || token == out.BrowserSecret {
		t.Fatal("secrets are invalid or not independent")
	}
	if len(h.challenges.c.TokenHash) != 32 || bytes.Contains(h.challenges.c.TokenHash, []byte(token)) || !out.ExpiresAt.Equal(clock.Now(h.ctx).Add(EmailLoginTTL)) {
		t.Fatal("invalid persisted challenge")
	}
	if _, err := h.p.Verify(h.ctx, token, ""); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
		t.Fatal(err)
	}
	wrong, _ := newLoginSecret()
	if _, err := h.p.Verify(h.ctx, token, wrong); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
		t.Fatal(err)
	}
	if got, err := h.p.Verify(h.ctx, token, out.BrowserSecret); err != nil || got != "signed-session" {
		t.Fatalf("%s %v", got, err)
	}
	if h.sessions.created != 1 || h.users.updates != 1 || !h.users.usr.IsEmailVerified() {
		t.Fatal("verification did not establish the session")
	}
	if _, err := h.p.Verify(h.ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
		t.Fatalf("reused link accepted: %v", err)
	}
	_, _ = h.begin(t)
	if h.sender.calls != 2 || h.sessions.created != 1 {
		t.Fatal("past email verification bypassed a later login")
	}
}

func TestEmailLoginPINSharedChallenge(t *testing.T) {
	for _, method := range []string{"PIN", "link"} {
		t.Run(method, func(t *testing.T) {
			h := newLoginHarness()
			out, token := h.begin(t)
			pin := h.sender.data["PIN"].(string)
			if !validLoginPIN(pin) || bytes.Equal(h.challenges.c.PINHash, []byte(pin)) {
				t.Fatal("PIN not securely generated/stored")
			}
			var err error
			if method == "PIN" {
				_, err = h.p.VerifyPIN(h.ctx, pin, out.BrowserSecret)
			} else {
				_, err = h.p.Verify(h.ctx, token, out.BrowserSecret)
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err = h.p.VerifyPIN(h.ctx, pin, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginPINError) {
				t.Fatal("PIN replay accepted", err)
			}
			if _, err = h.p.Verify(h.ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
				t.Fatal("link replay accepted", err)
			}
		})
	}
}

func TestEmailLoginPINAttemptsExpiryAndResend(t *testing.T) {
	h := newLoginHarness()
	out, token := h.begin(t)
	pin := h.sender.data["PIN"].(string)
	other, _ := newLoginSecret()
	if _, err := h.p.VerifyPIN(h.ctx, pin, other); !errors.Is(err, messages.InvalidEmailLoginPINError) || h.challenges.c.PINAttempts != 0 {
		t.Fatal("other browser affected guesses")
	}
	expired := clock.SetNow(h.ctx, out.ExpiresAt)
	if _, err := h.p.VerifyPIN(expired, pin, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginPINError) {
		t.Fatal("expired PIN accepted")
	}
	for i := 1; i <= 5; i++ {
		_, err := h.p.VerifyPIN(h.ctx, "invalid", out.BrowserSecret)
		expected := messages.InvalidEmailLoginPINError
		if i == 5 {
			expected = messages.EmailLoginPINAttemptsError
		}
		if !errors.Is(err, expected) || h.challenges.c.PINAttempts != i {
			t.Fatalf("guess %d: %v count %d", i, err, h.challenges.c.PINAttempts)
		}
	}
	if _, err := h.p.VerifyPIN(h.ctx, pin, out.BrowserSecret); !errors.Is(err, messages.EmailLoginPINAttemptsError) {
		t.Fatal("locked PIN accepted")
	}
	if _, err := h.p.Verify(h.ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
		t.Fatal("locked link accepted")
	}
	h.ctx = clock.SetNow(h.ctx, clock.Now(h.ctx).Add(time.Minute))
	resent, _ := h.begin(t)
	if h.challenges.c.PINAttempts != 0 {
		t.Fatal("resend retained attempts")
	}
	if _, err := h.p.VerifyPIN(h.ctx, pin, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginPINError) {
		t.Fatal("old PIN accepted")
	}
	if _, err := h.p.VerifyPIN(h.ctx, h.sender.data["PIN"].(string), resent.BrowserSecret); err != nil {
		t.Fatal(err)
	}
}

func TestLoginPINLeadingZerosAndKeyBinding(t *testing.T) {
	h := newLoginHarness()
	if !validLoginPIN("000001") || validLoginPIN("１２３４５６") || validLoginPIN("12345") {
		t.Fatal("PIN format validation")
	}
	a := h.p.pinHash("000001", []byte("browser"))
	if bytes.Equal(a, h.p.pinHash("000001", []byte("other"))) {
		t.Fatal("missing browser binding")
	}
	h.p.Config.TokenSecret = "different-key"
	if bytes.Equal(a, h.p.pinHash("000001", []byte("browser"))) {
		t.Fatal("missing server key binding")
	}
}

func TestPasswordLoginStartFailures(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(*loginTestHarness)
		want  error
	}{
		{"transaction", func(h *loginTestHarness) { h.tx.err = loginTestErr }, loginTestErr},
		{"unknown account", func(h *loginTestHarness) { h.users.usr = nil }, messages.InvalidPasswordOrEmailError},
		{"lookup", func(h *loginTestHarness) { h.users.findErr = loginTestErr }, messages.InvalidPasswordOrEmailError},
		{"wrong password", func(h *loginTestHarness) { h.passwords.valid = false }, messages.InvalidPasswordOrEmailError},
		{"password verification", func(h *loginTestHarness) { h.passwords.verifyErr = loginTestErr }, messages.InvalidPasswordOrEmailError},
		{"missing email", func(h *loginTestHarness) { h.users.usr.ChangeEmail("") }, messages.InvalidPasswordOrEmailError},
		{"credential lookup", func(h *loginTestHarness) { h.passwords.getErr = loginTestErr }, loginTestErr},
		{"save", func(h *loginTestHarness) { h.challenges.saveErr = loginTestErr }, loginTestErr},
		{"resend too soon", func(h *loginTestHarness) { h.challenges.saveErr = sql.ErrNoRows }, messages.EmailLoginRateLimitError},
		{"send", func(h *loginTestHarness) { h.sender.err = loginTestErr }, loginTestErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newLoginHarness()
			tc.setup(h)
			out, err := h.p.Execute(h.ctx, PasswordLoginInput{IDorEmail: "city@example.com", Password: "password"})
			if out != nil || !errors.Is(err, tc.want) {
				t.Fatalf("%#v %v", out, err)
			}
			if h.sessions.created != 0 || h.sessions.deactivated != 0 || h.tokens.calls != 0 {
				t.Fatal("failure affected authenticated sessions")
			}
			if tc.name != "send" && h.sender.calls != 0 {
				t.Fatal("sent mail on failure")
			}
		})
	}
}

func TestPasswordEmailConfirmationFailures(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(*loginTestHarness)
		want  error
	}{
		{"expiry boundary", func(h *loginTestHarness) { h.ctx = clock.SetNow(h.ctx, h.challenges.c.ExpiresAt) }, messages.InvalidEmailLoginLinkError},
		{"repository", func(h *loginTestHarness) { h.challenges.consumeErr = loginTestErr }, loginTestErr},
		{"deleted account", func(h *loginTestHarness) { h.users.usr = nil }, messages.InvalidEmailLoginLinkError},
		{"changed password", func(h *loginTestHarness) { h.passwords.auth.PasswordHash = "new-hash" }, messages.InvalidEmailLoginLinkError},
		{"changed email", func(h *loginTestHarness) { h.users.usr.ChangeEmail("other@example.com") }, messages.InvalidEmailLoginLinkError},
		{"deleted password", func(h *loginTestHarness) { h.passwords.getErr = sql.ErrNoRows }, messages.InvalidEmailLoginLinkError},
		{"credential error", func(h *loginTestHarness) { h.passwords.getErr = loginTestErr }, loginTestErr},
		{"user update", func(h *loginTestHarness) { h.users.updateErr = loginTestErr }, loginTestErr},
		{"deactivation", func(h *loginTestHarness) { h.sessions.deactivateErr = loginTestErr }, loginTestErr},
		{"session creation", func(h *loginTestHarness) { h.sessions.createErr = loginTestErr }, loginTestErr},
		{"token generation", func(h *loginTestHarness) { h.tokens.err = loginTestErr }, loginTestErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newLoginHarness()
			out, token := h.begin(t)
			tc.setup(h)
			got, err := h.p.Verify(h.ctx, token, out.BrowserSecret)
			if got != "" || !errors.Is(err, tc.want) {
				t.Fatalf("%s %v", got, err)
			}
		})
	}
}

func TestPasswordEmailResendInvalidatesPreviousLink(t *testing.T) {
	h := newLoginHarness()
	old, token := h.begin(t)
	h.ctx = clock.SetNow(h.ctx, clock.Now(h.ctx).Add(time.Minute))
	out, newToken := h.begin(t)
	if _, err := h.p.Verify(h.ctx, token, old.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
		t.Fatal("old link remained valid")
	}
	if _, err := h.p.Verify(h.ctx, newToken, out.BrowserSecret); err != nil {
		t.Fatal(err)
	}
}

func TestPasswordEmailURLConfiguration(t *testing.T) {
	for _, base := range []string{"http://example.com", "https://user:pass@example.com", "https://example.com?x=1", "https://example.com#fragment", "/relative", "https://"} {
		t.Run(base, func(t *testing.T) {
			h := newLoginHarness()
			h.p.Config.EMAIL_VERIFICATION_BASE_URL = base
			if _, err := h.p.verificationURL("token"); err == nil {
				t.Fatal("invalid URL accepted")
			}
		})
	}
	h := newLoginHarness()
	h.p.Config.Env = config.PROD
	if _, err := h.p.verificationURL("token"); err == nil {
		t.Fatal("production silently used local URL")
	}
}

func TestPasswordEmailUsesConfiguredWebPage(t *testing.T) {
	h := newLoginHarness()
	h.p.Config.EMAIL_VERIFICATION_PAGE_URL = "https://localhost:3000/auth/email/verify"
	u, err := h.p.verificationURL("one-time-token")
	if err != nil || u != "https://localhost:3000/auth/email/verify#one-time-token" {
		t.Fatalf("%s %v", u, err)
	}
	for _, invalid := range []string{"http://example.com/verify", "https://user:pass@example.com/verify", "https://example.com/verify?token=fixed", "https://example.com/verify#fixed", "/relative"} {
		h.p.Config.EMAIL_VERIFICATION_PAGE_URL = invalid
		if _, err := h.p.verificationURL("token"); err == nil {
			t.Fatalf("unsafe page URL accepted: %s", invalid)
		}
	}
}
