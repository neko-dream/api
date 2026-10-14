package handler

import (
	"context"

	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/presentation/oas"
)

func (m *manageHandler) CreateUserManage(ctx context.Context, req *oas.InviteManagedUserRequest) (*oas.ManagedUserInvitation, error) {
	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}
	if req == nil {
		return nil, messages.BadRequestError
	}
	out, err := m.userInvitations.Create(ctx, admin_usecase.InviteUserInput{DisplayName: req.DisplayName, Email: req.Email})
	if err != nil {
		return nil, err
	}
	return &oas.ManagedUserInvitation{UserID: out.UserID.String(), InvitationSent: out.InvitationSent}, nil
}

func (m *manageHandler) ResendUserInvitationManage(ctx context.Context, params oas.ResendUserInvitationManageParams) (*oas.ManagedUserInvitation, error) {
	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}
	id, err := shared.ParseUUID[user.User](params.UserID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	out, err := m.userInvitations.Resend(ctx, id)
	if err != nil {
		return nil, err
	}
	return &oas.ManagedUserInvitation{UserID: out.UserID.String(), InvitationSent: out.InvitationSent}, nil
}
