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
    talk_session_locations.talk_session_id as location_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
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
    talk_session_locations.talk_session_id as location_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
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
    CASE sqlc.narg('sort_key')::text
        WHEN 'latest' THEN EXTRACT(EPOCH FROM talk_sessions.created_at)
        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM talk_sessions.created_at)
        WHEN 'mostOpinions' THEN oc.opinion_count
        ELSE EXTRACT(EPOCH FROM talk_sessions.created_at)
    END ASC
LIMIT $1 OFFSET $2;

-- name: CountTalkSessions :one
SELECT
    COUNT(DISTINCT talk_sessions.talk_session_id) AS talk_session_count,
    sqlc.narg('status')::text AS status
FROM talk_sessions
-- talk_session_locationsがない場合も考慮
LEFT JOIN talk_session_locations
    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
LEFT JOIN votes
    ON votes.talk_session_id = talk_sessions.talk_session_id
WHERE
    CASE
        WHEN sqlc.narg('user_id')::uuid IS NOT NULL
            THEN votes.user_id = sqlc.narg('user_id')::uuid
        ELSE TRUE
    END
    AND
    CASE sqlc.narg('status')::text
        WHEN 'open' THEN talk_sessions.scheduled_end_time > now()
        WHEN 'finished' THEN talk_sessions.scheduled_end_time <= now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
        THEN talk_sessions.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END;

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
    talk_session_locations.talk_session_id as location_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
FROM talk_sessions
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON  oc.talk_session_id = talk_sessions.talk_session_id
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
LEFT JOIN votes
    ON votes.talk_session_id = talk_sessions.talk_session_id
LEFT JOIN talk_session_locations
    ON talk_session_locations.talk_session_id = talk_sessions.talk_session_id
WHERE
    votes.user_id = sqlc.narg('user_id')::uuid
    AND
    CASE sqlc.narg('status')::text IS NOT NULL
        WHEN sqlc.narg('status')::text = 'finished' THEN talk_sessions.scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN talk_sessions.scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
            THEN talk_sessions.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END
GROUP BY talk_sessions.talk_session_id, oc.opinion_count, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
LIMIT $1 OFFSET $2;
