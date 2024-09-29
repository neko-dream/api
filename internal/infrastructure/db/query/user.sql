-- name: CreateUser :exec
INSERT INTO users (user_id, created_at) VALUES ($1, $2);

-- name: CreateUserAuth :exec
INSERT INTO user_auths (user_auth_id, user_id, provider, subject, created_at, is_verified) VALUES ($1, $2, $3, $4, $5, false);

-- name: GetUserBySubject :one
SELECT
    "users".user_id,
    "users".display_id,
    "users".display_name,
    "user_auths".provider,
    "user_auths".subject,
    "user_auths".created_at,
    "users".picture,
    "user_auths".is_verified
FROM
    "users"
    JOIN "user_auths" ON "users".user_id = "user_auths".user_id
WHERE
    "user_auths".subject = $1;

-- name: GetUserByID :one
SELECT
    *
FROM
    "users"
WHERE
    users.user_id = $1;

-- name: GetUserAuthByUserID :one
SELECT
    *
FROM
    "user_auths"
WHERE
    user_id = $1;

-- name: VerifyUser :exec
UPDATE "user_auths" SET is_verified = true WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE "users" SET display_id = $2, display_name = $3, picture = $4 WHERE user_id = $1;

-- name: UpdateOrCreateUserDemographics :exec
INSERT INTO user_demographics (
    user_id,
    year_of_birth,
    occupation,
    gender,
    municipality,
    household_size,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, now(), now())
ON CONFLICT (user_id)
DO UPDATE SET
    year_of_birth = $2,
    occupation = $3,
    gender = $4,
    municipality = $5,
    household_size = $6,
    updated_at = now();



