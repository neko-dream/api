package admin_usecase

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
)

// ErrInvalidAdminCredentials Google 認証に失敗、または許可リストに無い。
var ErrInvalidAdminCredentials = errors.New("invalid admin credentials")

type AdminLogin interface {
	Execute(ctx context.Context, input AdminLoginInput) (*AdminAuthOutput, error)
}

type AdminLoginInput struct {
	// Google Identity Services で取得した ID トークン(credential)。
	IDToken string
}

type adminLoginInteractor struct {
	adminUserRepo     admin.AdminUserRepository
	adminSessionRepo  admin.AdminSessionRepository
	adminTokenManager admin.AdminTokenManager
	googleVerifier    admin.GoogleTokenVerifier
	cfg               *config.Config
}

func NewAdminLogin(
	adminUserRepo admin.AdminUserRepository,
	adminSessionRepo admin.AdminSessionRepository,
	adminTokenManager admin.AdminTokenManager,
	googleVerifier admin.GoogleTokenVerifier,
	cfg *config.Config,
) AdminLogin {
	return &adminLoginInteractor{
		adminUserRepo:     adminUserRepo,
		adminSessionRepo:  adminSessionRepo,
		adminTokenManager: adminTokenManager,
		googleVerifier:    googleVerifier,
		cfg:               cfg,
	}
}

func (a *adminLoginInteractor) Execute(ctx context.Context, input AdminLoginInput) (*AdminAuthOutput, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "adminLoginInteractor.Execute")
	defer span.End()

	claims, err := a.googleVerifier.Verify(ctx, input.IDToken)
	if err != nil {
		return nil, ErrInvalidAdminCredentials
	}

	// 許可リスト(admin_users)に登録された Google アカウントのみログイン可能。
	adminUser, err := a.adminUserRepo.FindByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if adminUser == nil {
		return nil, ErrInvalidAdminCredentials
	}

	return issueAccessAndRefresh(ctx, a.adminTokenManager, a.adminSessionRepo, a.cfg, *adminUser)
}
