package organization

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

// OrganizationAliasRepository リポジトリインターフェース
type OrganizationAliasRepository interface {
	Create(ctx context.Context, alias *OrganizationAlias) error
	FindByID(ctx context.Context, aliasID shared.UUID[OrganizationAlias]) (*OrganizationAlias, error)
	FindActiveByOrganizationID(ctx context.Context, organizationID shared.UUID[Organization]) ([]*OrganizationAlias, error)
	Deactivate(ctx context.Context, aliasID shared.UUID[OrganizationAlias], deactivatedBy shared.UUID[user.User]) error
	CountActiveByOrganizationID(ctx context.Context, organizationID shared.UUID[Organization]) (int64, error)
	ExistsActiveAliasName(ctx context.Context, organizationID shared.UUID[Organization], aliasName string) (bool, error)
}

// OrganizationAlias 組織のエイリアス（別名）
type OrganizationAlias struct {
	aliasID        shared.UUID[OrganizationAlias]
	organizationID shared.UUID[Organization]
	aliasName      string
	createdAt      time.Time
	updatedAt      time.Time
	createdBy      shared.UUID[user.User]
	deactivatedAt  *time.Time
	deactivatedBy  *shared.UUID[user.User]
}

// NewOrganizationAlias 新しい組織エイリアスを作成
func NewOrganizationAlias(
	name string,
	organizationID shared.UUID[Organization],
	createdBy shared.UUID[user.User],
) (*OrganizationAlias, error) {
	// UTF-8文字数バリデーション（3-100文字）
	runeCount := utf8.RuneCountInString(name)
	if runeCount < 3 || runeCount > 100 {
		return nil, errors.New("alias name must be between 3 and 100 characters")
	}

	now := time.Now()
	return &OrganizationAlias{
		aliasID:        shared.NewUUID[OrganizationAlias](),
		organizationID: organizationID,
		aliasName:      name,
		createdAt:      now,
		updatedAt:      now,
		createdBy:      createdBy,
	}, nil
}

// ReconstructOrganizationAlias DBから取得したデータでエイリアスを再構築
func ReconstructOrganizationAlias(
	aliasID shared.UUID[OrganizationAlias],
	organizationID shared.UUID[Organization],
	aliasName string,
	createdAt time.Time,
	updatedAt time.Time,
	createdBy shared.UUID[user.User],
	deactivatedAt *time.Time,
	deactivatedBy *shared.UUID[user.User],
) *OrganizationAlias {
	return &OrganizationAlias{
		aliasID:        aliasID,
		organizationID: organizationID,
		aliasName:      aliasName,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		createdBy:      createdBy,
		deactivatedAt:  deactivatedAt,
		deactivatedBy:  deactivatedBy,
	}
}

// Deactivate エイリアスを論理削除
func (a *OrganizationAlias) Deactivate(deactivatedBy shared.UUID[user.User]) error {
	if a.deactivatedAt != nil {
		return errors.New("alias is already deactivated")
	}
	now := time.Now()
	a.deactivatedAt = &now
	a.deactivatedBy = &deactivatedBy
	a.updatedAt = now
	return nil
}

// IsActive エイリアスがアクティブかどうか
func (a *OrganizationAlias) IsActive() bool {
	return a.deactivatedAt == nil
}

// AliasID エイリアスIDを取得
func (a *OrganizationAlias) AliasID() shared.UUID[OrganizationAlias] {
	return a.aliasID
}

// OrganizationID 組織IDを取得
func (a *OrganizationAlias) OrganizationID() shared.UUID[Organization] {
	return a.organizationID
}

// AliasName エイリアス名を取得
func (a *OrganizationAlias) AliasName() string {
	return a.aliasName
}

// CreatedAt 作成日時を取得
func (a *OrganizationAlias) CreatedAt() time.Time {
	return a.createdAt
}

// UpdatedAt 更新日時を取得
func (a *OrganizationAlias) UpdatedAt() time.Time {
	return a.updatedAt
}

// CreatedBy 作成者を取得
func (a *OrganizationAlias) CreatedBy() shared.UUID[user.User] {
	return a.createdBy
}

// DeactivatedAt 削除日時を取得
func (a *OrganizationAlias) DeactivatedAt() *time.Time {
	return a.deactivatedAt
}

// DeactivatedBy 削除者を取得
func (a *OrganizationAlias) DeactivatedBy() *shared.UUID[user.User] {
	return a.deactivatedBy
}
