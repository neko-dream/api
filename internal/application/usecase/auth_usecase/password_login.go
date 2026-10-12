package auth_usecase

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	organization_service "github.com/neko-dream/api/internal/domain/service/organization"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/hash"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type PasswordLogin interface {
	Execute(ctx context.Context, input PasswordLoginInput) (*PasswordLoginOutput, error)
	Verify(ctx context.Context, token, browserSecret string) (string, error)
	VerifyPIN(ctx context.Context, pin, browserSecret string) (string, error)
}

const EmailLoginTTL = 15 * time.Minute
const EmailLoginResendInterval = time.Minute
const EmailLoginPINMaxAttempts = 5

type loginTransactionRunner interface {
	ExecTx(context.Context, func(context.Context) error) error
}

type PasswordLoginInput struct {
	IDorEmail string
	Password  string
}

type PasswordLoginOutput struct {
	BrowserSecret string
	ExpiresAt     time.Time
}

type passwordLoginInteractor struct {
	user.UserRepository
	session.SessionService
	session.SessionRepository
	txRunner loginTransactionRunner
	*config.Config
	password_auth.PasswordAuthManager
	passwordRepository         password_auth.PasswordAuthRepository
	challenges                 password_auth.LoginChallengeRepository
	emailSender                email.EmailSender
	organizationRepository     organization.OrganizationRepository
	organizationUserRepository organization.OrganizationUserRepository
	session.TokenManager
}

func NewPasswordLogin(
	userRep user.UserRepository,
	sessService session.SessionService,
	sessionRepository session.SessionRepository,
	tm *db.DBManager,
	config *config.Config,
	passwordAuthManager password_auth.PasswordAuthManager,
	tokenManager session.TokenManager,
	passwordRepository password_auth.PasswordAuthRepository,
	challenges password_auth.LoginChallengeRepository,
	emailSender email.EmailSender,
	organizationRepository organization.OrganizationRepository,
	organizationUserRepository organization.OrganizationUserRepository,
) PasswordLogin {
	return &passwordLoginInteractor{
		UserRepository:             userRep,
		SessionService:             sessService,
		SessionRepository:          sessionRepository,
		txRunner:                   tm,
		Config:                     config,
		PasswordAuthManager:        passwordAuthManager,
		TokenManager:               tokenManager,
		passwordRepository:         passwordRepository,
		challenges:                 challenges,
		emailSender:                emailSender,
		organizationRepository:     organizationRepository,
		organizationUserRepository: organizationUserRepository,
	}
}

func IsEmail(s string) bool {
	// メールアドレスの正規表現
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(s)
}

func (p *passwordLoginInteractor) Execute(ctx context.Context, input PasswordLoginInput) (*PasswordLoginOutput, error) {
	ctx, span := otel.Tracer("auth_command").Start(ctx, "passwordLoginInteractor.Execute")
	defer span.End()

	var output PasswordLoginOutput
	if err := p.txRunner.ExecTx(ctx, func(ctx context.Context) error {
		var usr *user.User
		// IDorEmailがメールアドレスの場合は、ユーザーをメールアドレスで取得
		if IsEmail(input.IDorEmail) {
			emailHash, err := hash.HMACHash(input.IDorEmail, p.Config.HASH_PEPPER)
			if err != nil {
				utils.HandleError(ctx, err, "failed to hash email")
				return messages.InvalidPasswordOrEmailError
			}
			foundUser, err := p.UserRepository.FindBySubject(ctx, user.UserSubject(emailHash))
			if err != nil {
				return messages.InvalidPasswordOrEmailError
			}
			usr = foundUser
		} else {
			// DisplayIDからユーザーを取得
			foundUser, err := p.UserRepository.FindByDisplayID(ctx, input.IDorEmail)
			if err != nil {
				return messages.InvalidPasswordOrEmailError
			}
			usr = foundUser
		}
		if usr == nil {
			return messages.InvalidPasswordOrEmailError
		}
		if usr.IsReactivationPeriodExpired(clock.Now(ctx)) || usr.Email() == nil || *usr.Email() == "" {
			return messages.InvalidPasswordOrEmailError
		}

		// ユーザーからユーザーパスワードを取得
		userPassword, err := p.PasswordAuthManager.VerifyPassword(ctx, usr.UserID(), input.Password)
		if err != nil {
			return messages.InvalidPasswordOrEmailError
		}
		// パスワードが一致しない場合はエラーを返す
		if !userPassword {
			return messages.InvalidPasswordOrEmailError
		}

		credential, err := p.passwordRepository.GetPasswordAuthByUserID(ctx, usr.UserID())
		if err != nil {
			return err
		}
		if p.Config.TokenSecret == "" {
			return fmt.Errorf("TOKEN_SECRET is required for email login PIN protection")
		}
		pin, err := newLoginPIN()
		if err != nil {
			return err
		}
		token, err := newLoginSecret()
		if err != nil {
			return err
		}
		browser, err := newLoginSecret()
		if err != nil {
			return err
		}
		verificationURL, err := p.verificationURL(token)
		if err != nil {
			return err
		}
		now := clock.Now(ctx)
		expires := now.Add(EmailLoginTTL)
		browserHash := loginHash(browser)
		challenge := password_auth.LoginChallenge{UserID: usr.UserID(), TokenHash: loginHash(token), BrowserHash: browserHash, CredentialHash: credentialHash(credential, usr), PINHash: p.pinHash(pin, browserHash), CreatedAt: now, ExpiresAt: expires}
		if err := p.challenges.Save(ctx, challenge, now.Add(-EmailLoginResendInterval)); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return messages.EmailLoginRateLimitError
			}
			return err
		}
		if err := p.emailSender.Send(ctx, *usr.Email(), email_template.PasswordLoginEmailTemplate, map[string]any{"Title": "ログインの確認", "VerificationURL": verificationURL, "ExpiryMinutes": 15, "PIN": pin}); err != nil {
			return err
		}
		output = PasswordLoginOutput{BrowserSecret: browser, ExpiresAt: expires}
		return nil
	}); err != nil {
		return nil, err
	}

	return &output, nil
}

// Verify consumes the challenge and creates the session in the same transaction.
func (p *passwordLoginInteractor) Verify(ctx context.Context, token, browserSecret string) (string, error) {
	if !validLoginSecret(token) || !validLoginSecret(browserSecret) {
		return "", messages.InvalidEmailLoginLinkError
	}
	var sessionToken string
	err := p.txRunner.ExecTx(ctx, func(ctx context.Context) error {
		challenge, err := p.challenges.Consume(ctx, loginHash(token), loginHash(browserSecret), clock.Now(ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return messages.InvalidEmailLoginLinkError
		}
		if err != nil {
			return err
		}
		sessionToken, err = p.completeEmailLogin(ctx, challenge)
		return err
	})
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

// Failed guesses must be committed even though the HTTP response is an error.
// SELECT FOR UPDATE serializes PIN guesses and link/PIN consumption.
func (p *passwordLoginInteractor) VerifyPIN(ctx context.Context, pin, browserSecret string) (string, error) {
	if !validLoginSecret(browserSecret) {
		return "", messages.InvalidEmailLoginPINError
	}
	var sessionToken string
	var verificationError error
	err := p.txRunner.ExecTx(ctx, func(ctx context.Context) error {
		browserHash := loginHash(browserSecret)
		challenge, err := p.challenges.FindForPIN(ctx, browserHash)
		if errors.Is(err, sql.ErrNoRows) {
			verificationError = messages.InvalidEmailLoginPINError
			return nil
		}
		if err != nil {
			return err
		}
		if challenge.Consumed || !clock.Now(ctx).Before(challenge.ExpiresAt) || len(challenge.PINHash) != sha256.Size {
			verificationError = messages.InvalidEmailLoginPINError
			return nil
		}
		if challenge.PINAttempts >= EmailLoginPINMaxAttempts {
			verificationError = messages.EmailLoginPINAttemptsError
			return nil
		}
		if !validLoginPIN(pin) || !hmac.Equal(challenge.PINHash, p.pinHash(pin, browserHash)) {
			attempts, err := p.challenges.RecordPINFailure(ctx, browserHash)
			if err != nil {
				return err
			}
			verificationError = messages.InvalidEmailLoginPINError
			if attempts >= EmailLoginPINMaxAttempts {
				verificationError = messages.EmailLoginPINAttemptsError
			}
			return nil
		}
		consumed, err := p.challenges.Consume(ctx, challenge.TokenHash, browserHash, clock.Now(ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return messages.InvalidEmailLoginPINError
		}
		if err != nil {
			return err
		}
		sessionToken, err = p.completeEmailLogin(ctx, consumed)
		return err
	})
	if err != nil {
		return "", err
	}
	if verificationError != nil {
		return "", verificationError
	}
	return sessionToken, nil
}

func (p *passwordLoginInteractor) completeEmailLogin(ctx context.Context, challenge password_auth.LoginChallenge) (string, error) {
	usr, err := p.UserRepository.FindByID(ctx, challenge.UserID)
	if err != nil {
		return "", err
	}
	if usr == nil || usr.IsReactivationPeriodExpired(clock.Now(ctx)) || usr.Email() == nil {
		return "", messages.InvalidEmailLoginLinkError
	}
	credential, err := p.passwordRepository.GetPasswordAuthByUserID(ctx, challenge.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", messages.InvalidEmailLoginLinkError
	}
	if err != nil {
		return "", err
	}
	if subtle.ConstantTimeCompare(challenge.CredentialHash, credentialHash(credential, usr)) != 1 {
		return "", messages.InvalidEmailLoginLinkError
	}
	usr.SetEmailVerified(true)
	if err := p.UserRepository.Update(ctx, *usr); err != nil {
		return "", err
	}
	orgID, err := organization_service.ResolveLoginOrganization(ctx, p.organizationUserRepository, p.organizationRepository, usr.UserID(), nil)
	if err != nil {
		return "", err
	}
	// Existing password-login session policy is retained.
	if err := p.SessionRepository.DeactivateAllByUserID(ctx, usr.UserID()); err != nil {
		return "", err
	}
	sess := session.NewSessionWithOrganization(shared.NewUUID[session.Session](), usr.UserID(), usr.Provider(), session.SESSION_ACTIVE, *session.NewExpiresAt(ctx), clock.Now(ctx), orgID)
	if _, err := p.SessionRepository.Create(ctx, *sess); err != nil {
		return "", err
	}
	return p.TokenManager.Generate(ctx, *usr, sess.SessionID())
}

func newLoginPIN() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", number.Int64()), nil
}
func validLoginPIN(pin string) bool {
	if len(pin) != 6 {
		return false
	}
	for _, digit := range pin {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}
func (p *passwordLoginInteractor) pinHash(pin string, browserHash []byte) []byte {
	// Six-digit PINs are cheap to brute-force: a server-held key protects DB-only leaks.
	mac := hmac.New(sha256.New, []byte(p.Config.TokenSecret))
	mac.Write([]byte("password-email-pin\x00"))
	mac.Write(browserHash)
	mac.Write([]byte(pin))
	return mac.Sum(nil)
}

func newLoginSecret() (string, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(secret[:]), nil
}

func validLoginSecret(secret string) bool {
	if len(secret) != 43 {
		return false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(secret)
	return err == nil && len(decoded) == 32 && base64.RawURLEncoding.EncodeToString(decoded) == secret
}

func loginHash(value string) []byte { sum := sha256.Sum256([]byte(value)); return sum[:] }
func credentialHash(auth password_auth.PasswordAuth, usr *user.User) []byte {
	return loginHash(auth.PasswordHash + "\x00" + auth.Salt + "\x00" + *usr.Email())
}

func (p *passwordLoginInteractor) verificationURL(token string) (string, error) {
	if page := p.Config.EMAIL_VERIFICATION_PAGE_URL; page != "" {
		u, err := url.Parse(page)
		if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
			return "", fmt.Errorf("EMAIL_VERIFICATION_PAGE_URL must be an absolute HTTPS URL without credentials, query or fragment")
		}
		u.Fragment = token
		return u.String(), nil
	}
	base := p.Config.EMAIL_VERIFICATION_BASE_URL
	if base == "" && p.Config.Env == config.LOCAL {
		base = "https://localhost:3000/api"
	}
	u, err := url.Parse(base)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("EMAIL_VERIFICATION_BASE_URL must be an absolute HTTPS URL without credentials, query or fragment")
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/auth/password/verify"
	// Fragments stay in the browser: the token never enters server access logs.
	u.Fragment = token
	return u.String(), nil
}
