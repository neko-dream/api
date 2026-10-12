package organization_usecase

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	organization_svc "github.com/neko-dream/api/internal/domain/service/organization"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type InviteOrganizationCommand interface {
	Execute(ctx context.Context, input InviteOrganizationInput) (*InviteOrganizationOutput, error)
}
type InviteOrganizationInput struct {
	UserID         shared.UUID[user.User]
	OrganizationID shared.UUID[organization.Organization]
	Role           int
	Email          string
}

type InviteOrganizationOutput struct {
	Success bool
}

type inviteOrganizationInteractor struct {
	organizationService        organization_svc.OrganizationService
	userRepository             user.UserRepository
	organizationUserRepository organization.OrganizationUserRepository
	organization_svc.OrganizationMemberManager
	emailSender email.EmailSender
	cfg         *config.Config
}

func NewInviteOrganizationInteractor(
	organizationService organization_svc.OrganizationService,
	userRepository user.UserRepository,
	organizationUserRepository organization.OrganizationUserRepository,
	organizationMemberManager organization_svc.OrganizationMemberManager,
	emailSender email.EmailSender,
	cfg *config.Config,
) InviteOrganizationCommand {
	return &inviteOrganizationInteractor{
		organizationService:        organizationService,
		userRepository:             userRepository,
		organizationUserRepository: organizationUserRepository,
		OrganizationMemberManager:  organizationMemberManager,
		emailSender:                emailSender,
		cfg:                        cfg,
	}
}

func (i *inviteOrganizationInteractor) Execute(ctx context.Context, input InviteOrganizationInput) (*InviteOrganizationOutput, error) {
	ctx, span := otel.Tracer("organization_command").Start(ctx, "inviteOrganizationInteractor.Execute")
	defer span.End()

	// ユーザーの組織ユーザーを取得
	orgUser, err := i.organizationUserRepository.FindByOrganizationIDAndUserID(ctx, input.OrganizationID, input.UserID)
	if err != nil {
		return nil, err
	}
	// 組織ユーザーが存在しない場合はエラー
	if orgUser == nil {
		return nil, errors.New("organization user not found")
	}

	// 有効なロール値かチェック (10, 20, 30, 40)
	validRoles := []int{10, 20, 30, 40}
	isValidRole := false
	for _, validRole := range validRoles {
		if input.Role == validRole {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return nil, errors.New("invalid role")
	}

	// 組織の招待を送信
	_, err = i.OrganizationMemberManager.InviteUser(ctx, organization_svc.InviteUserParams{
		OrganizationID: input.OrganizationID,
		Role:           organization.OrganizationUserRole(input.Role),
		UserID:         input.UserID,
		Email:          input.Email,
	})
	if err != nil {
		utils.HandleError(ctx, err, "InviteUser")
		return nil, err
	}

	return &InviteOrganizationOutput{
		Success: true,
	}, nil
}
