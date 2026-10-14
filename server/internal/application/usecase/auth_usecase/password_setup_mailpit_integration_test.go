package auth_usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/service"
	sessiontokens "github.com/neko-dream/api/internal/infrastructure/auth/session"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/email"
	"github.com/neko-dream/api/internal/infrastructure/email/delivery"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	userqueryinfra "github.com/neko-dream/api/internal/infrastructure/persistence/query/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
)

type setupMailpitSender struct {
	email.EmailSender
	data map[string]any
}

func (s *setupMailpitSender) Send(ctx context.Context, to string, tmpl email_template.EmailTemplateType, data map[string]any) error {
	if err := s.EmailSender.Send(ctx, to, tmpl, data); err != nil {
		return err
	}
	s.data = data
	return nil
}

func TestManagedUserSetupAndLoginWithMailpit(t *testing.T) {
	address, mailbox := os.Getenv("MAIL_TEST_SMTP_ADDR"), os.Getenv("MAIL_TEST_HTTP_URL")
	if address == "" {
		t.Skip("set MAIL_TEST_SMTP_ADDR and MAIL_TEST_HTTP_URL to test real setup and login delivery")
	}
	mailboxURL, err := url.Parse(mailbox)
	if err != nil || mailboxURL.Scheme != "http" || (mailboxURL.Hostname() != "127.0.0.1" && mailboxURL.Hostname() != "localhost") {
		t.Fatal("MAIL_TEST_HTTP_URL must be a loopback Mailpit URL")
	}
	conn := emailLoginTestDB(t)
	dbm := db.NewDBManager(conn)
	cfg := &config.Config{Env: config.LOCAL, SMTP_ADDR: address, HASH_PEPPER: "local-test-pepper", HASH_ITERATIONS: 4, TokenSecret: "local-mailpit-flow-token-secret", EMAIL_FROM: "ことひろ <noreply@localhost.test>", APP_NAME: "ことひろ", WEBSITE_URL: "https://localhost:3002", EMAIL_VERIFICATION_PAGE_URL: "https://localhost:3002/auth/email/verify", ENCRYPTION_VERSION: "v1", ENCRYPTION_SECRET: "local-dev-key16!"}
	enc, err := crypto.NewEncryptor(cfg)
	if err != nil {
		t.Fatal(err)
	}
	transport, err := delivery.NewSender(cfg)
	if err != nil {
		t.Fatal(err)
	}
	sender := &setupMailpitSender{EmailSender: transport}
	users := repository.NewUserRepository(dbm, nil, enc)
	passwords := repository.NewPasswordAuthRepository(dbm)
	setupTokens := repository.NewPasswordSetupTokenRepository(dbm)
	manager := service.NewPasswordAuthManager(passwords, cfg)
	ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: "local-flow-test-operator"})
	recipient := "setup-flow-" + uuid.NewString() + "@localhost.test"
	invite := admin_usecase.NewManageUserInvitations(users, passwords, setupTokens, sender, cfg, dbm)
	out, err := invite.Create(ctx, admin_usecase.InviteUserInput{DisplayName: "ローカル検証", Email: recipient})
	if err != nil || !out.InvitationSent {
		t.Fatal("account invitation failed", err)
	}
	link, err := url.Parse(sender.data["SetupURL"].(string))
	if err != nil {
		t.Fatal(err)
	}
	assertSetupMailpitBody(t, mailbox, recipient, link.Fragment)
	setup := NewSetupPassword(users, passwords, manager, setupTokens, cfg, dbm)
	password := "local-test-new-password"
	if err := setup.Execute(ctx, SetupPasswordInput{Token: link.Fragment, Password: password}); err != nil {
		t.Fatal("initial password setup failed", err)
	}
	sessions := repository.NewSessionRepository(dbm)
	tokens := sessiontokens.NewSessionTokenManager(cfg, dbm, sessions, users, repository.NewOrganizationUserRepository(dbm), repository.NewOrganizationRepository(dbm))
	login := NewPasswordLogin(users, service.NewSessionService(sessions), sessions, dbm, cfg, manager, tokens, passwords, repository.NewPasswordLoginChallengeRepository(dbm), sender)
	started, err := login.Execute(ctx, PasswordLoginInput{IDorEmail: recipient, Password: password})
	if err != nil {
		t.Fatal("password login failed", err)
	}
	pin := sender.data["PIN"].(string)
	assertSetupMailpitBody(t, mailbox, recipient, pin)
	all, err := sessions.FindByUserID(ctx, out.UserID)
	if err != nil || len(all) != 0 {
		t.Fatal("password alone created a login session", err)
	}
	signed, err := login.VerifyPIN(ctx, pin, started.BrowserSecret)
	if err != nil || signed == "" {
		t.Fatal("mail confirmation failed", err)
	}
	all, err = sessions.FindByUserID(ctx, out.UserID)
	if err != nil || len(all) != 1 || !all[0].IsActive(ctx) {
		t.Fatal("mail confirmation did not create an active session", err)
	}
	profile, err := userqueryinfra.NewDetailHandler(dbm, enc).Execute(ctx, userquery.DetailInput{UserID: out.UserID})
	if err != nil {
		t.Fatal("newly invited user cannot load their profile after login", err)
	}
	if profile.User.DisplayName != "ローカル検証" || profile.User.Email == nil || *profile.User.Email != recipient || !profile.User.EmailVerified || profile.User.UserDemographic != nil {
		t.Fatal("newly invited user's profile is incomplete or invented demographics")
	}
	if _, err := dbm.GetQueries(ctx).GetUserDemographicByUserID(ctx, out.UserID.UUID()); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal("setup, login or profile loading created unnecessary demographics", err)
	}
}

func assertSetupMailpitBody(t *testing.T, base, recipient, secret string) {
	t.Helper()
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(base + "/api/v1/messages")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatal("Mailpit mailbox unavailable")
	}
	var mailbox struct {
		Messages []struct {
			ID string
			To []struct{ Address string }
		}
	}
	if err := json.NewDecoder(response.Body).Decode(&mailbox); err != nil {
		t.Fatal(err)
	}
	for _, message := range mailbox.Messages {
		if len(message.To) != 1 || message.To[0].Address != recipient {
			continue
		}
		body, err := client.Get(base + "/api/v1/message/" + message.ID)
		if err != nil {
			t.Fatal(err)
		}
		var rendered struct{ HTML string }
		err = json.NewDecoder(body.Body).Decode(&rendered)
		body.Body.Close()
		if err != nil || body.StatusCode != http.StatusOK {
			t.Fatal("Mailpit message unavailable", err)
		}
		if strings.Contains(rendered.HTML, secret) {
			return
		}
		// html/template URL-escapes characters such as '+' in the Base64
		// fragment. Compare the decoded link, not the HTML source bytes.
		for _, match := range regexp.MustCompile(`href="([^"]+)"`).FindAllStringSubmatch(rendered.HTML, -1) {
			link, err := url.Parse(html.UnescapeString(match[1]))
			if err == nil && link.Fragment == secret {
				return
			}
		}
	}
	t.Fatal("actual SMTP message did not contain the expected setup link or login number")
}
