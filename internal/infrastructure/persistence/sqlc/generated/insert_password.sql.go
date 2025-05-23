// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: insert_password.sql

package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPasswordAuth = `-- name: CreatePasswordAuth :exec
INSERT INTO password_auth (
  password_auth_id,
  user_id,
  password_hash,
  salt,
  required_password_change,
  last_changed,
  created_at,
  updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreatePasswordAuthParams struct {
	PasswordAuthID         uuid.UUID
	UserID                 uuid.UUID
	PasswordHash           string
	Salt                   sql.NullString
	RequiredPasswordChange bool
	LastChanged            time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// CreatePasswordAuth
//
//	INSERT INTO password_auth (
//	  password_auth_id,
//	  user_id,
//	  password_hash,
//	  salt,
//	  required_password_change,
//	  last_changed,
//	  created_at,
//	  updated_at
//	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
func (q *Queries) CreatePasswordAuth(ctx context.Context, arg CreatePasswordAuthParams) error {
	_, err := q.db.ExecContext(ctx, createPasswordAuth,
		arg.PasswordAuthID,
		arg.UserID,
		arg.PasswordHash,
		arg.Salt,
		arg.RequiredPasswordChange,
		arg.LastChanged,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
