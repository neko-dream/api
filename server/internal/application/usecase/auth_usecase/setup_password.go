package auth_usecase

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/hash"
	"go.opentelemetry.io/otel"
)

type SetupPassword interface {
	Execute(context.Context, SetupPasswordInput) error
}
type SetupPasswordInput struct{ Token, Password string }

type setupPassword struct {
	users     user.UserRepository
	passwords password_auth.PasswordAuthRepository
	manager   password_auth.PasswordAuthManager
	tokens    password_auth.SetupTokenRepository
	cfg       *config.Config
	tx        loginTransactionRunner
}

func NewSetupPassword(users user.UserRepository, passwords password_auth.PasswordAuthRepository, manager password_auth.PasswordAuthManager, tokens password_auth.SetupTokenRepository, cfg *config.Config, tx *db.DBManager) SetupPassword {
	return &setupPassword{users: users, passwords: passwords, manager: manager, tokens: tokens, cfg: cfg, tx: tx}
}

func (s *setupPassword) Execute(ctx context.Context, input SetupPasswordInput) error {
	ctx, span := otel.Tracer("auth_command").Start(ctx, "setupPassword.Execute")
	defer span.End()
	digest, valid := password_auth.SetupTokenHash(input.Token)
	if !valid {
		return messages.InvalidPasswordSetupLink
	}
	policy, err := password_auth.NewPasswordPolicy(len(s.cfg.HASH_PEPPER))
	if err != nil {
		return err
	}
	if err := policy.Validate(input.Password); err != nil {
		return err
	}
	return s.tx.ExecTx(ctx, func(ctx context.Context) error {
		id, err := s.tokens.FindUserIDByTokenHash(ctx, digest)
		if errors.Is(err, sql.ErrNoRows) {
			return messages.InvalidPasswordSetupLink
		}
		if err != nil {
			return err
		}
		if err := s.tokens.LockUser(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return messages.InvalidPasswordSetupLink
			}
			return err
		}
		pending, err := s.tokens.FindByUserID(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return messages.InvalidPasswordSetupLink
		}
		if err != nil {
			return err
		}
		if subtle.ConstantTimeCompare(pending.TokenHash, digest) != 1 || !clock.Now(ctx).Before(pending.ExpiresAt) {
			return messages.InvalidPasswordSetupLink
		}
		usr, err := s.users.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if usr == nil || usr.IsWithdrawn() || usr.Provider() != "PASSWORD" || usr.Email() == nil || usr.Subject() != pending.EmailSubject {
			return messages.InvalidPasswordSetupLink
		}
		subject, err := hash.HMACHash(*usr.Email(), s.cfg.HASH_PEPPER)
		if err != nil {
			return err
		}
		if subject != pending.EmailSubject {
			return messages.InvalidPasswordSetupLink
		}
		_, err = s.passwords.GetPasswordAuthByUserID(ctx, id)
		if err == nil {
			return messages.InvalidPasswordSetupLink
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err := s.manager.RegisterPassword(ctx, id, input.Password, false); err != nil {
			return err
		}
		usr.SetEmailVerified(true)
		if err := s.users.Update(ctx, *usr); err != nil {
			return err
		}
		// Consumption, credential creation and email verification share one transaction.
		return s.tokens.Delete(ctx, id)
	})
}
