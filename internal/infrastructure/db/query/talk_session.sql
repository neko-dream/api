-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, owner_id, created_at) VALUES ($1, $2, $3, $4);

-- name: GetTalkSessionByID :one
SELECT * FROM talk_sessions WHERE talk_session_id = $1;

-- name: EditTalkSession :exec
UPDATE talk_sessions SET theme = $2, finished_at = $3 WHERE talk_session_id = $1;

-- name: ListTalkSessions :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.finished_at,
    talk_sessions.created_at,
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url
FROM talk_sessions
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON talk_sessions.talk_session_id = oc.talk_session_id
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
WHERE
    CASE
        WHEN sqlc.narg('status')::text = 'finished' THEN finished_at IS NOT NULL
        WHEN sqlc.narg('status')::text = 'open' THEN finished_at IS NULL
        ELSE TRUE
    END
    AND
    (CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
        THEN talk_sessions.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END)
ORDER BY
    talk_sessions.created_at DESC
LIMIT $1 OFFSET $2;
