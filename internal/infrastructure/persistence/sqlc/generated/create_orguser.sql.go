// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: create_orguser.sql

package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createOrgUser = `-- name: CreateOrgUser :exec
INSERT INTO organization_users (
    organization_user_id,
    user_id,
    organization_id,
    role,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateOrgUserParams struct {
	OrganizationUserID uuid.UUID
	UserID             uuid.UUID
	OrganizationID     uuid.UUID
	Role               int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// CreateOrgUser
//
//	INSERT INTO organization_users (
//	    organization_user_id,
//	    user_id,
//	    organization_id,
//	    role,
//	    created_at,
//	    updated_at
//	) VALUES ($1, $2, $3, $4, $5, $6)
func (q *Queries) CreateOrgUser(ctx context.Context, arg CreateOrgUserParams) error {
	_, err := q.db.ExecContext(ctx, createOrgUser,
		arg.OrganizationUserID,
		arg.UserID,
		arg.OrganizationID,
		arg.Role,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
