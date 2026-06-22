package admin_usecase

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"go.opentelemetry.io/otel"
)

type ListAdminUsers interface {
	Execute(ctx context.Context) ([]admin.AdminUser, error)
}

type listAdminUsersInteractor struct {
	adminUserRepo admin.AdminUserRepository
}

func NewListAdminUsers(adminUserRepo admin.AdminUserRepository) ListAdminUsers {
	return &listAdminUsersInteractor{adminUserRepo: adminUserRepo}
}

func (a *listAdminUsersInteractor) Execute(ctx context.Context) ([]admin.AdminUser, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "listAdminUsersInteractor.Execute")
	defer span.End()

	return a.adminUserRepo.List(ctx)
}
