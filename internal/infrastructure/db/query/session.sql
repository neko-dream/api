-- name: FindActiveSessionsByUserID :many
SELECT *
FROM sessions
    WHERE user_id = $1
    AND session_status = 0;

-- name: CreateSession :exec
INSERT INTO sessions (session_id, user_id, provider, session_status, created_at, expires_at, last_activity_at) VALUES ($1, $2, $3, $4, $5, $6, $7);


-- name: UpdateSession :exec
UPDATE sessions
SET session_status = $2, last_activity_at = $3
WHERE session_id = $1;

-- name: FindSessionBySessionID :one
SELECT *
FROM sessions
WHERE session_id = $1;
