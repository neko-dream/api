-- name: CreateUser :exec
INSERT INTO users (user_id, display_id, display_name, created_at, picture) VALUES ($1, $2, $3, $4, $5);

-- name: CreateUserAuth :exec
INSERT INTO user_auths (user_id, provider, subject, created_at) VALUES ($1, $2, $3, $4);

-- name: GetUserBySubject :one
SELECT
    users.user_id,
    users.display_id,
    users.display_name,
    user_auths.provider,
    user_auths.subject,
    user_auths.created_at,
    users.picture
FROM
    users
    JOIN user_auths ON users.user_id = user_auths.user_id
WHERE
    user_auths.subject = $1;
