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
    "users".icon_url,
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

-- name: GetUserDetailByID :one
SELECT
    sqlc.embed(users),
    sqlc.embed(user_auths),
    sqlc.embed(user_demographics)
FROM
    users
LEFT JOIN user_auths ON users.user_id = user_auths.user_id
LEFT JOIN user_demographics ON users.user_id = user_demographics.user_id
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
UPDATE "users" SET display_id = $2, display_name = $3, icon_url = $4 WHERE user_id = $1;

-- name: UpdateOrCreateUserDemographic :exec
INSERT INTO user_demographics (
    user_demographics_id,
    user_id,
    year_of_birth,
    gender,
    city,
    prefecture,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6,  now(), now())
ON CONFLICT (user_id)
DO UPDATE SET
    year_of_birth = $3,
    gender = $4,
    city = $5,
    prefecture = $6,
    updated_at = now();

-- name: GetUserDemographicByUserID :one
SELECT
    *
FROM
    "user_demographics"
WHERE
    user_id = $1;

-- name: UserFindByDisplayID :one
SELECT
    *
FROM
    "users"
WHERE
    display_id = $1;
