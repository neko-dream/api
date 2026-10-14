package handler

import (
	"context"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

func (m *manageHandler) SetUserPermissionsManage(ctx context.Context, req *oas.UserPermissions, params oas.SetUserPermissionsManageParams) (*oas.UserPermissions, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.SetUserPermissionsManage")
	defer span.End()
	claim := admin.GetAdminSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	id, err := shared.ParseUUID[user.User](params.UserID)
	if err != nil || req == nil {
		return nil, messages.BadRequestError
	}
	if err := m.setCreationPermission.Execute(ctx, id, req.CanCreateTalkSession); err != nil {
		return nil, err
	}
	return &oas.UserPermissions{CanCreateTalkSession: req.CanCreateTalkSession}, nil
}
