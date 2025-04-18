// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_organization_by_id.sql

package model

import (
	"context"

	"github.com/google/uuid"
)

const findOrganizationByID = `-- name: FindOrganizationByID :one
SELECT
    organizations.organization_id, organizations.organization_type, organizations.name, organizations.owner_id
FROM organizations
WHERE organization_id = $1
`

type FindOrganizationByIDRow struct {
	Organization Organization
}

// FindOrganizationByID
//
//	SELECT
//	    organizations.organization_id, organizations.organization_type, organizations.name, organizations.owner_id
//	FROM organizations
//	WHERE organization_id = $1
func (q *Queries) FindOrganizationByID(ctx context.Context, organizationID uuid.UUID) (FindOrganizationByIDRow, error) {
	row := q.db.QueryRowContext(ctx, findOrganizationByID, organizationID)
	var i FindOrganizationByIDRow
	err := row.Scan(
		&i.Organization.OrganizationID,
		&i.Organization.OrganizationType,
		&i.Organization.Name,
		&i.Organization.OwnerID,
	)
	return i, err
}
