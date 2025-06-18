package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

// organizationAliasRepository 実装
type organizationAliasRepository struct {
	*db.DBManager
}

// NewOrganizationAliasRepository コンストラクタ
func NewOrganizationAliasRepository(dbManager *db.DBManager) organization.OrganizationAliasRepository {
	return &organizationAliasRepository{
		DBManager: dbManager,
	}
}

// Create エイリアスを作成
func (r *organizationAliasRepository) Create(ctx context.Context, alias *organization.OrganizationAlias) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationAliasRepository.Create")
	defer span.End()

	queries := r.GetQueries(ctx)
	params := model.CreateOrganizationAliasParams{
		OrganizationID: alias.OrganizationID().UUID(),
		AliasName:      alias.AliasName(),
		CreatedBy:      alias.CreatedBy().UUID(),
	}
	_, err := queries.CreateOrganizationAlias(ctx, params)
	return err
}

// FindByID IDでエイリアスを取得
func (r *organizationAliasRepository) FindByID(ctx context.Context, aliasID shared.UUID[organization.OrganizationAlias]) (*organization.OrganizationAlias, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationAliasRepository.FindByID")
	defer span.End()

	queries := r.GetQueries(ctx)
	row, err := queries.GetOrganizationAliasById(ctx, aliasID.UUID())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("organization alias not found")
		}
		return nil, err
	}
	return r.fromRow(&row), nil
}

// FindActiveByOrganizationID 組織のアクティブなエイリアスを取得
func (r *organizationAliasRepository) FindActiveByOrganizationID(ctx context.Context, organizationID shared.UUID[organization.Organization]) ([]*organization.OrganizationAlias, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationAliasRepository.FindActiveByOrganizationID")
	defer span.End()

	queries := r.GetQueries(ctx)
	rows, err := queries.GetActiveOrganizationAliases(ctx, organizationID.UUID())
	if err != nil {
		return nil, err
	}
	aliases := make([]*organization.OrganizationAlias, 0, len(rows))
	for _, row := range rows {
		aliases = append(aliases, r.fromRow(&row))
	}
	return aliases, nil
}

// Deactivate エイリアスを論理削除
func (r *organizationAliasRepository) Deactivate(ctx context.Context, aliasID shared.UUID[organization.OrganizationAlias], deactivatedBy shared.UUID[user.User]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationAliasRepository.Deactivate")
	defer span.End()

	queries := r.GetQueries(ctx)
	params := model.DeactivateOrganizationAliasParams{
		AliasID:       aliasID.UUID(),
		DeactivatedBy: uuid.NullUUID{UUID: deactivatedBy.UUID(), Valid: true},
	}
	return queries.DeactivateOrganizationAlias(ctx, params)
}

// CountActiveByOrganizationID 組織のアクティブなエイリアス数を取得
func (r *organizationAliasRepository) CountActiveByOrganizationID(ctx context.Context, organizationID shared.UUID[organization.Organization]) (int64, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationAliasRepository.CountActiveByOrganizationID")
	defer span.End()

	queries := r.GetQueries(ctx)
	return queries.CountActiveAliasesByOrganization(ctx, organizationID.UUID())
}

// ExistsActiveAliasName 同じ名前のアクティブなエイリアスが存在するか確認
func (r *organizationAliasRepository) ExistsActiveAliasName(ctx context.Context, organizationID shared.UUID[organization.Organization], aliasName string) (bool, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationAliasRepository.ExistsActiveAliasName")
	defer span.End()

	queries := r.GetQueries(ctx)
	params := model.CheckAliasNameExistsParams{
		OrganizationID: organizationID.UUID(),
		AliasName:      aliasName,
	}
	return queries.CheckAliasNameExists(ctx, params)
}

// fromRow SQLの行からドメインモデルに変換
func (r *organizationAliasRepository) fromRow(row *model.OrganizationAlias) *organization.OrganizationAlias {
	aliasID, _ := shared.ParseUUID[organization.OrganizationAlias](row.AliasID.String())
	organizationID, _ := shared.ParseUUID[organization.Organization](row.OrganizationID.String())
	createdBy, _ := shared.ParseUUID[user.User](row.CreatedBy.String())
	var deactivatedAt *time.Time
	if row.DeactivatedAt.Valid {
		deactivatedAt = &row.DeactivatedAt.Time
	}
	var deactivatedBy *shared.UUID[user.User]
	if row.DeactivatedBy.Valid {
		id, _ := shared.ParseUUID[user.User](row.DeactivatedBy.UUID.String())
		deactivatedBy = &id
	}
	return organization.ReconstructOrganizationAlias(
		aliasID,
		organizationID,
		row.AliasName,
		row.CreatedAt,
		row.UpdatedAt,
		createdBy,
		deactivatedAt,
		deactivatedBy,
	)
}
