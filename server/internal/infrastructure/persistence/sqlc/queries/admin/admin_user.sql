-- name: CreateAdminUser :exec
INSERT INTO admin_users (admin_user_id, email, display_name)
VALUES ($1, $2, $3);

-- name: GetAdminUserByEmail :one
SELECT * FROM admin_users WHERE email = $1;

-- name: GetAdminUserByID :one
SELECT * FROM admin_users WHERE admin_user_id = $1;

-- name: ListAdminUsers :many
SELECT * FROM admin_users ORDER BY created_at ASC;
