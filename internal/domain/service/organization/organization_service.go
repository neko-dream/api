package organization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type OrganizationService interface {
	// 組織の作成・更新・削除
	CreateOrganization(ctx context.Context, name string, orgType organization.OrganizationType, ownerID shared.UUID[user.User]) (*organization.Organization, error)

	// ユーザーの所属組織
	GetUserOrganizations(ctx context.Context, userID shared.UUID[user.User]) ([]*organization.Organization, error)
}

type organizationService struct {
	organizationRepo          organization.OrganizationRepository
	organizationUserRepo      organization.OrganizationUserRepository
	organizationMemberManager OrganizationMemberManager
}

func NewOrganizationService(
	organizationRepo organization.OrganizationRepository,
	organizationUserRepo organization.OrganizationUserRepository,
	organizationMemberManager OrganizationMemberManager,
) OrganizationService {
	return &organizationService{
		organizationRepo:          organizationRepo,
		organizationUserRepo:      organizationUserRepo,
		organizationMemberManager: organizationMemberManager,
	}
}

// GetUserOrganizations ユーザーの所属組織を取得
func (s *organizationService) GetUserOrganizations(ctx context.Context, userID shared.UUID[user.User]) ([]*organization.Organization, error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "organizationService.GetUserOrganizations")
	defer span.End()

	// ユーザーの組織ユーザーを取得
	orgUsers, err := s.organizationUserRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 組織IDのリストを作成
	orgIDs := make([]shared.UUID[organization.Organization], len(orgUsers))
	for i, orgUser := range orgUsers {
		orgIDs[i] = orgUser.OrganizationID
	}

	// 組織を取得
	orgs, err := s.organizationRepo.FindByIDs(ctx, orgIDs)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

// CreateOrganization implements OrganizationService.
func (s *organizationService) CreateOrganization(ctx context.Context, name string, orgType organization.OrganizationType, ownerID shared.UUID[user.User]) (*organization.Organization, error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "organizationService.CreateOrganization")
	defer span.End()

	// 現状はSuperAdminのみ作成可能
	isSuperAdmin, err := s.organizationMemberManager.IsSuperAdmin(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	if !isSuperAdmin {
		return nil, messages.OrganizationForbidden
	}

	// 名前が重複していないか確認
	existingOrg, err := s.organizationRepo.FindByName(ctx, name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	if existingOrg != nil {
		return nil, messages.OrganizationAlreadyExists
	}

	// 組織の作成
	orgID := shared.NewUUID[organization.Organization]()
	org := organization.NewOrganization(
		orgID,
		orgType,
		name,
		ownerID,
	)
	if err := s.organizationRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	// オーナーをメンバーとして追加
	orgUser := organization.OrganizationUser{
		OrganizationUserID: shared.NewUUID[organization.OrganizationUser](),
		OrganizationID: orgID,
		UserID:         ownerID,
		Role:           organization.OrganizationUserRoleOwner,
	}
	if err := s.organizationUserRepo.Create(ctx, orgUser); err != nil {
		return nil, err
	}

	return org, nil
}
