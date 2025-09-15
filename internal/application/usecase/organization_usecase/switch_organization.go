package organization_usecase

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	organization_svc "github.com/neko-dream/api/internal/domain/service/organization"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type SwitchOrganizationUseCaseInput struct {
	Code      string
	UserID    shared.UUID[user.User]
	SessionID shared.UUID[session.Session]
}

type SwitchOrganizationUseCaseOutput struct {
	SessionTokenStr string
}

type SwitchOrganizationUseCase interface {
	Execute(ctx context.Context, input SwitchOrganizationUseCaseInput) (*SwitchOrganizationUseCaseOutput, error)
}

type switchOrganizationInteractor struct {
	organizationService  organization_svc.OrganizationService
	organizationUserRepo organization.OrganizationUserRepository
	sessionService       session.SessionService
	tokenManager         session.TokenManager
	userRepository       user.UserRepository
}

// Execute 組織を切り替える
func (s *switchOrganizationInteractor) Execute(
	ctx context.Context,
	input SwitchOrganizationUseCaseInput,
) (*SwitchOrganizationUseCaseOutput, error) {
	ctx, span := otel.Tracer("usecase").Start(ctx, "switchOrganizationInteractor.Execute")
	defer span.End()

	// 組織コードから組織IDを解決
	orgIDAny, err := s.organizationService.ResolveOrganizationIDFromCode(ctx, &input.Code)
	if err != nil {
		utils.HandleError(ctx, err, "organizationService.ResolveOrganizationIDFromCode")
		return nil, errtrace.Wrap(err)
	}
	if orgIDAny == nil {
		return nil, messages.OrganizationPermissionDenied
	}

	// UUID[any]をUUID[organization.Organization]に変換
	orgID := shared.UUID[organization.Organization](orgIDAny.UUID())

	// ユーザーがその組織に所属しているか確認
	_, err = s.organizationUserRepo.FindByOrganizationIDAndUserID(ctx, orgID, input.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, messages.OrganizationPermissionDenied
		}
		utils.HandleError(ctx, err, "organizationUserRepo.FindByOrganizationIDAndUserID")
		return nil, errtrace.Wrap(err)
	}

	// 現在のセッションを無効化し、新しい組織付きセッションを作成
	newSession, err := s.sessionService.SwitchOrganization(ctx, input.UserID, orgID, input.SessionID)
	if err != nil {
		utils.HandleError(ctx, err, "sessionService.SwitchOrganization")
		return nil, errtrace.Wrap(err)
	}

	// ユーザー情報を取得
	userInfo, err := s.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "userRepository.FindByID")
		return nil, errtrace.Wrap(err)
	}

	// セッショントークンを生成
	tokenStr, err := s.tokenManager.Generate(
		ctx,
		*userInfo,
		newSession.SessionID(),
	)
	if err != nil {
		utils.HandleError(ctx, err, "sessionTokenManager.GenerateTokenWithOrganization")
		return nil, errtrace.Wrap(err)
	}

	return &SwitchOrganizationUseCaseOutput{
		SessionTokenStr: tokenStr,
	}, nil
}

func NewSwitchOrganizationUseCase(
	organizationService organization_svc.OrganizationService,
	organizationUserRepo organization.OrganizationUserRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
	userRepository user.UserRepository,
) SwitchOrganizationUseCase {
	return &switchOrganizationInteractor{
		organizationService:  organizationService,
		organizationUserRepo: organizationUserRepo,
		sessionService:       sessionService,
		tokenManager:         tokenManager,
		userRepository:       userRepository,
	}
}
