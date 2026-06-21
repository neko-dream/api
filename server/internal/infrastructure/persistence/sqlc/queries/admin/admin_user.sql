-- name: CreateAdminUser :exec
INSERT INTO admin_users (admin_user_id, email, display_name, password_hash, salt, role)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetAdminUserByEmail :one
SELECT * FROM admin_users WHERE email = $1;

-- name: GetAdminUserByID :one
SELECT * FROM admin_users WHERE admin_user_id = $1;
