package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	password_auth "github.com/neko-dream/server/internal/domain/model/auth/password"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/hash"
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

	salt, err := hash.GenerateSalt(16)
	if err != nil {
		return err
	}

	hashedPassword, err := hash.HashPassword(plainPassword, salt, u.cfg.HASH_PEPPER, u.cfg.HASH_ITERATIONS)
	if err != nil {
		return err
	}

	return u.passwordAuthRepo.CreatePasswordAuth(ctx, userID, hashedPassword, salt, requiredChange)
}

// VerifyPassword パスワードを検証
func (u *passwordAuthManager) VerifyPassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) (bool, error) {
	ctx, span := otel.Tracer("usecase").Start(ctx, "passwordAuthUseCase.VerifyPassword")
	defer span.End()

	// ユーザーのパスワード認証情報を取得
	auth, err := u.passwordAuthRepo.GetPasswordAuthByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	log.Println("auth", auth.PasswordHash)
	// パスワードを検証
	bs, err := hash.VerifyPassword(plainPassword, auth.PasswordHash), nil
	log.Println("err", bs, err)
	if err != nil {
		return false, err
	}
	return bs, err
}

// UpdatePassword パスワードを更新
func (u *passwordAuthManager) UpdatePassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) error {
	ctx, span := otel.Tracer("usecase").Start(ctx, "passwordAuthUseCase.UpdatePassword")
	defer span.End()

	salt, err := hash.GenerateSalt(16)
	if err != nil {
		return err
	}

	hashedPassword, err := hash.HashPassword(plainPassword, salt, u.cfg.HASH_PEPPER, u.cfg.HASH_ITERATIONS)
	if err != nil {
		return err
	}

	return u.passwordAuthRepo.UpdatePasswordAuth(ctx, userID, hashedPassword, salt, false)
}
