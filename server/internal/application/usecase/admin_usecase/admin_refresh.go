package admin_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/hash"
	"go.opentelemetry.io/otel"
)

// ErrInvalidRefreshToken refresh トークンが無効・失効・期限切れ。
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type AdminRefresh interface {
	Execute(ctx context.Context, refreshToken string) (*AdminAuthOutput, error)
}

type adminRefreshInteractor struct {
	adminUserRepo     admin.AdminUserRepository
	adminSessionRepo  admin.AdminSessionRepository
	adminTokenManager admin.AdminTokenManager
	cfg               *config.Config
}

func NewAdminRefresh(
	adminUserRepo admin.AdminUserRepository,
	adminSessionRepo admin.AdminSessionRepository,
	adminTokenManager admin.AdminTokenManager,
	cfg *config.Config,
) AdminRefresh {
	return &adminRefreshInteractor{
		adminUserRepo:     adminUserRepo,
		adminSessionRepo:  adminSessionRepo,
		adminTokenManager: adminTokenManager,
		cfg:               cfg,
	}
}

func (a *adminRefreshInteractor) Execute(ctx context.Context, refreshToken string) (*AdminAuthOutput, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "adminRefreshInteractor.Execute")
	defer span.End()

	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	tokenHash, err := hashRefreshToken(refreshToken, a.cfg)
	if err != nil {
		return nil, err
	}

	now := clock.Now(ctx)
	session, err := a.adminSessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidRefreshToken
	}
	if session.IsExpired(now) {
		_ = a.adminSessionRepo.DeleteByTokenHash(ctx, tokenHash)
		return nil, ErrInvalidRefreshToken
	}

	adminUser, err := a.adminUserRepo.FindByID(ctx, session.AdminUserID)
	if err != nil {
		return nil, err
	}
	if adminUser == nil {
		return nil, ErrInvalidRefreshToken
	}

	// ローテーション直前のトークン。別タブが先にリフレッシュした場合なので、猶予期間内ならアクセストークンだけ発行する
	if session.TokenHash != tokenHash {
		return a.accessOnly(ctx, session, tokenHash, adminUser, now)
	}

	// access を再発行し、refresh をローテーションする（旧トークンは猶予期間を過ぎると使えなくなる）
	access, err := a.adminTokenManager.Generate(ctx, *adminUser)
	if err != nil {
		return nil, err
	}
	newRefresh, err := hash.GenerateRandomToken(48)
	if err != nil {
		return nil, err
	}
	newHash, err := hashRefreshToken(newRefresh, a.cfg)
	if err != nil {
		return nil, err
	}
	rotated, err := a.adminSessionRepo.Rotate(ctx, session.AdminSessionID, tokenHash, newHash, now.Add(refreshTokenTTL), now)
	if err != nil {
		return nil, err
	}
	if !rotated {
		// 読み込んでからローテーションするまでに、別タブが先にローテーションした
		latest, err := a.adminSessionRepo.FindByTokenHash(ctx, tokenHash)
		if err != nil {
			return nil, err
		}
		if latest == nil {
			return nil, ErrInvalidRefreshToken
		}
		return a.accessOnly(ctx, latest, tokenHash, adminUser, now)
	}

	return &AdminAuthOutput{AccessToken: access, RefreshToken: newRefresh, AdminUser: *adminUser}, nil
}

// accessOnly は猶予期間内の直前のトークンに対して、アクセストークンだけを発行する。
// リフレッシュトークンは先にローテーションしたタブが受け取った新しいものを使うので、返さない（Cookie を上書きしない）。
func (a *adminRefreshInteractor) accessOnly(ctx context.Context, session *admin.AdminSession, tokenHash string, adminUser *admin.AdminUser, now time.Time) (*AdminAuthOutput, error) {
	if !session.IsPreviousTokenInGracePeriod(tokenHash, now) {
		return nil, ErrInvalidRefreshToken
	}
	access, err := a.adminTokenManager.Generate(ctx, *adminUser)
	if err != nil {
		return nil, err
	}
	return &AdminAuthOutput{AccessToken: access, AdminUser: *adminUser}, nil
}
