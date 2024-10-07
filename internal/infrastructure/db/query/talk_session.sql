-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, owner_id, scheduled_end_time, created_at) VALUES ($1, $2, $3, $4, $5);

-- name: CreateTalkSessionLocation :exec
INSERT INTO talk_session_locations (talk_session_id, location, city, prefecture) VALUES ($1, ST_GeographyFromText($2), $3, $4);

-- name: UpdateTalkSessionLocation :exec
UPDATE talk_session_locations SET location = ST_GeographyFromText($2), city = $3, prefecture = $4 WHERE talk_session_id = $1;

-- name: GetTalkSessionByID :one
SELECT * FROM talk_sessions WHERE talk_session_id = $1;

-- name: EditTalkSession :exec
UPDATE talk_sessions
    SET theme = $2,
        finished_at = $3,
        scheduled_end_time = $4
    WHERE talk_session_id = $1;

-- name: ListTalkSessions :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.finished_at,
    talk_sessions.created_at,
    talk_sessions.scheduled_end_time,
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))) AS latitude,
    ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))) AS longitude,
    talk_session_locations.city AS city,
    talk_session_locations.prefecture AS prefecture
FROM talk_sessions
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON talk_sessions.talk_session_id = oc.talk_session_id
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
LEFT JOIN talk_session_locations
    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
WHERE
    CASE
        WHEN sqlc.narg('status')::text = 'finished' THEN finished_at IS NOT NULL
        WHEN sqlc.narg('status')::text = 'open' THEN finished_at IS NULL AND scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    (CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
        THEN talk_sessions.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END)
ORDER BY
    CASE
        WHEN sqlc.narg('status')::text = 'finished' THEN finished_at IS NOT NULL
        WHEN sqlc.narg('status')::text = 'open' THEN scheduled_end_time > now()
        ELSE TRUE
    END DESC
LIMIT $1 OFFSET $2;
