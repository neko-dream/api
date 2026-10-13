package service

import (
	"context"
	"database/sql"
	"errors"

	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/hash"
	"go.opentelemetry.io/otel"
)

type passwordAuthManager struct {
	passwordAuthRepo password_auth.PasswordAuthRepository
	cfg              *config.Config
}

func NewPasswordAuthManager(
	repo password_auth.PasswordAuthRepository,
	cfg *config.Config,
) password_auth.PasswordAuthManager {
	return &passwordAuthManager{
		passwordAuthRepo: repo,
		cfg:              cfg,
	}
}

// RegisterPassword 新しいパスワードを登録
func (u *passwordAuthManager) RegisterPassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string, requiredChange bool) error {
	ctx, span := otel.Tracer("usecase").Start(ctx, "passwordAuthUseCase.RegisterPassword")
	defer span.End()

	policy, err := password_auth.NewPasswordPolicy(len(u.cfg.HASH_PEPPER))
	if err != nil {
		return err
	}
	if err := policy.Validate(plainPassword); err != nil {
		return err
	}

	salt, err := hash.GenerateRandomToken(16)
	if err != nil {
		return err
	}

	hashedPassword, err := hash.BcryptHash(plainPassword, u.cfg.HASH_PEPPER, u.cfg.HASH_ITERATIONS)
	if err != nil {
		return err
	}

	return u.passwordAuthRepo.CreatePasswordAuth(ctx, userID, hashedPassword, salt, requiredChange)
}

// VerifyPassword パスワードを検証
func (u *passwordAuthManager) VerifyPassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) (bool, error) {
	ctx, span := otel.Tracer("usecase").Start(ctx, "passwordAuthUseCase.VerifyPassword")
	defer span.End()

	auth, err := u.passwordAuthRepo.GetPasswordAuthByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return hash.BcryptVerify(plainPassword, u.cfg.HASH_PEPPER, auth.PasswordHash), nil
}

// UpdatePassword パスワードを更新
func (u *passwordAuthManager) UpdatePassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) error {
	ctx, span := otel.Tracer("usecase").Start(ctx, "passwordAuthUseCase.UpdatePassword")
	defer span.End()

	policy, err := password_auth.NewPasswordPolicy(len(u.cfg.HASH_PEPPER))
	if err != nil {
		return err
	}
	if err := policy.Validate(plainPassword); err != nil {
		return err
	}

	salt, err := hash.GenerateRandomToken(16)
	if err != nil {
		return err
	}

	hashedPassword, err := hash.BcryptHash(plainPassword, u.cfg.HASH_PEPPER, u.cfg.HASH_ITERATIONS)
	if err != nil {
		return err
	}

	return u.passwordAuthRepo.UpdatePasswordAuth(ctx, userID, hashedPassword, salt, false)
}
