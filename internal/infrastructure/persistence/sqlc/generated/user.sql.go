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

// CreateUser
//
//	INSERT INTO users (user_id, created_at) VALUES ($1, $2)
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

// CreateUserAuth
//
//	INSERT INTO user_auths (user_auth_id, user_id, provider, subject, created_at, is_verified) VALUES ($1, $2, $3, $4, $5, false)
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

const getUserAuthByUserID = `-- name: GetUserAuthByUserID :one
SELECT
    user_auth_id, user_id, provider, subject, is_verified, created_at
FROM
    "user_auths"
WHERE
    user_id = $1
`

// GetUserAuthByUserID
//
//	SELECT
//	    user_auth_id, user_id, provider, subject, is_verified, created_at
//	FROM
//	    "user_auths"
//	WHERE
//	    user_id = $1
func (q *Queries) GetUserAuthByUserID(ctx context.Context, userID uuid.UUID) (UserAuth, error) {
	row := q.db.QueryRowContext(ctx, getUserAuthByUserID, userID)
	var i UserAuth
	err := row.Scan(
		&i.UserAuthID,
		&i.UserID,
		&i.Provider,
		&i.Subject,
		&i.IsVerified,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT
    user_id, display_id, display_name, icon_url, created_at, updated_at
FROM
    "users"
WHERE
    users.user_id = $1
`

// GetUserByID
//
//	SELECT
//	    user_id, display_id, display_name, icon_url, created_at, updated_at
//	FROM
//	    "users"
//	WHERE
//	    users.user_id = $1
func (q *Queries) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.DisplayID,
		&i.DisplayName,
		&i.IconUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserBySubject = `-- name: GetUserBySubject :one
SELECT
    "users".user_id,
    "users".display_id,
    "users".display_name,
    "user_auths".provider,
    "user_auths".subject,
    "user_auths".created_at,
    "users".icon_url,
    "user_auths".is_verified
FROM
    "users"
    JOIN "user_auths" ON "users".user_id = "user_auths".user_id
WHERE
    "user_auths".subject = $1
`

type GetUserBySubjectRow struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	Provider    string
	Subject     string
	CreatedAt   time.Time
	IconUrl     sql.NullString
	IsVerified  bool
}

// GetUserBySubject
//
//	SELECT
//	    "users".user_id,
//	    "users".display_id,
//	    "users".display_name,
//	    "user_auths".provider,
//	    "user_auths".subject,
//	    "user_auths".created_at,
//	    "users".icon_url,
//	    "user_auths".is_verified
//	FROM
//	    "users"
//	    JOIN "user_auths" ON "users".user_id = "user_auths".user_id
//	WHERE
//	    "user_auths".subject = $1
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
		&i.IconUrl,
		&i.IsVerified,
	)
	return i, err
}

const getUserDemographicByUserID = `-- name: GetUserDemographicByUserID :one
SELECT
    user_demographics_id, user_id, year_of_birth, occupation, gender, city, household_size, created_at, updated_at, prefecture
FROM
    "user_demographics"
WHERE
    user_id = $1
`

// GetUserDemographicByUserID
//
//	SELECT
//	    user_demographics_id, user_id, year_of_birth, occupation, gender, city, household_size, created_at, updated_at, prefecture
//	FROM
//	    "user_demographics"
//	WHERE
//	    user_id = $1
func (q *Queries) GetUserDemographicByUserID(ctx context.Context, userID uuid.UUID) (UserDemographic, error) {
	row := q.db.QueryRowContext(ctx, getUserDemographicByUserID, userID)
	var i UserDemographic
	err := row.Scan(
		&i.UserDemographicsID,
		&i.UserID,
		&i.YearOfBirth,
		&i.Occupation,
		&i.Gender,
		&i.City,
		&i.HouseholdSize,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Prefecture,
	)
	return i, err
}

const getUserDetailByID = `-- name: GetUserDetailByID :one
SELECT
    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
    user_auths.user_auth_id, user_auths.user_id, user_auths.provider, user_auths.subject, user_auths.is_verified, user_auths.created_at,
    user_demographics.user_demographics_id, user_demographics.user_id, user_demographics.year_of_birth, user_demographics.occupation, user_demographics.gender, user_demographics.city, user_demographics.household_size, user_demographics.created_at, user_demographics.updated_at, user_demographics.prefecture
FROM
    users
LEFT JOIN user_auths ON users.user_id = user_auths.user_id
LEFT JOIN user_demographics ON users.user_id = user_demographics.user_id
WHERE
    users.user_id = $1
`

type GetUserDetailByIDRow struct {
	User            User
	UserAuth        UserAuth
	UserDemographic UserDemographic
}

// GetUserDetailByID
//
//	SELECT
//	    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
//	    user_auths.user_auth_id, user_auths.user_id, user_auths.provider, user_auths.subject, user_auths.is_verified, user_auths.created_at,
//	    user_demographics.user_demographics_id, user_demographics.user_id, user_demographics.year_of_birth, user_demographics.occupation, user_demographics.gender, user_demographics.city, user_demographics.household_size, user_demographics.created_at, user_demographics.updated_at, user_demographics.prefecture
//	FROM
//	    users
//	LEFT JOIN user_auths ON users.user_id = user_auths.user_id
//	LEFT JOIN user_demographics ON users.user_id = user_demographics.user_id
//	WHERE
//	    users.user_id = $1
func (q *Queries) GetUserDetailByID(ctx context.Context, userID uuid.UUID) (GetUserDetailByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserDetailByID, userID)
	var i GetUserDetailByIDRow
	err := row.Scan(
		&i.User.UserID,
		&i.User.DisplayID,
		&i.User.DisplayName,
		&i.User.IconUrl,
		&i.User.CreatedAt,
		&i.User.UpdatedAt,
		&i.UserAuth.UserAuthID,
		&i.UserAuth.UserID,
		&i.UserAuth.Provider,
		&i.UserAuth.Subject,
		&i.UserAuth.IsVerified,
		&i.UserAuth.CreatedAt,
		&i.UserDemographic.UserDemographicsID,
		&i.UserDemographic.UserID,
		&i.UserDemographic.YearOfBirth,
		&i.UserDemographic.Occupation,
		&i.UserDemographic.Gender,
		&i.UserDemographic.City,
		&i.UserDemographic.HouseholdSize,
		&i.UserDemographic.CreatedAt,
		&i.UserDemographic.UpdatedAt,
		&i.UserDemographic.Prefecture,
	)
	return i, err
}

const updateOrCreateUserDemographic = `-- name: UpdateOrCreateUserDemographic :exec
INSERT INTO user_demographics (
    user_demographics_id,
    user_id,
    year_of_birth,
    occupation,
    gender,
    city,
    household_size,
    prefecture,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
ON CONFLICT (user_id)
DO UPDATE SET
    year_of_birth = $3,
    occupation = $4,
    gender = $5,
    city = $6,
    household_size = $7,
    prefecture = $8,
    updated_at = now()
`

type UpdateOrCreateUserDemographicParams struct {
	UserDemographicsID uuid.UUID
	UserID             uuid.UUID
	YearOfBirth        sql.NullString
	Occupation         sql.NullInt16
	Gender             sql.NullString
	City               sql.NullString
	HouseholdSize      sql.NullInt16
	Prefecture         sql.NullString
}

// UpdateOrCreateUserDemographic
//
//	INSERT INTO user_demographics (
//	    user_demographics_id,
//	    user_id,
//	    year_of_birth,
//	    occupation,
//	    gender,
//	    city,
//	    household_size,
//	    prefecture,
//	    created_at,
//	    updated_at
//	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
//	ON CONFLICT (user_id)
//	DO UPDATE SET
//	    year_of_birth = $3,
//	    occupation = $4,
//	    gender = $5,
//	    city = $6,
//	    household_size = $7,
//	    prefecture = $8,
//	    updated_at = now()
func (q *Queries) UpdateOrCreateUserDemographic(ctx context.Context, arg UpdateOrCreateUserDemographicParams) error {
	_, err := q.db.ExecContext(ctx, updateOrCreateUserDemographic,
		arg.UserDemographicsID,
		arg.UserID,
		arg.YearOfBirth,
		arg.Occupation,
		arg.Gender,
		arg.City,
		arg.HouseholdSize,
		arg.Prefecture,
	)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE "users" SET display_id = $2, display_name = $3, icon_url = $4 WHERE user_id = $1
`

type UpdateUserParams struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	IconUrl     sql.NullString
}

// UpdateUser
//
//	UPDATE "users" SET display_id = $2, display_name = $3, icon_url = $4 WHERE user_id = $1
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.UserID,
		arg.DisplayID,
		arg.DisplayName,
		arg.IconUrl,
	)
	return err
}

const userFindByDisplayID = `-- name: UserFindByDisplayID :one
SELECT
    user_id, display_id, display_name, icon_url, created_at, updated_at
FROM
    "users"
WHERE
    display_id = $1
`

// UserFindByDisplayID
//
//	SELECT
//	    user_id, display_id, display_name, icon_url, created_at, updated_at
//	FROM
//	    "users"
//	WHERE
//	    display_id = $1
func (q *Queries) UserFindByDisplayID(ctx context.Context, displayID sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, userFindByDisplayID, displayID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.DisplayID,
		&i.DisplayName,
		&i.IconUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const verifyUser = `-- name: VerifyUser :exec
UPDATE "user_auths" SET is_verified = true WHERE user_id = $1
`

// VerifyUser
//
//	UPDATE "user_auths" SET is_verified = true WHERE user_id = $1
func (q *Queries) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, verifyUser, userID)
	return err
}
