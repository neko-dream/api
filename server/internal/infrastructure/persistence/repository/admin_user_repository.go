package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type adminUserRepository struct {
	*db.DBManager
}

func NewAdminUserRepository(dbManager *db.DBManager) admin.AdminUserRepository {
	return &adminUserRepository{DBManager: dbManager}
}

func (r *adminUserRepository) Create(ctx context.Context, adminUser admin.AdminUser) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminUserRepository.Create")
	defer span.End()

	if err := r.DBManager.GetQueries(ctx).CreateAdminUser(ctx, model.CreateAdminUserParams{
		AdminUserID: adminUser.AdminUserID.UUID(),
		Email:       adminUser.Email,
		DisplayName: adminUser.DisplayName,
		Role:        int32(adminUser.Role),
	}); err != nil {
		utils.HandleError(ctx, err, "CreateAdminUser")
		return err
	}
	return nil
}

func (r *adminUserRepository) FindByEmail(ctx context.Context, email string) (*admin.AdminUser, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminUserRepository.FindByEmail")
	defer span.End()

	row, err := r.DBManager.GetQueries(ctx).GetAdminUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.HandleError(ctx, err, "GetAdminUserByEmail")
		return nil, err
	}
	return toAdminUser(row), nil
}

func (r *adminUserRepository) FindByID(ctx context.Context, id shared.UUID[admin.AdminUser]) (*admin.AdminUser, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminUserRepository.FindByID")
	defer span.End()

	row, err := r.DBManager.GetQueries(ctx).GetAdminUserByID(ctx, id.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.HandleError(ctx, err, "GetAdminUserByID")
		return nil, err
	}
	return toAdminUser(row), nil
}

func (r *adminUserRepository) List(ctx context.Context) ([]admin.AdminUser, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminUserRepository.List")
	defer span.End()

	rows, err := r.DBManager.GetQueries(ctx).ListAdminUsers(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "ListAdminUsers")
		return nil, err
	}
	users := make([]admin.AdminUser, 0, len(rows))
	for _, row := range rows {
		users = append(users, *toAdminUser(row))
	}
	return users, nil
}

func toAdminUser(row model.AdminUser) *admin.AdminUser {
	return &admin.AdminUser{
		AdminUserID: shared.UUID[admin.AdminUser](row.AdminUserID),
		Email:       row.Email,
		DisplayName: row.DisplayName,
		Role:        admin.AdminRole(row.Role),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
