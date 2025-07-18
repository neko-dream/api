// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_orguser_by_user_id.sql

package model

import (
	"context"

	"github.com/google/uuid"
)

const findOrgUserByUserID = `-- name: FindOrgUserByUserID :many
SELECT
    organization_users.organization_user_id, organization_users.user_id, organization_users.organization_id, organization_users.created_at, organization_users.updated_at, organization_users.role
FROM organization_users
WHERE organization_users.user_id = $1
`

type FindOrgUserByUserIDRow struct {
	OrganizationUser OrganizationUser
}

// FindOrgUserByUserID
//
//	SELECT
//	    organization_users.organization_user_id, organization_users.user_id, organization_users.organization_id, organization_users.created_at, organization_users.updated_at, organization_users.role
//	FROM organization_users
//	WHERE organization_users.user_id = $1
func (q *Queries) FindOrgUserByUserID(ctx context.Context, userID uuid.UUID) ([]FindOrgUserByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findOrgUserByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindOrgUserByUserIDRow
	for rows.Next() {
		var i FindOrgUserByUserIDRow
		if err := rows.Scan(
			&i.OrganizationUser.OrganizationUserID,
			&i.OrganizationUser.UserID,
			&i.OrganizationUser.OrganizationID,
			&i.OrganizationUser.CreatedAt,
			&i.OrganizationUser.UpdatedAt,
			&i.OrganizationUser.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findOrgUserByUserIDWithOrganization = `-- name: FindOrgUserByUserIDWithOrganization :many
SELECT
    organization_users.organization_user_id, organization_users.user_id, organization_users.organization_id, organization_users.created_at, organization_users.updated_at, organization_users.role,
    organizations.organization_id, organizations.organization_type, organizations.name, organizations.owner_id, organizations.code
FROM organization_users
LEFT JOIN organizations ON organization_users.organization_id = organizations.organization_id
WHERE organization_users.user_id = $1
`

type FindOrgUserByUserIDWithOrganizationRow struct {
	OrganizationUser OrganizationUser
	Organization     Organization
}

// FindOrgUserByUserIDWithOrganization
//
//	SELECT
//	    organization_users.organization_user_id, organization_users.user_id, organization_users.organization_id, organization_users.created_at, organization_users.updated_at, organization_users.role,
//	    organizations.organization_id, organizations.organization_type, organizations.name, organizations.owner_id, organizations.code
//	FROM organization_users
//	LEFT JOIN organizations ON organization_users.organization_id = organizations.organization_id
//	WHERE organization_users.user_id = $1
func (q *Queries) FindOrgUserByUserIDWithOrganization(ctx context.Context, userID uuid.UUID) ([]FindOrgUserByUserIDWithOrganizationRow, error) {
	rows, err := q.db.QueryContext(ctx, findOrgUserByUserIDWithOrganization, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindOrgUserByUserIDWithOrganizationRow
	for rows.Next() {
		var i FindOrgUserByUserIDWithOrganizationRow
		if err := rows.Scan(
			&i.OrganizationUser.OrganizationUserID,
			&i.OrganizationUser.UserID,
			&i.OrganizationUser.OrganizationID,
			&i.OrganizationUser.CreatedAt,
			&i.OrganizationUser.UpdatedAt,
			&i.OrganizationUser.Role,
			&i.Organization.OrganizationID,
			&i.Organization.OrganizationType,
			&i.Organization.Name,
			&i.Organization.OwnerID,
			&i.Organization.Code,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
