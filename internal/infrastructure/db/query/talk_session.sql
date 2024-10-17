-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, owner_id, scheduled_end_time, created_at, city, prefecture) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateTalkSessionLocation :exec
INSERT INTO talk_session_locations (talk_session_id, location) VALUES ($1, ST_GeographyFromText($2));

-- name: UpdateTalkSessionLocation :exec
UPDATE talk_session_locations SET location = ST_GeographyFromText($2) WHERE talk_session_id = $1;

-- name: EditTalkSession :exec
UPDATE talk_sessions
    SET theme = $2,
        scheduled_end_time = $3
    WHERE talk_session_id = $1;

-- name: GetTalkSessionByID :one
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.created_at,
    talk_sessions.scheduled_end_time,
    talk_sessions.city AS city,
    talk_sessions.prefecture AS prefecture,
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    CASE
        WHEN talk_session_locations.location IS NULL THEN NULL
        ELSE ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location)))
    END AS latitude,
    CASE
        WHEN talk_session_locations.location IS NULL THEN NULL
        ELSE ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location)))
    END AS longitude
FROM talk_sessions
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
LEFT JOIN (
    SELECT opinions.talk_session_id, COUNT(opinions.opinion_id) AS opinion_count
    FROM opinions
    GROUP BY opinions.talk_session_id
) oc ON talk_sessions.talk_session_id = oc.talk_session_id
LEFT JOIN talk_session_locations
    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
WHERE talk_sessions.talk_session_id = $1;

-- name: ListTalkSessions :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.scheduled_end_time,
    talk_sessions.city AS city,
    talk_sessions.prefecture AS prefecture,
    talk_sessions.created_at,
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    CASE
        WHEN talk_session_locations.location IS NULL THEN NULL
        ELSE ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location)))
    END AS latitude,
    CASE
        WHEN talk_session_locations.location IS NULL THEN NULL
        ELSE ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location)))
    END AS longitude
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
        WHEN sqlc.narg('status')::text = 'finished' THEN scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN scheduled_end_time > now()
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
        WHEN sqlc.narg('status')::text = 'finished' THEN scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN scheduled_end_time > now()
        ELSE TRUE
    END DESC
LIMIT $1 OFFSET $2;

-- name: CountTalkSessions :one
SELECT
    COUNT(talk_sessions.*) AS talk_session_count
FROM talk_sessions
LEFT JOIN talk_session_locations
    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
WHERE
    CASE
        WHEN sqlc.narg('status')::text = 'finished' THEN scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    (CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
        THEN talk_sessions.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END);



-- name: GetTalkSessionByUserID :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.scheduled_end_time,
    talk_sessions.city AS city,
    talk_sessions.prefecture AS prefecture,
    talk_sessions.created_at,
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    CASE
        WHEN talk_session_locations.location IS NULL THEN NULL
        ELSE ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location)))
    END AS latitude,
    CASE
        WHEN talk_session_locations.location IS NULL THEN NULL
        ELSE ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location)))
    END AS longitude
FROM talk_sessions
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON talk_sessions.talk_session_id = oc.talk_session_id
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
LEFT JOIN votes
    ON talk_sessions.talk_session_id = votes.talk_session_id
LEFT JOIN talk_session_locations
    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
WHERE
    votes.user_id = sqlc.narg('user_id')::uuid AND
    CASE
        WHEN sqlc.narg('status')::text = 'finished' THEN scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
            THEN talk_sessions.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END
ORDER BY
    CASE
        WHEN sqlc.narg('status')::text = 'finished' THEN scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN scheduled_end_time > now()
        ELSE TRUE
    END DESC
LIMIT $1 OFFSET $2;
