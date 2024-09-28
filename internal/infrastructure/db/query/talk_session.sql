-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, owner_id, created_at) VALUES ($1, $2, $3, $4);

-- name: GetTalkSessionByID :one
SELECT * FROM talk_sessions WHERE talk_session_id = $1;

-- name: EditTalkSession :exec
UPDATE talk_sessions SET theme = $2, finished_at = $3 WHERE talk_session_id = $1;
