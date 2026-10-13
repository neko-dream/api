-- name: LockPasswordSetupEmail :exec
SELECT pg_advisory_xact_lock(hashtextextended(sqlc.arg(email_subject)::text, 0));

-- name: LockPasswordSetupUser :one
SELECT user_id FROM users WHERE user_id = $1 FOR UPDATE;

-- name: FindPasswordSetupUserByTokenHash :one
SELECT user_id FROM password_setup_tokens WHERE token_hash = $1;

-- name: GetPasswordSetupToken :one
SELECT * FROM password_setup_tokens WHERE user_id = $1 FOR UPDATE;

-- name: SavePasswordSetupToken :exec
INSERT INTO password_setup_tokens (user_id, token_hash, email_subject, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE SET token_hash = EXCLUDED.token_hash,
    email_subject = EXCLUDED.email_subject, created_at = EXCLUDED.created_at, expires_at = EXCLUDED.expires_at;

-- name: DeletePasswordSetupToken :exec
DELETE FROM password_setup_tokens WHERE user_id = $1;
