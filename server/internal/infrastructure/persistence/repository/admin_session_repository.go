package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type adminSessionRepository struct {
	*db.DBManager
}

func NewAdminSessionRepository(dbManager *db.DBManager) admin.AdminSessionRepository {
	return &adminSessionRepository{DBManager: dbManager}
}

func (r *adminSessionRepository) Create(ctx context.Context, session admin.AdminSession) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminSessionRepository.Create")
	defer span.End()

	if err := r.DBManager.GetQueries(ctx).CreateAdminSession(ctx, model.CreateAdminSessionParams{
		AdminSessionID: session.AdminSessionID.UUID(),
		AdminUserID:    session.AdminUserID.UUID(),
		TokenHash:      session.TokenHash,
		ExpiresAt:      session.ExpiresAt,
	}); err != nil {
		utils.HandleError(ctx, err, "CreateAdminSession")
		return err
	}
	return nil
}

func (r *adminSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*admin.AdminSession, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminSessionRepository.FindByTokenHash")
	defer span.End()

	row, err := r.DBManager.GetQueries(ctx).GetAdminSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.HandleError(ctx, err, "GetAdminSessionByTokenHash")
		return nil, err
	}
	return &admin.AdminSession{
		AdminSessionID: shared.UUID[admin.AdminSession](row.AdminSessionID),
		AdminUserID:    shared.UUID[admin.AdminUser](row.AdminUserID),
		TokenHash:      row.TokenHash,
		ExpiresAt:      row.ExpiresAt,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}, nil
}

func (r *adminSessionRepository) Rotate(ctx context.Context, id shared.UUID[admin.AdminSession], newTokenHash string, expiresAt time.Time) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminSessionRepository.Rotate")
	defer span.End()

	if err := r.DBManager.GetQueries(ctx).RotateAdminSession(ctx, model.RotateAdminSessionParams{
		AdminSessionID: id.UUID(),
		TokenHash:      newTokenHash,
		ExpiresAt:      expiresAt,
	}); err != nil {
		utils.HandleError(ctx, err, "RotateAdminSession")
		return err
	}
	return nil
}

func (r *adminSessionRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "adminSessionRepository.DeleteByTokenHash")
	defer span.End()

	if err := r.DBManager.GetQueries(ctx).DeleteAdminSessionByTokenHash(ctx, tokenHash); err != nil {
		utils.HandleError(ctx, err, "DeleteAdminSessionByTokenHash")
		return err
	}
	return nil
}
