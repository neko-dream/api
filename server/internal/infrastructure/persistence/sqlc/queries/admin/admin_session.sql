-- name: CreateAdminSession :exec
INSERT INTO admin_sessions (admin_session_id, admin_user_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetAdminSessionByTokenHash :one
SELECT * FROM admin_sessions WHERE token_hash = $1;

-- name: RotateAdminSession :execrows
UPDATE admin_sessions
SET token_hash = $2, expires_at = $3, updated_at = CURRENT_TIMESTAMP
WHERE admin_session_id = $1 AND token_hash = $4 AND expires_at > CURRENT_TIMESTAMP;

-- name: DeleteAdminSessionByTokenHash :exec
DELETE FROM admin_sessions WHERE token_hash = $1;
