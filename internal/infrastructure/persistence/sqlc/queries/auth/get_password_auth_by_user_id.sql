-- name: GetPasswordAuthByUserId :one
SELECT
  sqlc.embed(password_auth)
FROM password_auth
WHERE user_id = $1;
