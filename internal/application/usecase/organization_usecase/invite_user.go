package organization_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	organization_svc "github.com/neko-dream/api/internal/domain/service/organization"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type InviteOrganizationForUserCommand interface {
	Execute(ctx context.Context, input InviteOrganizationForUserInput) (*InviteOrganizationForUserOutput, error)
}

type InviteOrganizationForUserInput struct {
	UserID         shared.UUID[user.User]
	OrganizationID shared.UUID[organization.Organization]
	DisplayID      string
	Role           int
}

type InviteOrganizationForUserOutput struct {
	Success bool
}

type inviteOrganizationForUserInteractor struct {
	userRepository             user.UserRepository
	organizationUserRepository organization.OrganizationUserRepository
	organizationMemberManager  organization_svc.OrganizationMemberManager
}

func NewInviteOrganizationForUserInteractor(
	userRepository user.UserRepository,
	organizationUserRepository organization.OrganizationUserRepository,
	organizationMemberManager organization_svc.OrganizationMemberManager,
) InviteOrganizationForUserCommand {
	return &inviteOrganizationForUserInteractor{
		userRepository:             userRepository,
		organizationUserRepository: organizationUserRepository,
		organizationMemberManager:  organizationMemberManager,
	}
}

func (i *inviteOrganizationForUserInteractor) Execute(ctx context.Context, input InviteOrganizationForUserInput) (*InviteOrganizationForUserOutput, error) {
	ctx, span := otel.Tracer("organization_command").Start(ctx, "inviteOrganizationForUserInteractor.Execute")
	defer span.End()

	// DisplayIDからユーザーを取得
	user, err := i.userRepository.FindByDisplayID(ctx, input.DisplayID)
	if err != nil {
		utils.HandleError(ctx, err, "userRepository.FindByDisplayID")
		return nil, messages.UserNotFound
	}
	if user == nil {
		return nil, messages.UserNotFound
	}

	// ログインユーザーが組織の管理者であることを確認
	orgUser, err := i.organizationUserRepository.FindByOrganizationIDAndUserID(ctx, input.OrganizationID, input.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, messages.OrganizationPermissionDenied
		}
		return nil, messages.OrganizationInternalServerError
	}

	// ドメインロジックを使用して権限チェックを行うのじゃ
	targetRole := organization.NewOrganizationUserRole(input.Role)
	if !orgUser.HasPermissionToChangeRoleTo(targetRole) {
		return nil, messages.OrganizationPermissionDenied
	}

	if err := i.organizationMemberManager.AddUser(ctx, organization_svc.InviteUserParams{
		OrganizationID: input.OrganizationID,
		Role:           targetRole,
		UserID:         user.UserID(),
	}); err != nil {
		utils.HandleError(ctx, err, "organizationMemberManager.AddUser")
		return nil, err
	}

	return &InviteOrganizationForUserOutput{Success: true}, nil
}
