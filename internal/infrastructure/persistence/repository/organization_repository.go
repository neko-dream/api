package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type organizationRepository struct {
	*db.DBManager
}

func NewOrganizationRepository(dbManager *db.DBManager) organization.OrganizationRepository {
	return &organizationRepository{
		DBManager: dbManager,
	}
}

// Create
func (o *organizationRepository) Create(ctx context.Context, org *organization.Organization) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationRepository.Create")
	defer span.End()

	if err := o.GetQueries(ctx).CreateOrganization(ctx, model.CreateOrganizationParams{
		OrganizationID:   org.OrganizationID.UUID(),
		OrganizationType: int32(org.OrganizationType),
		Name:             org.Name,
		OwnerID:          org.OwnerID.UUID(),
		Code:             org.Code,
	}); err != nil {
		return err
	}
	return nil
}

// FindByID implements organization.OrganizationRepository.
func (o *organizationRepository) FindByID(ctx context.Context, id shared.UUID[organization.Organization]) (*organization.Organization, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationRepository.FindByID")
	defer span.End()

	org, err := o.GetQueries(ctx).FindOrganizationByID(ctx, id.UUID())
	if err != nil {
		return nil, err
	}

	return organization.NewOrganization(
		shared.UUID[organization.Organization](org.Organization.OrganizationID),
		organization.OrganizationType(org.Organization.OrganizationType),
		org.Organization.Name,
		org.Organization.Code,
		shared.UUID[user.User](org.Organization.OwnerID),
	), nil
}

// FindByIDs implements organization.OrganizationRepository.
func (o *organizationRepository) FindByIDs(ctx context.Context, ids []shared.UUID[organization.Organization]) ([]*organization.Organization, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationRepository.FindByIDs")
	defer span.End()

	if len(ids) == 0 {
		return nil, nil
	}
	// Convert shared.UUID[organization.Organization] to uuid.UUID
	uuidIDs := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		uuidIDs[i] = id.UUID()
	}

	orgs, err := o.GetQueries(ctx).FindOrganizationsByIDs(ctx, uuidIDs)
	if err != nil {
		return nil, err
	}

	var result []*organization.Organization
	for _, org := range orgs {
		result = append(result, organization.NewOrganization(
			shared.UUID[organization.Organization](org.Organization.OrganizationID),
			organization.OrganizationType(org.Organization.OrganizationType),
			org.Organization.Name,
			org.Organization.Code,
			shared.UUID[user.User](org.Organization.OwnerID),
		))
	}

	return result, nil
}

// FindByName implements organization.OrganizationRepository.
func (o *organizationRepository) FindByName(ctx context.Context, name string) (*organization.Organization, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationRepository.FindByName")
	defer span.End()

	org, err := o.GetQueries(ctx).FindOrganizationByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return organization.NewOrganization(
		shared.UUID[organization.Organization](org.Organization.OrganizationID),
		organization.OrganizationType(org.Organization.OrganizationType),
		org.Organization.Name,
		org.Organization.Code,
		shared.UUID[user.User](org.Organization.OwnerID),
	), nil
}

// FindByCode implements organization.OrganizationRepository.
func (o *organizationRepository) FindByCode(ctx context.Context, code string) (*organization.Organization, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationRepository.FindByCode")
	defer span.End()

	org, err := o.GetQueries(ctx).FindOrganizationByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return organization.NewOrganization(
		shared.UUID[organization.Organization](org.Organization.OrganizationID),
		organization.OrganizationType(org.Organization.OrganizationType),
		org.Organization.Name,
		org.Organization.Code,
		shared.UUID[user.User](org.Organization.OwnerID),
	), nil
}
