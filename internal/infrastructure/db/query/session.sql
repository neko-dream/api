-- name: FindActiveSessionsByUserID :many
SELECT *
FROM sessions
    WHERE user_id = $1
    AND session_status = 0;

-- name: CreateSession :exec
INSERT INTO sessions (session_id, user_id, provider, session_status, created_at, expires_at, last_activity_at) VALUES ($1, $2, $3, $4, $5, $6, $7);
