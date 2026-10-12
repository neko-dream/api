package repository

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type sessionCreationPermissionRepository struct{ *db.DBManager }

func NewSessionCreationPermissionRepository(dbm *db.DBManager) user.SessionCreationPermissionRepository {
	return &sessionCreationPermissionRepository{DBManager: dbm}
}

func (r *sessionCreationPermissionRepository) Get(ctx context.Context, id shared.UUID[user.User]) (bool, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionCreationPermissionRepository.Get")
	defer span.End()
	return r.GetQueries(ctx).GetSessionCreationPermission(ctx, id.UUID())
}

func (r *sessionCreationPermissionRepository) Set(ctx context.Context, id shared.UUID[user.User], allowed bool) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionCreationPermissionRepository.Set")
	defer span.End()
	_, err := r.GetQueries(ctx).SetSessionCreationPermission(ctx, model.SetSessionCreationPermissionParams{UserID: id.UUID(), CanCreateTalkSession: allowed})
	return err
}
