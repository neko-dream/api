package admin_usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type SetSessionCreationPermission interface {
	Execute(context.Context, shared.UUID[user.User], bool) error
}

type setSessionCreationPermission struct {
	permissions user.SessionCreationPermissionRepository
}

func NewSetSessionCreationPermission(permissions user.SessionCreationPermissionRepository) SetSessionCreationPermission {
	return &setSessionCreationPermission{permissions: permissions}
}

func (s *setSessionCreationPermission) Execute(ctx context.Context, id shared.UUID[user.User], allowed bool) error {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "setSessionCreationPermission.Execute")
	defer span.End()
	claim := admin.GetAdminSession(ctx)
	if claim == nil {
		slog.WarnContext(ctx, "session creation permission change denied", "user_id", id.String(), "allowed", allowed)
		return messages.ForbiddenError
	}
	if err := s.permissions.Set(ctx, id, allowed); err != nil {
		slog.WarnContext(ctx, "session creation permission change failed", "admin_user_id", claim.Sub, "user_id", id.String(), "allowed", allowed)
		if errors.Is(err, sql.ErrNoRows) {
			return messages.UserNotFoundError
		}
		if errors.Is(err, user.ErrSessionCreationPermissionIneligible) {
			return messages.BadRequestError
		}
		return errtrace.Wrap(err)
	}
	slog.InfoContext(ctx, "session creation permission changed", "admin_user_id", claim.Sub, "user_id", id.String(), "allowed", allowed)
	return nil
}
