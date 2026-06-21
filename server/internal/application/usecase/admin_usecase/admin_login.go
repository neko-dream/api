package admin_usecase

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/hash"
	"go.opentelemetry.io/otel"
)

// ErrInvalidAdminCredentials メールアドレスまたはパスワードが正しくない。
var ErrInvalidAdminCredentials = errors.New("invalid admin credentials")

type AdminLogin interface {
	Execute(ctx context.Context, input AdminLoginInput) (*AdminLoginOutput, error)
}

type AdminLoginInput struct {
	Email    string
	Password string
}

type AdminLoginOutput struct {
	Token     string
	AdminUser admin.AdminUser
}

type adminLoginInteractor struct {
	adminUserRepo     admin.AdminUserRepository
	adminTokenManager admin.AdminTokenManager
	cfg               *config.Config
}

func NewAdminLogin(
	adminUserRepo admin.AdminUserRepository,
	adminTokenManager admin.AdminTokenManager,
	cfg *config.Config,
) AdminLogin {
	return &adminLoginInteractor{
		adminUserRepo:     adminUserRepo,
		adminTokenManager: adminTokenManager,
		cfg:               cfg,
	}
}

func (a *adminLoginInteractor) Execute(ctx context.Context, input AdminLoginInput) (*AdminLoginOutput, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "adminLoginInteractor.Execute")
	defer span.End()

	adminUser, err := a.adminUserRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if adminUser == nil {
		return nil, ErrInvalidAdminCredentials
	}

	if !hash.BcryptVerify(input.Password, a.cfg.HASH_PEPPER, adminUser.PasswordHash) {
		return nil, ErrInvalidAdminCredentials
	}

	token, err := a.adminTokenManager.Generate(ctx, *adminUser)
	if err != nil {
		return nil, err
	}

	return &AdminLoginOutput{
		Token:     token,
		AdminUser: *adminUser,
	}, nil
}
