package auth_command

import (
	"context"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	password_auth "github.com/neko-dream/server/internal/domain/model/auth/password"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/hash"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type PasswordRegister interface {
	Execute(ctx context.Context, input PasswordRegisterInput) (*PasswordRegisterOutput, error)
}
type PasswordRegisterInput struct {
	Email    string
	Password string
}

type PasswordRegisterOutput struct {
	Token string
}

type passwordRegisterInteractor struct {
	userRepository      user.UserRepository
	consentService      consent.ConsentService
	policyRepository    consent.PolicyRepository
	cfg                 *config.Config
	passwordAuthManager password_auth.PasswordAuthManager
	session.TokenManager
	session.SessionRepository

	*db.DBManager
}

func NewPasswordRegister(
	userRepository user.UserRepository,
	consentService consent.ConsentService,
	policyRepository consent.PolicyRepository,
	cfg *config.Config,
	passwordAuthManager password_auth.PasswordAuthManager,
	tokenManager session.TokenManager,
	sessionRepository session.SessionRepository,
	dbManager *db.DBManager,
) PasswordRegister {
	return &passwordRegisterInteractor{
		userRepository:      userRepository,
		consentService:      consentService,
		policyRepository:    policyRepository,
		passwordAuthManager: passwordAuthManager,
		cfg:                 cfg,
		TokenManager:        tokenManager,
		SessionRepository:   sessionRepository,
		DBManager:           dbManager,
	}
}

func (p *passwordRegisterInteractor) Execute(ctx context.Context, input PasswordRegisterInput) (*PasswordRegisterOutput, error) {
	ctx, span := otel.Tracer("auth_command").Start(ctx, "passwordRegisterInteractor.Execute")
	defer span.End()
	if p.cfg.Env != config.PROD {
		return nil, messages.InternalServerError
	}

	var tokenRes string
	if err := p.ExecTx(ctx, func(ctx context.Context) error {
		authProviderName, err := auth.NewAuthProviderName("password")
		if err != nil {
			return errtrace.Wrap(err)
		}
		subject, err := hash.HashEmail(input.Email, p.cfg.HASH_PEPPER)
		if err != nil {
			utils.HandleError(ctx, err, "HashEmail")
			return messages.InvalidPasswordOrEmailError
		}
		existUser, err := p.userRepository.FindBySubject(ctx, user.UserSubject(subject))
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindBySubject")
			return errtrace.Wrap(err)
		}
		if existUser != nil {
			return errors.New("既に登録済みです。")
		}

		newUser := user.NewUser(
			shared.NewUUID[user.User](),
			nil,
			nil,
			subject,
			authProviderName,
			nil,
		)
		newUser.ChangeEmail(input.Email)
		version, err := p.policyRepository.FetchLatestPolicy(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "PolicyRepository.GetLatestVersion")
			return errtrace.Wrap(err)
		}
		_, err = p.consentService.RecordConsent(
			ctx,
			newUser.UserID(),
			version.Version,
			"",
			"",
		)
		if err != nil {
			utils.HandleError(ctx, err, "ConsentService.RecordConsent")
			return errtrace.Wrap(err)
		}
		if err := p.userRepository.Create(ctx, newUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Create")
			return errtrace.Wrap(err)
		}
		if err := p.passwordAuthManager.RegisterPassword(ctx, newUser.UserID(), input.Password, false); err != nil {
			utils.HandleError(ctx, err, "PasswordAuthManager.RegisterPassword")
			return errtrace.Wrap(err)
		}

		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			newUser.UserID(),
			newUser.Provider(),
			session.SESSION_ACTIVE,
			*session.NewExpiresAt(ctx),
			clock.Now(ctx),
		)

		if _, err := p.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "failed to create session")
			return errtrace.Wrap(err)
		}

		token, err := p.TokenManager.Generate(ctx, newUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return errtrace.Wrap(err)
		}

		tokenRes = token
		return nil
	}); err != nil {
		return nil, err
	}

	return &PasswordRegisterOutput{
		Token: tokenRes,
	}, nil
}
