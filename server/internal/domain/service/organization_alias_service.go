package service

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

// OrganizationAliasService 組織エイリアス管理サービス
type OrganizationAliasService struct {
	orgRepo     organization.OrganizationRepository
	orgUserRepo organization.OrganizationUserRepository
	aliasRepo   organization.OrganizationAliasRepository
}

// NewOrganizationAliasService コンストラクタ
func NewOrganizationAliasService(
	orgRepo organization.OrganizationRepository,
	orgUserRepo organization.OrganizationUserRepository,
	aliasRepo organization.OrganizationAliasRepository,
) *OrganizationAliasService {
	return &OrganizationAliasService{
		orgRepo:     orgRepo,
		orgUserRepo: orgUserRepo,
		aliasRepo:   aliasRepo,
	}
}

// CreateAlias エイリアスを作成
func (s *OrganizationAliasService) CreateAlias(
	ctx context.Context,
	name string,
	organizationID shared.UUID[organization.Organization],
	createdBy shared.UUID[user.User],
) (*organization.OrganizationAlias, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "OrganizationAliasService.CreateAlias")
	defer span.End()

	// 権限チェック
	canManage, err := s.CanManageAlias(ctx, createdBy, organizationID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, errors.New("permission denied: cannot manage organization aliases")
	}

	// 既存のエイリアス名チェック
	exists, err := s.aliasRepo.ExistsActiveAliasName(ctx, organizationID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("alias name already exists")
	}

	// エイリアス作成
	alias, err := organization.NewOrganizationAlias(name, organizationID, createdBy)
	if err != nil {
		return nil, err
	}

	if err := s.aliasRepo.Create(ctx, alias); err != nil {
		return nil, err
	}

	return alias, nil
}

// DeactivateAlias エイリアスを論理削除
func (s *OrganizationAliasService) DeactivateAlias(
	ctx context.Context,
	aliasID shared.UUID[organization.OrganizationAlias],
	deactivatedBy shared.UUID[user.User],
) error {
	ctx, span := otel.Tracer("service").Start(ctx, "OrganizationAliasService.DeactivateAlias")
	defer span.End()

	// エイリアスを取得
	alias, err := s.aliasRepo.FindByID(ctx, aliasID)
	if err != nil {
		return err
	}

	// 権限チェック
	canManage, err := s.CanManageAlias(ctx, deactivatedBy, alias.OrganizationID())
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("permission denied: cannot manage organization aliases")
	}
	// 論理削除実行
	return s.aliasRepo.Deactivate(ctx, aliasID, deactivatedBy)
}

// GetActiveAliases 組織のアクティブなエイリアスを取得
func (s *OrganizationAliasService) GetActiveAliases(
	ctx context.Context,
	organizationID shared.UUID[organization.Organization],
) ([]*organization.OrganizationAlias, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "OrganizationAliasService.GetActiveAliases")
	defer span.End()

	return s.aliasRepo.FindActiveByOrganizationID(ctx, organizationID)
}

// CanManageAlias ユーザーがエイリアスを管理できるか確認
func (s *OrganizationAliasService) CanManageAlias(
	ctx context.Context,
	userID shared.UUID[user.User],
	organizationID shared.UUID[organization.Organization],
) (bool, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "OrganizationAliasService.CanManageAlias")
	defer span.End()

	// ユーザーの組織内ロールを取得
	orgUser, err := s.orgUserRepo.FindByOrganizationIDAndUserID(ctx, organizationID, userID)
	if err != nil {
		return false, err
	}
	if orgUser == nil {
		return false, nil
	}

	// AdminまたはOwner権限が必要
	return orgUser.Role <= organization.OrganizationUserRoleAdmin, nil
}
