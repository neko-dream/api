package admin_usecase

import (
	"context"
	"database/sql"
	"errors"

	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type SetSessionCreationPermission interface {
	Execute(context.Context, shared.UUID[user.User], bool) error
}

type setSessionCreationPermission struct {
	users       userquery.Current
	permissions user.SessionCreationPermissionRepository
}

func NewSetSessionCreationPermission(users userquery.Current, permissions user.SessionCreationPermissionRepository) SetSessionCreationPermission {
	return &setSessionCreationPermission{users: users, permissions: permissions}
}

func (s *setSessionCreationPermission) Execute(ctx context.Context, id shared.UUID[user.User], allowed bool) error {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "setSessionCreationPermission.Execute")
	defer span.End()
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return messages.UserNotFoundError
	}
	if allowed && !user.CanCreateTalkSession(true, u.IsGuest, u.User.WithdrawalDate != nil) {
		return messages.BadRequestError
	}
	if err := s.permissions.Set(ctx, id, allowed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages.UserNotFoundError
		}
		return err
	}
	return nil
}
