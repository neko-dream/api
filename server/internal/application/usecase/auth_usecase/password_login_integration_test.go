package auth_usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	tokens "github.com/neko-dream/api/internal/infrastructure/auth/session"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/pkg/hash"
)

type failingLoginTokens struct {
	session.TokenManager
	fail bool
}

func (t *failingLoginTokens) Generate(ctx context.Context, u user.User, id shared.UUID[session.Session]) (string, error) {
	if t.fail {
		return "", loginTestErr
	}
	return t.TokenManager.Generate(ctx, u, id)
}

func emailLoginTestDB(t *testing.T) *sql.DB {
	t.Helper()
	raw := os.Getenv("MAIL_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("set MAIL_TEST_DATABASE_URL to run real login transaction tests")
	}
	admin, err := sql.Open("pgx", raw)
	if err != nil {
		t.Fatal(err)
	}
	schema := "email_login_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.Exec(`CREATE SCHEMA "` + schema + `"`); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(`DROP SCHEMA "` + schema + `" CASCADE`); err != nil {
			t.Error(err)
		}
		admin.Close()
	})
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema+",public")
	u.RawQuery = q.Encode()
	if err := db.NewMigrator(&config.Config{DatabaseURL: u.String()}).Up(); err != nil {
		t.Fatal(err)
	}
	conn, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestPasswordEmailLoginTransactionsIntegration(t *testing.T) {
	conn := emailLoginTestDB(t)
	manager := db.NewDBManager(conn)
	cfg := &config.Config{Env: config.LOCAL, TokenSecret: "local-integration-secret", HASH_PEPPER: "local-integration-pepper", HASH_ITERATIONS: 4, EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://localhost:3000", ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
	encryptor, err := crypto.NewEncryptor(cfg)
	if err != nil {
		t.Fatal(err)
	}
	users := repository.NewUserRepository(manager, nil, encryptor)
	passwords := repository.NewPasswordAuthRepository(manager)
	passwordManager := service.NewPasswordAuthManager(passwords, cfg)
	sessions := repository.NewSessionRepository(manager)
	sessionService := service.NewSessionService(sessions)
	realTokens := tokens.NewSessionTokenManager(cfg, manager, sessions, users, repository.NewOrganizationUserRepository(manager), repository.NewOrganizationRepository(manager))
	for _, scenario := range []string{"single use concurrency", "session failure rollback", "resend and expiry", "changed password", "changed email", "send failure rollback", "PIN concurrency", "PIN attempts and resend", "PIN session rollback"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := clock.SetNow(context.Background(), time.Now().UTC())
			address := "login-" + uuid.NewString() + "@example.com"
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
			previous := session.NewSession(shared.NewUUID[session.Session](), usr.UserID(), usr.Provider(), session.SESSION_ACTIVE, *session.NewExpiresAt(ctx), clock.Now(ctx))
			if _, err := sessions.Create(ctx, *previous); err != nil {
				t.Fatal(err)
			}
			sender := &loginTestSender{cfg: cfg}
			tokenManager := &failingLoginTokens{TokenManager: realTokens}
			p := NewPasswordLogin(users, sessionService, sessions, manager, cfg, passwordManager, tokenManager, passwords, repository.NewPasswordLoginChallengeRepository(manager), sender)
			input := PasswordLoginInput{IDorEmail: address, Password: "original-password"}
			if scenario == "send failure rollback" {
				sender.err = loginTestErr
				if out, err := p.Execute(ctx, input); out != nil || !errors.Is(err, loginTestErr) {
					t.Fatalf("%#v %v", out, err)
				}
				var count int
				if err := conn.QueryRow("SELECT count(*) FROM password_login_challenges WHERE user_id=$1", usr.UserID().UUID()).Scan(&count); err != nil || count != 0 {
					t.Fatalf("failed send persisted a challenge: %d %v", count, err)
				}
				sender.err = nil
			}
			out, err := p.Execute(ctx, input)
			if err != nil {
				t.Fatal(err)
			}
			link, err := url.Parse(sender.data["VerificationURL"].(string))
			if err != nil {
				t.Fatal(err)
			}
			token := link.Fragment
			pin := sender.data["PIN"].(string)
			all, err := sessions.FindByUserID(ctx, usr.UserID())
			if err != nil || len(all) != 1 || !all[0].IsActive(ctx) {
				t.Fatalf("password affected previous sessions: %#v %v", all, err)
			}
			if _, err := p.Verify(ctx, token, ""); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
				t.Fatal(err)
			}
			other, _ := newLoginSecret()
			if _, err := p.Verify(ctx, token, other); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
				t.Fatal(err)
			}
			switch scenario {
			case "PIN attempts and resend":
				var wg sync.WaitGroup
				errs := make(chan error, 10)
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() { defer wg.Done(); _, err := p.VerifyPIN(ctx, "wrong", out.BrowserSecret); errs <- err }()
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					if !errors.Is(err, messages.InvalidEmailLoginPINError) && !errors.Is(err, messages.EmailLoginPINAttemptsError) {
						t.Fatal(err)
					}
				}
				var count int
				if err := conn.QueryRow("SELECT pin_attempts FROM password_login_challenges WHERE user_id=$1", usr.UserID().UUID()).Scan(&count); err != nil || count != 5 {
					t.Fatalf("attempts not persisted/capped: %d %v", count, err)
				}
				if _, err := p.VerifyPIN(ctx, pin, out.BrowserSecret); !errors.Is(err, messages.EmailLoginPINAttemptsError) {
					t.Fatal("locked PIN accepted", err)
				}
				if _, err := p.Verify(ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
					t.Fatal("locked link accepted", err)
				}
				later := clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
				resent, err := p.Execute(later, input)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := p.VerifyPIN(later, pin, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginPINError) {
					t.Fatal("old PIN survived resend")
				}
				ctx = later
				out = resent
				u, _ := url.Parse(sender.data["VerificationURL"].(string))
				token = u.Fragment
				pin = sender.data["PIN"].(string)
			case "PIN session rollback":
				tokenManager.fail = true
				if got, err := p.VerifyPIN(ctx, pin, out.BrowserSecret); got != "" || !errors.Is(err, loginTestErr) {
					t.Fatal("expected failed session", err)
				}
				var consumed sql.NullTime
				if err := conn.QueryRow("SELECT consumed_at FROM password_login_challenges WHERE user_id=$1", usr.UserID().UUID()).Scan(&consumed); err != nil || consumed.Valid {
					t.Fatal("failed session consumed PIN", err)
				}
				read, err := users.FindByID(ctx, usr.UserID())
				if err != nil || read.IsEmailVerified() {
					t.Fatal("email verification not rolled back")
				}
				tokenManager.fail = false
			case "PIN concurrency":
				var wg sync.WaitGroup
				errs := make(chan error, 8)
				for i := 0; i < 8; i++ {
					wg.Add(1)
					go func(usePIN bool) {
						defer wg.Done()
						var err error
						if usePIN {
							_, err = p.VerifyPIN(ctx, pin, out.BrowserSecret)
						} else {
							_, err = p.Verify(ctx, token, out.BrowserSecret)
						}
						errs <- err
					}(i%2 == 0)
				}
				wg.Wait()
				close(errs)
				success := 0
				for err := range errs {
					if err == nil {
						success++
					} else if !errors.Is(err, messages.InvalidEmailLoginPINError) && !errors.Is(err, messages.InvalidEmailLoginLinkError) {
						t.Fatal(err)
					}
				}
				if success != 1 {
					t.Fatalf("link/PIN race created %d sessions", success)
				}
			case "single use concurrency":
				var wg sync.WaitGroup
				results := make(chan error, 8)
				for i := 0; i < 8; i++ {
					wg.Add(1)
					go func() { defer wg.Done(); _, err := p.Verify(ctx, token, out.BrowserSecret); results <- err }()
				}
				wg.Wait()
				close(results)
				success := 0
				for err := range results {
					if err == nil {
						success++
					} else if !errors.Is(err, messages.InvalidEmailLoginLinkError) {
						t.Fatal(err)
					}
				}
				if success != 1 {
					t.Fatalf("concurrent verification succeeded %d times", success)
				}
			case "session failure rollback":
				tokenManager.fail = true
				if got, err := p.Verify(ctx, token, out.BrowserSecret); got != "" || !errors.Is(err, loginTestErr) {
					t.Fatalf("%s %v", got, err)
				}
				all, err := sessions.FindByUserID(ctx, usr.UserID())
				if err != nil || len(all) != 1 || !all[0].IsActive(ctx) {
					t.Fatalf("session failure was not rolled back: %#v %v", all, err)
				}
				read, err := users.FindByID(ctx, usr.UserID())
				if err != nil || read.IsEmailVerified() {
					t.Fatal("email verification update was not rolled back")
				}
				tokenManager.fail = false
			case "resend and expiry":
				if _, err := p.Execute(ctx, input); !errors.Is(err, messages.EmailLoginRateLimitError) {
					t.Fatalf("resend cooldown ignored: %v", err)
				}
				later := clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
				resent, err := p.Execute(later, input)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := p.Verify(later, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
					t.Fatal("resend left old link active")
				}
				u, _ := url.Parse(sender.data["VerificationURL"].(string))
				token = u.Fragment
				out = resent
				ctx = later
				expired := clock.SetNow(ctx, out.ExpiresAt)
				if _, err := p.Verify(expired, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
					t.Fatalf("expired link accepted: %v", err)
				}
			case "changed password":
				if err := passwordManager.UpdatePassword(ctx, usr.UserID(), "new-password"); err != nil {
					t.Fatal(err)
				}
				if _, err := p.Verify(ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
					t.Fatal("password change did not invalidate link")
				}
				return
			case "changed email":
				usr.ChangeEmail("changed@example.com")
				if err := users.Update(ctx, usr); err != nil {
					t.Fatal(err)
				}
				if _, err := p.Verify(ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
					t.Fatal("email change did not invalidate link")
				}
				return
			}
			if scenario != "single use concurrency" && scenario != "PIN concurrency" {
				var signed string
				var err error
				if strings.HasPrefix(scenario, "PIN") {
					signed, err = p.VerifyPIN(ctx, pin, out.BrowserSecret)
				} else {
					signed, err = p.Verify(ctx, token, out.BrowserSecret)
				}
				if err != nil {
					t.Fatal(err)
				}
				claim, err := realTokens.Parse(ctx, signed)
				if err != nil || claim == nil || !claim.IsEmailVerified {
					t.Fatalf("invalid authenticated token: %#v %v", claim, err)
				}
			}
			if _, err := p.Verify(ctx, token, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginLinkError) {
				t.Fatal("reused link accepted")
			}
			if _, err := p.VerifyPIN(ctx, pin, out.BrowserSecret); !errors.Is(err, messages.InvalidEmailLoginPINError) {
				t.Fatal("reused PIN accepted", err)
			}
			all, err = sessions.FindByUserID(ctx, usr.UserID())
			if err != nil || len(all) != 1 {
				t.Fatalf("wrong session count: %d %v", len(all), err)
			}
			old, err := sessions.FindBySessionID(ctx, previous.SessionID())
			if err != nil || old == nil || old.IsActive(ctx) {
				t.Fatal("previous session remained active")
			}
			active := 0
			for _, s := range all {
				if s.IsActive(ctx) {
					active++
				}
			}
			if active != 1 {
				t.Fatalf("expected exactly one active session, got %d", active)
			}
		})
	}
}
