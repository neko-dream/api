-- name: CreateUser :exec
INSERT INTO users (user_id, created_at) VALUES ($1, $2);

-- name: CreateUserAuth :exec
INSERT INTO user_auths (user_auth_id, user_id, provider, subject, created_at, is_verified) VALUES ($1, $2, $3, $4, $5, false);

-- name: GetUserBySubject :one
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
    user_auths.subject = $1;

-- name: GetUserByID :one
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
    users.user_id = $1;


-- name: VerifyUser :exec
UPDATE user_auths SET is_verified = true WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE users SET display_id = $2, display_name = $3, picture = $4 WHERE user_id = $1;
