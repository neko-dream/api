// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (user_id, created_at) VALUES ($1, $2)
`

type CreateUserParams struct {
	UserID    uuid.UUID
	CreatedAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.UserID, arg.CreatedAt)
	return err
}

const createUserAuth = `-- name: CreateUserAuth :exec
INSERT INTO user_auths (user_auth_id, user_id, provider, subject, created_at, is_verified) VALUES ($1, $2, $3, $4, $5, false)
`

type CreateUserAuthParams struct {
	UserAuthID uuid.UUID
	UserID     uuid.UUID
	Provider   string
	Subject    string
	CreatedAt  time.Time
}

func (q *Queries) CreateUserAuth(ctx context.Context, arg CreateUserAuthParams) error {
	_, err := q.db.ExecContext(ctx, createUserAuth,
		arg.UserAuthID,
		arg.UserID,
		arg.Provider,
		arg.Subject,
		arg.CreatedAt,
	)
	return err
}

const getUserByID = `-- name: GetUserByID :one
SELECT
    users.user_id,
    users.display_id,
    users.display_name,
    user_auths.provider,
    user_auths.subject,
    user_auths.created_at,
    users.picture,
    user_auths.is_verified
FROM
    users
    JOIN user_auths ON users.user_id = user_auths.user_id
WHERE
    users.user_id = $1
`

type GetUserByIDRow struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	Provider    string
	Subject     string
	CreatedAt   time.Time
	Picture     sql.NullString
	IsVerified  bool
}

func (q *Queries) GetUserByID(ctx context.Context, userID uuid.UUID) (GetUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var i GetUserByIDRow
	err := row.Scan(
		&i.UserID,
		&i.DisplayID,
		&i.DisplayName,
		&i.Provider,
		&i.Subject,
		&i.CreatedAt,
		&i.Picture,
		&i.IsVerified,
	)
	return i, err
}

const getUserBySubject = `-- name: GetUserBySubject :one
SELECT
    users.user_id,
    users.display_id,
    users.display_name,
    user_auths.provider,
    user_auths.subject,
    user_auths.created_at,
    users.picture,
    user_auths.is_verified
FROM
    users
    JOIN user_auths ON users.user_id = user_auths.user_id
WHERE
    user_auths.subject = $1
`

type GetUserBySubjectRow struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	Provider    string
	Subject     string
	CreatedAt   time.Time
	Picture     sql.NullString
	IsVerified  bool
}

func (q *Queries) GetUserBySubject(ctx context.Context, subject string) (GetUserBySubjectRow, error) {
	row := q.db.QueryRowContext(ctx, getUserBySubject, subject)
	var i GetUserBySubjectRow
	err := row.Scan(
		&i.UserID,
		&i.DisplayID,
		&i.DisplayName,
		&i.Provider,
		&i.Subject,
		&i.CreatedAt,
		&i.Picture,
		&i.IsVerified,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET display_id = $2, display_name = $3, picture = $4 WHERE user_id = $1
`

type UpdateUserParams struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	Picture     sql.NullString
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.UserID,
		arg.DisplayID,
		arg.DisplayName,
		arg.Picture,
	)
	return err
}

const verifyUser = `-- name: VerifyUser :exec
UPDATE user_auths SET is_verified = true WHERE user_id = $1
`

func (q *Queries) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, verifyUser, userID)
	return err
}
