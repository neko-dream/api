package admin_usecase

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
)

type AdminLogout interface {
	Execute(ctx context.Context, refreshToken string) error
}

type adminLogoutInteractor struct {
	adminSessionRepo admin.AdminSessionRepository
	cfg              *config.Config
}

func NewAdminLogout(
	adminSessionRepo admin.AdminSessionRepository,
	cfg *config.Config,
) AdminLogout {
	return &adminLogoutInteractor{
		adminSessionRepo: adminSessionRepo,
		cfg:              cfg,
	}
}

func (a *adminLogoutInteractor) Execute(ctx context.Context, refreshToken string) error {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "adminLogoutInteractor.Execute")
	defer span.End()

	if refreshToken == "" {
		return nil
	}
	tokenHash, err := hashRefreshToken(refreshToken, a.cfg)
	if err != nil {
		return err
	}
	return a.adminSessionRepo.DeleteByTokenHash(ctx, tokenHash)
}
