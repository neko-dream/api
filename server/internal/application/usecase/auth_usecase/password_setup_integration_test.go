package auth_usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	userquery "github.com/neko-dream/api/internal/infrastructure/persistence/query/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/pkg/hash"
)

type failingSetupManager struct {
	password_auth.PasswordAuthManager
	fail bool
}

type failingSetupTokenSave struct {
	password_auth.SetupTokenRepository
}

func (failingSetupTokenSave) Save(context.Context, password_auth.SetupToken) error {
	return loginTestErr
}

func (m *failingSetupManager) RegisterPassword(ctx context.Context, id shared.UUID[user.User], password string, requireChange bool) error {
	if err := m.PasswordAuthManager.RegisterPassword(ctx, id, password, requireChange); err != nil {
		return err
	}
	if m.fail {
		return loginTestErr
	}
	return nil
}

func TestManagedUserPasswordSetupIntegration(t *testing.T) {
	conn := emailLoginTestDB(t)
	dbm := db.NewDBManager(conn)
	cfg := &config.Config{Env: config.LOCAL, HASH_PEPPER: "setup-pepper", HASH_ITERATIONS: 4, EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://localhost:3000", ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
	enc, err := crypto.NewEncryptor(cfg)
	if err != nil {
		t.Fatal(err)
	}
	users := repository.NewUserRepository(dbm, nil, enc)
	passwords := repository.NewPasswordAuthRepository(dbm)
	tokens := repository.NewPasswordSetupTokenRepository(dbm)
	manager := service.NewPasswordAuthManager(passwords, cfg)
	query := userquery.NewManagementQuery(dbm, enc)
	t.Run("creation rollback", func(t *testing.T) {
		ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: "test-operator"})
		sender := &loginTestSender{cfg: cfg}
		invite := admin_usecase.NewManageUserInvitations(users, passwords, failingSetupTokenSave{tokens}, sender, cfg, dbm)
		address := "rollback-" + uuid.NewString() + "@example.com"
		if _, err := invite.Create(ctx, admin_usecase.InviteUserInput{DisplayName: "表示名", Email: address}); !errors.Is(err, loginTestErr) {
			t.Fatal(err)
		}
		subject, _ := hash.HMACHash(address, cfg.HASH_PEPPER)
		if u, err := users.FindBySubject(ctx, user.UserSubject(subject)); !errors.Is(err, sql.ErrNoRows) || u != nil {
			t.Fatal("failed creation left an account", err)
		}
		if sender.calls != 0 {
			t.Fatal("failed transaction sent an invitation")
		}
	})
	for _, scenario := range []string{"success and replay", "expiry", "resend", "delivery failure", "password failure rollback", "withdrawn", "changed email", "changed subject", "existing password", "completion concurrency", "resend completion concurrency", "password length"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := clock.SetNow(admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: "test-operator"}), time.Now().UTC())
			sender := &loginTestSender{cfg: cfg}
			if scenario == "delivery failure" {
				sender.err = loginTestErr
			}
			invite := admin_usecase.NewManageUserInvitations(users, passwords, tokens, sender, cfg, dbm)
			address := "setup-" + uuid.NewString() + "@example.com"
			out, err := invite.Create(ctx, admin_usecase.InviteUserInput{DisplayName: "表示名", Email: address})
			if err != nil || out.InvitationSent != (scenario != "delivery failure") {
				t.Fatalf("create: %#v %v", out, err)
			}
			if _, err := invite.Create(ctx, admin_usecase.InviteUserInput{DisplayName: "別の名前", Email: address}); !errors.Is(err, messages.ManagedUserEmailExists) {
				t.Fatal("duplicate accepted", err)
			}
			u, err := users.FindByID(ctx, out.UserID)
			if err != nil {
				t.Fatal(err)
			}
			if !u.Verify() || u.IsEmailVerified() {
				t.Fatal("invalid pending profile")
			}
			pending, err := query.FindByID(ctx, out.UserID)
			if err != nil || !pending.PasswordSetupRequired || pending.CanCreateTalkSession || pending.Email == nil || *pending.Email != address {
				t.Fatal("incorrect pending projection", err)
			}
			for _, table := range []string{"password_auth", "sessions", "organization_users", "policy_consents"} {
				var count int
				if err := conn.QueryRow("SELECT count(*) FROM "+table+" WHERE user_id=$1", out.UserID.UUID()).Scan(&count); err != nil || count != 0 {
					t.Fatalf("creation changed %s: %d %v", table, count, err)
				}
			}
			if scenario == "delivery failure" {
				sender.err = nil
				ctx = clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
				if resent, err := invite.Resend(ctx, out.UserID); err != nil || !resent.InvitationSent {
					t.Fatal("failed delivery not retryable", err)
				}
			}
			link, _ := url.Parse(sender.data["SetupURL"].(string))
			input := SetupPasswordInput{Token: link.Fragment, Password: "new-password-for-setup"}
			fault := &failingSetupManager{PasswordAuthManager: manager}
			setup := NewSetupPassword(users, passwords, fault, tokens, cfg, dbm)
			switch scenario {
			case "expiry":
				ctx = clock.SetNow(ctx, clock.Now(ctx).Add(24*time.Hour))
				if err := setup.Execute(ctx, input); !errors.Is(err, messages.InvalidPasswordSetupLink) {
					t.Fatal("expired accepted", err)
				}
				return
			case "resend":
				if _, err := invite.Resend(ctx, out.UserID); !errors.Is(err, messages.PasswordSetupRateLimit) {
					t.Fatal("resend limit missing", err)
				}
				ctx = clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
				if _, err := invite.Resend(ctx, out.UserID); err != nil {
					t.Fatal(err)
				}
				if err := setup.Execute(ctx, input); !errors.Is(err, messages.InvalidPasswordSetupLink) {
					t.Fatal("old link accepted", err)
				}
				link, _ = url.Parse(sender.data["SetupURL"].(string))
				input.Token = link.Fragment
			case "password failure rollback":
				fault.fail = true
				if err := setup.Execute(ctx, input); !errors.Is(err, loginTestErr) {
					t.Fatal(err)
				}
				if _, err := passwords.GetPasswordAuthByUserID(ctx, out.UserID); !errors.Is(err, sql.ErrNoRows) {
					t.Fatal("credential not rolled back")
				}
				u, _ = users.FindByID(ctx, out.UserID)
				if u.IsEmailVerified() {
					t.Fatal("verification not rolled back")
				}
				if _, err := tokens.FindByUserID(ctx, out.UserID); err != nil {
					t.Fatal("link consumed on failure")
				}
				fault.fail = false
			case "withdrawn", "changed email", "changed subject", "existing password":
				switch scenario {
				case "withdrawn":
					u.Withdraw(clock.Now(ctx))
					err = users.Update(ctx, *u)
				case "changed email":
					u.ChangeEmail("changed@example.com")
					err = users.Update(ctx, *u)
				case "changed subject":
					err = users.ChangeSubject(ctx, out.UserID, "changed-subject")
				case "existing password":
					err = manager.RegisterPassword(ctx, out.UserID, "previous-password", false)
				}
				if err != nil {
					t.Fatal(err)
				}
				if err := setup.Execute(ctx, input); !errors.Is(err, messages.InvalidPasswordSetupLink) {
					t.Fatal("ineligible account accepted", err)
				}
				ctx = clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
				if _, err := invite.Resend(ctx, out.UserID); !errors.Is(err, messages.PasswordSetupNotPending) {
					t.Fatal("ineligible resend accepted", err)
				}
				return
			case "completion concurrency":
				var wg sync.WaitGroup
				errs := make(chan error, 8)
				for range 8 {
					wg.Add(1)
					go func() { defer wg.Done(); errs <- setup.Execute(ctx, input) }()
				}
				wg.Wait()
				close(errs)
				successes := 0
				for err := range errs {
					if err == nil {
						successes++
					} else if !errors.Is(err, messages.InvalidPasswordSetupLink) {
						t.Fatal(err)
					}
				}
				if successes != 1 {
					t.Fatalf("single use failed: %d", successes)
				}
			case "resend completion concurrency":
				ctx = clock.SetNow(ctx, clock.Now(ctx).Add(time.Minute))
				var wg sync.WaitGroup
				wg.Add(2)
				errs := make(chan error, 2)
				go func() { defer wg.Done(); errs <- setup.Execute(ctx, input) }()
				go func() { defer wg.Done(); _, err := invite.Resend(ctx, out.UserID); errs <- err }()
				wg.Wait()
				close(errs)
				for err := range errs {
					if err != nil && !errors.Is(err, messages.InvalidPasswordSetupLink) && !errors.Is(err, messages.PasswordSetupNotPending) {
						t.Fatal(err)
					}
				}
				if _, err := passwords.GetPasswordAuthByUserID(ctx, out.UserID); errors.Is(err, sql.ErrNoRows) {
					link, _ = url.Parse(sender.data["SetupURL"].(string))
					input.Token = link.Fragment
					if err := setup.Execute(ctx, input); err != nil {
						t.Fatal(err)
					}
				}
			case "password length":
				for _, tc := range []struct {
					password string
					want     error
				}{{"short", messages.PasswordTooShort}, {strings.Repeat("a", 73-len(cfg.HASH_PEPPER)), messages.PasswordTooLong}, {strings.Repeat("あ", 24), messages.PasswordTooLong}} {
					if err := setup.Execute(ctx, SetupPasswordInput{Token: input.Token, Password: tc.password}); !errors.Is(err, tc.want) {
						t.Fatal("wrong password validation", err)
					}
				}
				input.Password = strings.Repeat("a", 72-len(cfg.HASH_PEPPER))
			}
			if scenario != "completion concurrency" && scenario != "resend completion concurrency" {
				if err := setup.Execute(ctx, input); err != nil {
					t.Fatal(err)
				}
			}
			if err := setup.Execute(ctx, input); !errors.Is(err, messages.InvalidPasswordSetupLink) {
				t.Fatal("replay accepted", err)
			}
			if valid, err := manager.VerifyPassword(ctx, out.UserID, input.Password); err != nil || !valid {
				t.Fatal("credential cannot log in", err)
			}
			credential, _ := passwords.GetPasswordAuthByUserID(ctx, out.UserID)
			if credential.RequiredPasswordChange || credential.PasswordHash == input.Password {
				t.Fatal("unsafe credential")
			}
			u, _ = users.FindByID(ctx, out.UserID)
			if !u.IsEmailVerified() {
				t.Fatal("email not verified")
			}
			pending, _ = query.FindByID(ctx, out.UserID)
			if pending.PasswordSetupRequired {
				t.Fatal("pending state not cleared")
			}
			if _, err := invite.Resend(ctx, out.UserID); !errors.Is(err, messages.PasswordSetupNotPending) {
				t.Fatal("completed account can be reset", err)
			}
			var sessions int
			if err := conn.QueryRow("SELECT count(*) FROM sessions WHERE user_id=$1", out.UserID.UUID()).Scan(&sessions); err != nil || sessions != 0 {
				t.Fatal("setup logged user in", err)
			}
			if scenario == "success and replay" {
				cfg.TokenSecret = "local-setup-login-test-secret"
				sessionRepo := repository.NewSessionRepository(dbm)
				login := NewPasswordLogin(users, service.NewSessionService(sessionRepo), sessionRepo, dbm, cfg, manager, &loginTestTokens{}, passwords, repository.NewPasswordLoginChallengeRepository(dbm), sender)
				started, err := login.Execute(ctx, PasswordLoginInput{IDorEmail: address, Password: input.Password})
				if err != nil {
					t.Fatal("new account cannot start ordinary login", err)
				}
				if err := conn.QueryRow("SELECT count(*) FROM sessions WHERE user_id=$1", out.UserID.UUID()).Scan(&sessions); err != nil || sessions != 0 {
					t.Fatal("password alone logged in new account", err)
				}
				if _, err := login.VerifyPIN(ctx, sender.data["PIN"].(string), started.BrowserSecret); err != nil {
					t.Fatal("new account cannot complete 6-digit login", err)
				}
				if err := conn.QueryRow("SELECT count(*) FROM sessions WHERE user_id=$1", out.UserID.UUID()).Scan(&sessions); err != nil || sessions != 1 {
					t.Fatal("confirmation did not create session", err)
				}
			}
		})
	}
	t.Run("duplicate creation concurrency", func(t *testing.T) {
		ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: "test-operator"})
		address := "concurrent-" + uuid.NewString() + "@example.com"
		var wg sync.WaitGroup
		errs := make(chan error, 8)
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				invite := admin_usecase.NewManageUserInvitations(users, passwords, tokens, &loginTestSender{cfg: cfg}, cfg, dbm)
				_, err := invite.Create(ctx, admin_usecase.InviteUserInput{DisplayName: "表示名", Email: address})
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)
		successes := 0
		for err := range errs {
			if err == nil {
				successes++
			} else if !errors.Is(err, messages.ManagedUserEmailExists) {
				t.Fatal(err)
			}
		}
		if successes != 1 {
			t.Fatalf("duplicate accounts created: %d", successes)
		}
	})
}
