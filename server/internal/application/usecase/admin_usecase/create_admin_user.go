package admin_usecase

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"go.opentelemetry.io/otel"
)

// ErrAdminUserAlreadyExists 同じ email の管理者が既に存在する。
var ErrAdminUserAlreadyExists = errors.New("admin user already exists")

type CreateAdminUserInput struct {
	Email       string
	DisplayName string
}

type CreateAdminUser interface {
	Execute(ctx context.Context, input CreateAdminUserInput) (*admin.AdminUser, error)
}

type createAdminUserInteractor struct {
	adminUserRepo admin.AdminUserRepository
}

func NewCreateAdminUser(adminUserRepo admin.AdminUserRepository) CreateAdminUser {
	return &createAdminUserInteractor{adminUserRepo: adminUserRepo}
}

func (a *createAdminUserInteractor) Execute(ctx context.Context, input CreateAdminUserInput) (*admin.AdminUser, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "createAdminUserInteractor.Execute")
	defer span.End()

	existing, err := a.adminUserRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAdminUserAlreadyExists
	}

	displayName := input.DisplayName
	if displayName == "" {
		displayName = input.Email
	}

	adminUser := admin.AdminUser{
		AdminUserID: shared.NewUUID[admin.AdminUser](),
		Email:       input.Email,
		DisplayName: displayName,
	}
	if err := a.adminUserRepo.Create(ctx, adminUser); err != nil {
		return nil, err
	}
	return &adminUser, nil
}
