-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, description, thumbnail_url, owner_id, scheduled_end_time, created_at, city, prefecture, restrictions, hide_report, organization_id, organization_alias_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: CreateTalkSessionLocation :exec
INSERT INTO talk_session_locations (talk_session_id, location) VALUES ($1, ST_GeographyFromText($2));

-- name: UpdateTalkSessionLocation :exec
UPDATE talk_session_locations SET location = ST_GeographyFromText($2) WHERE talk_session_id = $1;

-- name: UpdateTalkSessionHideReport :exec
UPDATE talk_sessions SET hide_report = $2 WHERE talk_session_id = $1;

-- name: EditTalkSession :exec
UPDATE talk_sessions
    SET theme = $2,
        description = $3,
        scheduled_end_time = $4,
        thumbnail_url = $5,
        city = $6,
        prefecture = $7,
        restrictions = $8,
        hide_report = $9,
        updated_at = NOW(),
        organization_id = $10,
        organization_alias_id = $11
    WHERE talk_session_id = $1;

-- name: GetTalkSessionByID :one
SELECT
    sqlc.embed(ts),
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    sqlc.embed(users),
    COALESCE(organization_aliases.alias_name, '') AS alias_name,
    COALESCE(organization_aliases.alias_id, '00000000-0000-0000-0000-000000000000'::uuid) AS alias_id,
    COALESCE(organization_aliases.organization_id, '00000000-0000-0000-0000-000000000000'::uuid) AS organization_id,
    talk_session_locations.talk_session_id as location_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
FROM talk_sessions ts
LEFT JOIN users
    ON ts.owner_id = users.user_id
LEFT JOIN (
    SELECT opinions.talk_session_id, COUNT(opinions.opinion_id) AS opinion_count
    FROM opinions
    GROUP BY opinions.talk_session_id
) oc ON ts.talk_session_id = oc.talk_session_id
LEFT JOIN organization_aliases
    ON ts.organization_alias_id = organization_aliases.alias_id
LEFT JOIN talk_session_locations
    ON ts.talk_session_id = talk_session_locations.talk_session_id
WHERE ts.talk_session_id = $1;

-- name: ListTalkSessions :many
SELECT
    sqlc.embed(ts),
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    sqlc.embed(users),
    COALESCE(organization_aliases.alias_name, '') AS alias_name,
    COALESCE(organization_aliases.alias_id, '00000000-0000-0000-0000-000000000000'::uuid) AS alias_id,
    COALESCE(organization_aliases.organization_id, '00000000-0000-0000-0000-000000000000'::uuid) AS organization_id,
    COALESCE(votes.vote_count, 0) AS vote_count,
    COALESCE(vote_users.vote_count, 0) AS vote_user_count,
    talk_session_locations.talk_session_id as location_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude,
    CASE
        WHEN sqlc.narg('latitude')::float IS NOT NULL AND sqlc.narg('longitude')::float IS NOT NULL AND talk_session_locations.location IS NOT NULL
            THEN ('SRID=4326;POINT(' ||
            ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
            ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
            <->
            ('SRID=4326;POINT(' || sqlc.narg('latitude')::float || ' ' || sqlc.narg('longitude')::float || ')')::geometry
        ELSE NULL
    END AS distance
FROM talk_sessions ts
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON ts.talk_session_id = oc.talk_session_id
LEFT JOIN users
    ON ts.owner_id = users.user_id
LEFT JOIN (
    SELECT talk_session_id, COUNT(DISTINCT vote_id) AS vote_count
    FROM votes
    GROUP BY talk_session_id
) votes ON ts.talk_session_id = votes.talk_session_id
LEFT JOIN (
    SELECT talk_session_id, COUNT(DISTINCT user_id) AS vote_count
    FROM votes
    GROUP BY talk_session_id
) vote_users ON ts.talk_session_id = vote_users.talk_session_id
LEFT JOIN organization_aliases
    ON ts.organization_alias_id = organization_aliases.alias_id
LEFT JOIN talk_session_locations
    ON ts.talk_session_id = talk_session_locations.talk_session_id
WHERE
    CASE sqlc.narg('status')::text
        WHEN 'finished' THEN ts.scheduled_end_time <= now()
        WHEN 'open' THEN ts.scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    (CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
        THEN ts.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END)
    AND
    (CASE sqlc.narg('sort_key')::text
            WHEN 'nearlest' THEN
                sqlc.narg('latitude')::float IS NOT NULL AND sqlc.narg('longitude')::float IS NOT NULL
                AND
                ('SRID=4326;POINT(' ||
                ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
                ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
                <->
                ('SRID=4326;POINT(' || sqlc.narg('latitude')::float || ' ' || sqlc.narg('longitude')::float || ')')::geometry <= 100000
            ELSE TRUE
    END)
ORDER BY
    CASE sqlc.narg('sort_key')::text
        WHEN 'oldest' THEN (EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM ts.created_at))*-1
        WHEN 'mostReplies' THEN -oc.opinion_count
        WHEN 'nearest' THEN (
        CASE
            WHEN sqlc.narg('latitude')::float IS NOT NULL AND sqlc.narg('longitude')::float IS NOT NULL AND talk_session_locations.location IS NOT NULL
                THEN ('SRID=4326;POINT(' ||
                ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
                ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
                <->
                ('SRID=4326;POINT(' || sqlc.narg('latitude')::float || ' ' || sqlc.narg('longitude')::float || ')')::geometry
            ELSE NULL
        END
        )
        ELSE EXTRACT(EPOCH FROM ts.created_at)*-1
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

-- name: GetOwnTalkSessionByUserID :many
SELECT
    sqlc.embed(ts),
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    sqlc.embed(users),
    COALESCE(organization_aliases.alias_name, '') AS alias_name,
    COALESCE(organization_aliases.alias_id, '00000000-0000-0000-0000-000000000000'::uuid) AS alias_id,
    COALESCE(organization_aliases.organization_id, '00000000-0000-0000-0000-000000000000'::uuid) AS organization_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
FROM talk_sessions ts
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON  oc.talk_session_id = ts.talk_session_id
LEFT JOIN users
    ON ts.owner_id = users.user_id
LEFT JOIN talk_session_locations
    ON talk_session_locations.talk_session_id = ts.talk_session_id
LEFT JOIN organization_aliases
    ON ts.organization_alias_id = organization_aliases.alias_id
WHERE
    ts.owner_id = sqlc.narg('user_id')::uuid
    AND
    CASE sqlc.narg('status')::text
        WHEN 'finished' THEN ts.scheduled_end_time <= now()
        WHEN 'open' THEN ts.scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
            THEN ts.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END
GROUP BY ts.talk_session_id, oc.opinion_count, users.user_id, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
ORDER BY ts.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRespondTalkSessionByUserID :many
SELECT
    sqlc.embed(ts),
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    sqlc.embed(users),
    COALESCE(organization_aliases.alias_name, '') AS alias_name,
    COALESCE(organization_aliases.alias_id, '00000000-0000-0000-0000-000000000000'::uuid) AS alias_id,
    COALESCE(organization_aliases.organization_id, '00000000-0000-0000-0000-000000000000'::uuid) AS organization_id,
    organization_aliases.created_at,
    talk_session_locations.talk_session_id as location_id,
    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
FROM talk_sessions ts
LEFT JOIN (
    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
    FROM opinions
    GROUP BY talk_session_id
) oc ON  oc.talk_session_id = ts.talk_session_id
LEFT JOIN users
    ON ts.owner_id = users.user_id
LEFT JOIN votes
    ON votes.talk_session_id = ts.talk_session_id
LEFT JOIN talk_session_locations
    ON talk_session_locations.talk_session_id = ts.talk_session_id
LEFT JOIN organization_aliases
    ON ts.organization_alias_id = organization_aliases.alias_id
WHERE
    votes.user_id = sqlc.narg('user_id')::uuid
    AND
    CASE sqlc.narg('status')::text IS NOT NULL
        WHEN sqlc.narg('status')::text = 'finished' THEN ts.scheduled_end_time <= now()
        WHEN sqlc.narg('status')::text = 'open' THEN ts.scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('theme')::text IS NOT NULL
            THEN ts.theme LIKE '%' || sqlc.narg('theme')::text || '%'
        ELSE TRUE
    END
GROUP BY ts.talk_session_id, oc.opinion_count, users.user_id, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
ORDER BY ts.created_at DESC
LIMIT $1 OFFSET $2;


-- name: GetAllTalkSessionCount :one
SELECT
    COUNT(DISTINCT talk_sessions.talk_session_id) AS talk_session_count
FROM talk_sessions;
