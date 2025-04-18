package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type organizationUserRepository struct {
	*db.DBManager
}

func NewOrganizationUserRepository(dbManager *db.DBManager) organization.OrganizationUserRepository {
	return &organizationUserRepository{
		DBManager: dbManager,
	}
}

// Create implements organization.OrganizationUserRepository.
func (o *organizationUserRepository) Create(ctx context.Context, orgUser organization.OrganizationUser) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationUserRepository.Create")
	defer span.End()

	if err := o.GetQueries(ctx).CreateOrgUser(ctx, model.CreateOrgUserParams{
		OrganizationUserID: orgUser.OrganizationUserID.UUID(),
		OrganizationID:     orgUser.OrganizationID.UUID(),
		UserID:             orgUser.UserID.UUID(),
		Role:               int32(orgUser.Role),
		UpdatedAt:          clock.Now(ctx),
		CreatedAt:          clock.Now(ctx),
	}); err != nil {
		return err
	}

	return nil
}

// FindByOrganizationID implements organization.OrganizationUserRepository.
func (o *organizationUserRepository) FindByOrganizationID(ctx context.Context, orgID shared.UUID[organization.Organization]) ([]*organization.OrganizationUser, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationUserRepository.FindByOrganizationID")
	defer span.End()

	orgUsers, err := o.GetQueries(ctx).FindOrgUserByOrganizationID(ctx, orgID.UUID())
	if err != nil {
		return nil, err
	}

	result := make([]*organization.OrganizationUser, len(orgUsers))
	for i, orgUser := range orgUsers {
		result[i] = &organization.OrganizationUser{
			OrganizationUserID: shared.UUID[organization.OrganizationUser](orgUser.OrganizationUser.OrganizationUserID),
			OrganizationID:     shared.UUID[organization.Organization](orgUser.OrganizationUser.OrganizationID),
			UserID:             shared.UUID[user.User](orgUser.OrganizationUser.UserID),
			Role:               organization.OrganizationUserRole(orgUser.OrganizationUser.Role),
		}
	}
	return result, nil
}

// FindByOrganizationIDAndUserID implements organization.OrganizationUserRepository.
func (o *organizationUserRepository) FindByOrganizationIDAndUserID(ctx context.Context, orgID shared.UUID[organization.Organization], userID shared.UUID[user.User]) (*organization.OrganizationUser, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationUserRepository.FindByOrganizationIDAndUserID")
	defer span.End()

	orgUser, err := o.GetQueries(ctx).FindOrgUserByOrganizationIDAndUserID(ctx, model.FindOrgUserByOrganizationIDAndUserIDParams{
		OrganizationID: orgID.UUID(),
		UserID:         userID.UUID(),
	})
	if err != nil {
		return nil, err
	}

	return &organization.OrganizationUser{
		OrganizationUserID: shared.UUID[organization.OrganizationUser](orgUser.OrganizationUser.OrganizationUserID),
		OrganizationID:     shared.UUID[organization.Organization](orgUser.OrganizationUser.OrganizationID),
		UserID:             shared.UUID[user.User](orgUser.OrganizationUser.UserID),
		Role:               organization.OrganizationUserRole(orgUser.OrganizationUser.Role),
	}, nil
}

// FindByUserID implements organization.OrganizationUserRepository.
func (o *organizationUserRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]*organization.OrganizationUser, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "organizationUserRepository.FindByUserID")
	defer span.End()

	orgUsers, err := o.GetQueries(ctx).FindOrgUserByUserID(ctx, userID.UUID())
	if err != nil {
		return nil, err
	}

	result := make([]*organization.OrganizationUser, len(orgUsers))
	for i, orgUser := range orgUsers {
		result[i] = &organization.OrganizationUser{
			OrganizationUserID: shared.UUID[organization.OrganizationUser](orgUser.OrganizationUser.OrganizationUserID),
			OrganizationID:     shared.UUID[organization.Organization](orgUser.OrganizationUser.OrganizationID),
			UserID:             shared.UUID[user.User](orgUser.OrganizationUser.UserID),
			Role:               organization.OrganizationUserRole(orgUser.OrganizationUser.Role),
		}
	}
	return result, nil
}
