// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: organization_alias.sql

package model

import (
	"context"

	"github.com/google/uuid"
)

const checkAliasNameExists = `-- name: CheckAliasNameExists :one
SELECT EXISTS(
    SELECT 1 FROM organization_aliases
    WHERE organization_id = $1 AND alias_name = $2 AND deactivated_at IS NULL
)
`

type CheckAliasNameExistsParams struct {
	OrganizationID uuid.UUID
	AliasName      string
}

// CheckAliasNameExists
//
//	SELECT EXISTS(
//	    SELECT 1 FROM organization_aliases
//	    WHERE organization_id = $1 AND alias_name = $2 AND deactivated_at IS NULL
//	)
func (q *Queries) CheckAliasNameExists(ctx context.Context, arg CheckAliasNameExistsParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkAliasNameExists, arg.OrganizationID, arg.AliasName)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const countActiveAliasesByOrganization = `-- name: CountActiveAliasesByOrganization :one
SELECT COUNT(*) FROM organization_aliases
WHERE organization_id = $1 AND deactivated_at IS NULL
`

// CountActiveAliasesByOrganization
//
//	SELECT COUNT(*) FROM organization_aliases
//	WHERE organization_id = $1 AND deactivated_at IS NULL
func (q *Queries) CountActiveAliasesByOrganization(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, countActiveAliasesByOrganization, organizationID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createOrganizationAlias = `-- name: CreateOrganizationAlias :one
INSERT INTO organization_aliases (
    alias_id,
    organization_id,
    alias_name,
    created_by
) VALUES ($1, $2, $3, $4)
RETURNING alias_id, organization_id, alias_name, created_at, updated_at, created_by, deactivated_at, deactivated_by
`

type CreateOrganizationAliasParams struct {
	AliasID        uuid.UUID
	OrganizationID uuid.UUID
	AliasName      string
	CreatedBy      uuid.UUID
}

// CreateOrganizationAlias
//
//	INSERT INTO organization_aliases (
//	    alias_id,
//	    organization_id,
//	    alias_name,
//	    created_by
//	) VALUES ($1, $2, $3, $4)
//	RETURNING alias_id, organization_id, alias_name, created_at, updated_at, created_by, deactivated_at, deactivated_by
func (q *Queries) CreateOrganizationAlias(ctx context.Context, arg CreateOrganizationAliasParams) (OrganizationAlias, error) {
	row := q.db.QueryRowContext(ctx, createOrganizationAlias,
		arg.AliasID,
		arg.OrganizationID,
		arg.AliasName,
		arg.CreatedBy,
	)
	var i OrganizationAlias
	err := row.Scan(
		&i.AliasID,
		&i.OrganizationID,
		&i.AliasName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.DeactivatedAt,
		&i.DeactivatedBy,
	)
	return i, err
}

const deactivateOrganizationAlias = `-- name: DeactivateOrganizationAlias :exec
UPDATE organization_aliases
SET deactivated_at = CURRENT_TIMESTAMP,
    deactivated_by = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE alias_id = $1 AND deactivated_at IS NULL
`

type DeactivateOrganizationAliasParams struct {
	AliasID       uuid.UUID
	DeactivatedBy uuid.NullUUID
}

// DeactivateOrganizationAlias
//
//	UPDATE organization_aliases
//	SET deactivated_at = CURRENT_TIMESTAMP,
//	    deactivated_by = $2,
//	    updated_at = CURRENT_TIMESTAMP
//	WHERE alias_id = $1 AND deactivated_at IS NULL
func (q *Queries) DeactivateOrganizationAlias(ctx context.Context, arg DeactivateOrganizationAliasParams) error {
	_, err := q.db.ExecContext(ctx, deactivateOrganizationAlias, arg.AliasID, arg.DeactivatedBy)
	return err
}

const getActiveOrganizationAliases = `-- name: GetActiveOrganizationAliases :many
SELECT alias_id, organization_id, alias_name, created_at, updated_at, created_by, deactivated_at, deactivated_by FROM organization_aliases
WHERE organization_id = $1 AND deactivated_at IS NULL
ORDER BY created_at ASC
`

// GetActiveOrganizationAliases
//
//	SELECT alias_id, organization_id, alias_name, created_at, updated_at, created_by, deactivated_at, deactivated_by FROM organization_aliases
//	WHERE organization_id = $1 AND deactivated_at IS NULL
//	ORDER BY created_at ASC
func (q *Queries) GetActiveOrganizationAliases(ctx context.Context, organizationID uuid.UUID) ([]OrganizationAlias, error) {
	rows, err := q.db.QueryContext(ctx, getActiveOrganizationAliases, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationAlias
	for rows.Next() {
		var i OrganizationAlias
		if err := rows.Scan(
			&i.AliasID,
			&i.OrganizationID,
			&i.AliasName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.DeactivatedAt,
			&i.DeactivatedBy,
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

const getOrganizationAliasById = `-- name: GetOrganizationAliasById :one
SELECT alias_id, organization_id, alias_name, created_at, updated_at, created_by, deactivated_at, deactivated_by FROM organization_aliases
WHERE alias_id = $1
`

// GetOrganizationAliasById
//
//	SELECT alias_id, organization_id, alias_name, created_at, updated_at, created_by, deactivated_at, deactivated_by FROM organization_aliases
//	WHERE alias_id = $1
func (q *Queries) GetOrganizationAliasById(ctx context.Context, aliasID uuid.UUID) (OrganizationAlias, error) {
	row := q.db.QueryRowContext(ctx, getOrganizationAliasById, aliasID)
	var i OrganizationAlias
	err := row.Scan(
		&i.AliasID,
		&i.OrganizationID,
		&i.AliasName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.DeactivatedAt,
		&i.DeactivatedBy,
	)
	return i, err
}
