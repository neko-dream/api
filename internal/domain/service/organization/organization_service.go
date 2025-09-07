package organization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
)

type OrganizationService interface {
	// 組織の作成・更新・削除
	CreateOrganization(ctx context.Context, name string, code string, iconURL *string, orgType organization.OrganizationType, ownerID shared.UUID[user.User]) (*organization.Organization, error)

	// ユーザーの所属組織
	GetUserOrganizations(ctx context.Context, userID shared.UUID[user.User]) ([]*organization.Organization, error)

	// 組織コードから組織IDを解決
	ResolveOrganizationIDFromCode(ctx context.Context, code *string) (*shared.UUID[any], error)
}

type organizationService struct {
	organizationRepo          organization.OrganizationRepository
	organizationUserRepo      organization.OrganizationUserRepository
	organizationMemberManager OrganizationMemberManager
	config                    *config.Config
}

func NewOrganizationService(
	organizationRepo organization.OrganizationRepository,
	organizationUserRepo organization.OrganizationUserRepository,
	organizationMemberManager OrganizationMemberManager,
	cfg *config.Config,
) OrganizationService {
	return &organizationService{
		organizationRepo:          organizationRepo,
		organizationUserRepo:      organizationUserRepo,
		organizationMemberManager: organizationMemberManager,
		config:                    cfg,
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

// CreateOrganizationWithCode implements OrganizationService.
func (s *organizationService) CreateOrganization(ctx context.Context, name string, code string, iconURL *string, orgType organization.OrganizationType, ownerID shared.UUID[user.User]) (*organization.Organization, error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "organizationService.CreateOrganizationWithCode")
	defer span.End()

	// 組織種別のバリデーション（0は無効）
	if orgType < organization.OrganizationTypeNormal || orgType > organization.OrganizationTypeCouncillor {
		return nil, messages.OrganizationTypeInvalid
	}

	// 開発環境以外はSuperAdminのみ作成可能
	if s.config.Env != config.LOCAL && s.config.Env != config.DEV {
		isSuperAdmin, err := s.organizationMemberManager.IsSuperAdmin(ctx, ownerID)
		if err != nil {
			return nil, err
		}
		if !isSuperAdmin {
			return nil, messages.OrganizationForbidden
		}
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

	// Validate provided code
	if err := organization.ValidateOrganizationCode(code); err != nil {
		return nil, err
	}
	// Check if code already exists
	_, err = s.organizationRepo.FindByCode(ctx, code)
	if err == nil {
		// Code exists
		return nil, messages.OrganizationCodeAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	org := organization.NewOrganization(
		orgID,
		orgType,
		name,
		code,
		iconURL,
		ownerID,
	)
	if err := s.organizationRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	// オーナーをメンバーとして追加
	orgUser := organization.OrganizationUser{
		OrganizationUserID: shared.NewUUID[organization.OrganizationUser](),
		OrganizationID:     orgID,
		UserID:             ownerID,
		Role:               organization.OrganizationUserRoleOwner,
	}
	if err := s.organizationUserRepo.Create(ctx, orgUser); err != nil {
		return nil, err
	}

	return org, nil
}

// 組織が見つからない場合はnilを返し、エラーとはしない
// 組織コードが無効な場合もnilを返す
func (s *organizationService) ResolveOrganizationIDFromCode(ctx context.Context, code *string) (*shared.UUID[any], error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "organizationService.ResolveOrganizationIDFromCode")
	defer span.End()

	if code == nil || *code == "" {
		return nil, nil
	}

	// 組織コードのバリデーション
	if err := organization.ValidateOrganizationCode(*code); err != nil {
		return nil, nil
	}

	org, err := s.organizationRepo.FindByCode(ctx, *code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 組織が見つからない場合は組織なしで続行
			return nil, nil
		}
		// その他のエラーは返す
		return nil, err
	}

	orgID := shared.UUID[any](org.OrganizationID)
	return &orgID, nil
}
