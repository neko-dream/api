-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, created_at) VALUES ($1, $2, $3);

-- name: GetTalkSessionByID :one
SELECT * FROM talk_sessions WHERE talk_session_id = $1;

