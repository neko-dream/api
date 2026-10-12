package admin_usecase

import (
	"context"
	"errors"

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

	session, err := a.adminSessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidRefreshToken
	}
	if session.IsExpired(clock.Now(ctx)) {
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

	// access を再発行し、refresh をローテーション(旧トークンは即時無効化)。
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
	rotated, err := a.adminSessionRepo.Rotate(ctx, session.AdminSessionID, tokenHash, newHash, clock.Now(ctx).Add(refreshTokenTTL))
	if err != nil {
		return nil, err
	}
	if !rotated {
		return nil, ErrInvalidRefreshToken
	}

	return &AdminAuthOutput{AccessToken: access, RefreshToken: newRefresh, AdminUser: *adminUser}, nil
}
