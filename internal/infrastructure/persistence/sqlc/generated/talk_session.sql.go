// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: talk_session.sql

package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const countTalkSessions = `-- name: CountTalkSessions :one
SELECT
    COUNT(DISTINCT talk_sessions.talk_session_id) AS talk_session_count,
    $1::text AS status
FROM talk_sessions
LEFT JOIN talk_session_locations
    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
LEFT JOIN votes
    ON votes.talk_session_id = talk_sessions.talk_session_id
WHERE
    CASE
        WHEN $2::uuid IS NOT NULL
            THEN votes.user_id = $2::uuid
        ELSE TRUE
    END
    AND
    CASE $1::text
        WHEN 'open' THEN talk_sessions.scheduled_end_time > now()
        WHEN 'finished' THEN talk_sessions.scheduled_end_time <= now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN $3::text IS NOT NULL
        THEN talk_sessions.theme LIKE '%' || $3::text || '%'
        ELSE TRUE
    END
`

type CountTalkSessionsParams struct {
	Status sql.NullString
	UserID uuid.NullUUID
	Theme  sql.NullString
}

type CountTalkSessionsRow struct {
	TalkSessionCount int64
	Status           string
}

// talk_session_locationsがない場合も考慮
//
//	SELECT
//	    COUNT(DISTINCT talk_sessions.talk_session_id) AS talk_session_count,
//	    $1::text AS status
//	FROM talk_sessions
//	LEFT JOIN talk_session_locations
//	    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
//	LEFT JOIN votes
//	    ON votes.talk_session_id = talk_sessions.talk_session_id
//	WHERE
//	    CASE
//	        WHEN $2::uuid IS NOT NULL
//	            THEN votes.user_id = $2::uuid
//	        ELSE TRUE
//	    END
//	    AND
//	    CASE $1::text
//	        WHEN 'open' THEN talk_sessions.scheduled_end_time > now()
//	        WHEN 'finished' THEN talk_sessions.scheduled_end_time <= now()
//	        ELSE TRUE
//	    END
//	    AND
//	    CASE
//	        WHEN $3::text IS NOT NULL
//	        THEN talk_sessions.theme LIKE '%' || $3::text || '%'
//	        ELSE TRUE
//	    END
func (q *Queries) CountTalkSessions(ctx context.Context, arg CountTalkSessionsParams) (CountTalkSessionsRow, error) {
	row := q.db.QueryRowContext(ctx, countTalkSessions, arg.Status, arg.UserID, arg.Theme)
	var i CountTalkSessionsRow
	err := row.Scan(&i.TalkSessionCount, &i.Status)
	return i, err
}

const createTalkSession = `-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, description, owner_id, scheduled_end_time, created_at, city, prefecture) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreateTalkSessionParams struct {
	TalkSessionID    uuid.UUID
	Theme            string
	Description      sql.NullString
	OwnerID          uuid.UUID
	ScheduledEndTime time.Time
	CreatedAt        time.Time
	City             sql.NullString
	Prefecture       sql.NullString
}

// CreateTalkSession
//
//	INSERT INTO talk_sessions (talk_session_id, theme, description, owner_id, scheduled_end_time, created_at, city, prefecture) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
func (q *Queries) CreateTalkSession(ctx context.Context, arg CreateTalkSessionParams) error {
	_, err := q.db.ExecContext(ctx, createTalkSession,
		arg.TalkSessionID,
		arg.Theme,
		arg.Description,
		arg.OwnerID,
		arg.ScheduledEndTime,
		arg.CreatedAt,
		arg.City,
		arg.Prefecture,
	)
	return err
}

const createTalkSessionLocation = `-- name: CreateTalkSessionLocation :exec
INSERT INTO talk_session_locations (talk_session_id, location) VALUES ($1, ST_GeographyFromText($2))
`

type CreateTalkSessionLocationParams struct {
	TalkSessionID       uuid.UUID
	StGeographyfromtext interface{}
}

// CreateTalkSessionLocation
//
//	INSERT INTO talk_session_locations (talk_session_id, location) VALUES ($1, ST_GeographyFromText($2))
func (q *Queries) CreateTalkSessionLocation(ctx context.Context, arg CreateTalkSessionLocationParams) error {
	_, err := q.db.ExecContext(ctx, createTalkSessionLocation, arg.TalkSessionID, arg.StGeographyfromtext)
	return err
}

const editTalkSession = `-- name: EditTalkSession :exec
UPDATE talk_sessions
    SET theme = $2,
        description = $3,
        scheduled_end_time = $4
    WHERE talk_session_id = $1
`

type EditTalkSessionParams struct {
	TalkSessionID    uuid.UUID
	Theme            string
	Description      sql.NullString
	ScheduledEndTime time.Time
}

// EditTalkSession
//
//	UPDATE talk_sessions
//	    SET theme = $2,
//	        description = $3,
//	        scheduled_end_time = $4
//	    WHERE talk_session_id = $1
func (q *Queries) EditTalkSession(ctx context.Context, arg EditTalkSessionParams) error {
	_, err := q.db.ExecContext(ctx, editTalkSession,
		arg.TalkSessionID,
		arg.Theme,
		arg.Description,
		arg.ScheduledEndTime,
	)
	return err
}

const getOwnTalkSessionByUserID = `-- name: GetOwnTalkSessionByUserID :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.description,
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
LEFT JOIN talk_session_locations
    ON talk_session_locations.talk_session_id = talk_sessions.talk_session_id
WHERE
    talk_sessions.owner_id = $3::uuid
    AND
    CASE $4::text
        WHEN 'finished' THEN talk_sessions.scheduled_end_time <= now()
        WHEN 'open' THEN talk_sessions.scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN $5::text IS NOT NULL
            THEN talk_sessions.theme LIKE '%' || $5::text || '%'
        ELSE TRUE
    END
GROUP BY talk_sessions.talk_session_id, oc.opinion_count, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
ORDER BY talk_sessions.created_at DESC
LIMIT $1 OFFSET $2
`

type GetOwnTalkSessionByUserIDParams struct {
	Limit  int32
	Offset int32
	UserID uuid.NullUUID
	Status sql.NullString
	Theme  sql.NullString
}

type GetOwnTalkSessionByUserIDRow struct {
	TalkSessionID    uuid.UUID
	Theme            string
	Description      sql.NullString
	ScheduledEndTime time.Time
	City             sql.NullString
	Prefecture       sql.NullString
	CreatedAt        time.Time
	OpinionCount     int64
	DisplayName      sql.NullString
	DisplayID        sql.NullString
	IconUrl          sql.NullString
	LocationID       uuid.NullUUID
	Latitude         float64
	Longitude        float64
}

// GetOwnTalkSessionByUserID
//
//	SELECT
//	    talk_sessions.talk_session_id,
//	    talk_sessions.theme,
//	    talk_sessions.description,
//	    talk_sessions.scheduled_end_time,
//	    talk_sessions.city AS city,
//	    talk_sessions.prefecture AS prefecture,
//	    talk_sessions.created_at,
//	    COALESCE(oc.opinion_count, 0) AS opinion_count,
//	    users.display_name AS display_name,
//	    users.display_id AS display_id,
//	    users.icon_url AS icon_url,
//	    talk_session_locations.talk_session_id as location_id,
//	    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
//	    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
//	FROM talk_sessions
//	LEFT JOIN (
//	    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
//	    FROM opinions
//	    GROUP BY talk_session_id
//	) oc ON  oc.talk_session_id = talk_sessions.talk_session_id
//	LEFT JOIN users
//	    ON talk_sessions.owner_id = users.user_id
//	LEFT JOIN talk_session_locations
//	    ON talk_session_locations.talk_session_id = talk_sessions.talk_session_id
//	WHERE
//	    talk_sessions.owner_id = $3::uuid
//	    AND
//	    CASE $4::text
//	        WHEN 'finished' THEN talk_sessions.scheduled_end_time <= now()
//	        WHEN 'open' THEN talk_sessions.scheduled_end_time > now()
//	        ELSE TRUE
//	    END
//	    AND
//	    CASE
//	        WHEN $5::text IS NOT NULL
//	            THEN talk_sessions.theme LIKE '%' || $5::text || '%'
//	        ELSE TRUE
//	    END
//	GROUP BY talk_sessions.talk_session_id, oc.opinion_count, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
//	ORDER BY talk_sessions.created_at DESC
//	LIMIT $1 OFFSET $2
func (q *Queries) GetOwnTalkSessionByUserID(ctx context.Context, arg GetOwnTalkSessionByUserIDParams) ([]GetOwnTalkSessionByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getOwnTalkSessionByUserID,
		arg.Limit,
		arg.Offset,
		arg.UserID,
		arg.Status,
		arg.Theme,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOwnTalkSessionByUserIDRow
	for rows.Next() {
		var i GetOwnTalkSessionByUserIDRow
		if err := rows.Scan(
			&i.TalkSessionID,
			&i.Theme,
			&i.Description,
			&i.ScheduledEndTime,
			&i.City,
			&i.Prefecture,
			&i.CreatedAt,
			&i.OpinionCount,
			&i.DisplayName,
			&i.DisplayID,
			&i.IconUrl,
			&i.LocationID,
			&i.Latitude,
			&i.Longitude,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRespondTalkSessionByUserID = `-- name: GetRespondTalkSessionByUserID :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.description,
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
    votes.user_id = $3::uuid
    AND
    CASE $4::text IS NOT NULL
        WHEN $4::text = 'finished' THEN talk_sessions.scheduled_end_time <= now()
        WHEN $4::text = 'open' THEN talk_sessions.scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    CASE
        WHEN $5::text IS NOT NULL
            THEN talk_sessions.theme LIKE '%' || $5::text || '%'
        ELSE TRUE
    END
GROUP BY talk_sessions.talk_session_id, oc.opinion_count, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
ORDER BY talk_sessions.created_at DESC
LIMIT $1 OFFSET $2
`

type GetRespondTalkSessionByUserIDParams struct {
	Limit  int32
	Offset int32
	UserID uuid.NullUUID
	Status sql.NullString
	Theme  sql.NullString
}

type GetRespondTalkSessionByUserIDRow struct {
	TalkSessionID    uuid.UUID
	Theme            string
	Description      sql.NullString
	ScheduledEndTime time.Time
	City             sql.NullString
	Prefecture       sql.NullString
	CreatedAt        time.Time
	OpinionCount     int64
	DisplayName      sql.NullString
	DisplayID        sql.NullString
	IconUrl          sql.NullString
	LocationID       uuid.NullUUID
	Latitude         float64
	Longitude        float64
}

// GetRespondTalkSessionByUserID
//
//	SELECT
//	    talk_sessions.talk_session_id,
//	    talk_sessions.theme,
//	    talk_sessions.description,
//	    talk_sessions.scheduled_end_time,
//	    talk_sessions.city AS city,
//	    talk_sessions.prefecture AS prefecture,
//	    talk_sessions.created_at,
//	    COALESCE(oc.opinion_count, 0) AS opinion_count,
//	    users.display_name AS display_name,
//	    users.display_id AS display_id,
//	    users.icon_url AS icon_url,
//	    talk_session_locations.talk_session_id as location_id,
//	    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
//	    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
//	FROM talk_sessions
//	LEFT JOIN (
//	    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
//	    FROM opinions
//	    GROUP BY talk_session_id
//	) oc ON  oc.talk_session_id = talk_sessions.talk_session_id
//	LEFT JOIN users
//	    ON talk_sessions.owner_id = users.user_id
//	LEFT JOIN votes
//	    ON votes.talk_session_id = talk_sessions.talk_session_id
//	LEFT JOIN talk_session_locations
//	    ON talk_session_locations.talk_session_id = talk_sessions.talk_session_id
//	WHERE
//	    votes.user_id = $3::uuid
//	    AND
//	    CASE $4::text IS NOT NULL
//	        WHEN $4::text = 'finished' THEN talk_sessions.scheduled_end_time <= now()
//	        WHEN $4::text = 'open' THEN talk_sessions.scheduled_end_time > now()
//	        ELSE TRUE
//	    END
//	    AND
//	    CASE
//	        WHEN $5::text IS NOT NULL
//	            THEN talk_sessions.theme LIKE '%' || $5::text || '%'
//	        ELSE TRUE
//	    END
//	GROUP BY talk_sessions.talk_session_id, oc.opinion_count, users.display_name, users.display_id, users.icon_url, talk_session_locations.talk_session_id
//	ORDER BY talk_sessions.created_at DESC
//	LIMIT $1 OFFSET $2
func (q *Queries) GetRespondTalkSessionByUserID(ctx context.Context, arg GetRespondTalkSessionByUserIDParams) ([]GetRespondTalkSessionByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getRespondTalkSessionByUserID,
		arg.Limit,
		arg.Offset,
		arg.UserID,
		arg.Status,
		arg.Theme,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRespondTalkSessionByUserIDRow
	for rows.Next() {
		var i GetRespondTalkSessionByUserIDRow
		if err := rows.Scan(
			&i.TalkSessionID,
			&i.Theme,
			&i.Description,
			&i.ScheduledEndTime,
			&i.City,
			&i.Prefecture,
			&i.CreatedAt,
			&i.OpinionCount,
			&i.DisplayName,
			&i.DisplayID,
			&i.IconUrl,
			&i.LocationID,
			&i.Latitude,
			&i.Longitude,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTalkSessionByID = `-- name: GetTalkSessionByID :one
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.description,
    talk_sessions.created_at,
    talk_sessions.scheduled_end_time,
    talk_sessions.city AS city,
    talk_sessions.prefecture AS prefecture,
    COALESCE(oc.opinion_count, 0) AS opinion_count,
    users.user_id AS user_id,
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
WHERE talk_sessions.talk_session_id = $1
`

type GetTalkSessionByIDRow struct {
	TalkSessionID    uuid.UUID
	Theme            string
	Description      sql.NullString
	CreatedAt        time.Time
	ScheduledEndTime time.Time
	City             sql.NullString
	Prefecture       sql.NullString
	OpinionCount     int64
	UserID           uuid.NullUUID
	DisplayName      sql.NullString
	DisplayID        sql.NullString
	IconUrl          sql.NullString
	LocationID       uuid.NullUUID
	Latitude         float64
	Longitude        float64
}

// GetTalkSessionByID
//
//	SELECT
//	    talk_sessions.talk_session_id,
//	    talk_sessions.theme,
//	    talk_sessions.description,
//	    talk_sessions.created_at,
//	    talk_sessions.scheduled_end_time,
//	    talk_sessions.city AS city,
//	    talk_sessions.prefecture AS prefecture,
//	    COALESCE(oc.opinion_count, 0) AS opinion_count,
//	    users.user_id AS user_id,
//	    users.display_name AS display_name,
//	    users.display_id AS display_id,
//	    users.icon_url AS icon_url,
//	    talk_session_locations.talk_session_id as location_id,
//	    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
//	    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude
//	FROM talk_sessions
//	LEFT JOIN users
//	    ON talk_sessions.owner_id = users.user_id
//	LEFT JOIN (
//	    SELECT opinions.talk_session_id, COUNT(opinions.opinion_id) AS opinion_count
//	    FROM opinions
//	    GROUP BY opinions.talk_session_id
//	) oc ON talk_sessions.talk_session_id = oc.talk_session_id
//	LEFT JOIN talk_session_locations
//	    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
//	WHERE talk_sessions.talk_session_id = $1
func (q *Queries) GetTalkSessionByID(ctx context.Context, talkSessionID uuid.UUID) (GetTalkSessionByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getTalkSessionByID, talkSessionID)
	var i GetTalkSessionByIDRow
	err := row.Scan(
		&i.TalkSessionID,
		&i.Theme,
		&i.Description,
		&i.CreatedAt,
		&i.ScheduledEndTime,
		&i.City,
		&i.Prefecture,
		&i.OpinionCount,
		&i.UserID,
		&i.DisplayName,
		&i.DisplayID,
		&i.IconUrl,
		&i.LocationID,
		&i.Latitude,
		&i.Longitude,
	)
	return i, err
}

const listTalkSessions = `-- name: ListTalkSessions :many
SELECT
    talk_sessions.talk_session_id,
    talk_sessions.theme,
    talk_sessions.description,
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
    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude,
    CASE
        WHEN $3::float IS NOT NULL AND $4::float IS NOT NULL AND talk_session_locations.location IS NOT NULL
            THEN ('SRID=4326;POINT(' ||
            ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
            ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
            <->
            ('SRID=4326;POINT(' || $3::float || ' ' || $4::float || ')')::geometry
        ELSE NULL
    END AS distance
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
        WHEN $5::text = 'finished' THEN scheduled_end_time <= now()
        WHEN $5::text = 'open' THEN scheduled_end_time > now()
        ELSE TRUE
    END
    AND
    (CASE
        WHEN $6::text IS NOT NULL
        THEN talk_sessions.theme LIKE '%' || $6::text || '%'
        ELSE TRUE
    END)
    AND
    (CASE $7::text
            WHEN 'nearlest' THEN
                $3::float IS NOT NULL AND $4::float IS NOT NULL
                AND
                ('SRID=4326;POINT(' ||
                ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
                ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
                <->
                ('SRID=4326;POINT(' || $3::float || ' ' || $4::float || ')')::geometry <= 100000
            ELSE TRUE
    END)
ORDER BY
    CASE $7::text
        WHEN 'oldest' THEN (EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM talk_sessions.created_at))*-1
        WHEN 'mostReplies' THEN -oc.opinion_count
        WHEN 'nearest' THEN (
        CASE
            WHEN $3::float IS NOT NULL AND $4::float IS NOT NULL AND talk_session_locations.location IS NOT NULL
                THEN ('SRID=4326;POINT(' ||
                ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
                ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
                <->
                ('SRID=4326;POINT(' || $3::float || ' ' || $4::float || ')')::geometry
            ELSE NULL
        END
        )
        ELSE EXTRACT(EPOCH FROM talk_sessions.created_at)*-1
    END ASC
LIMIT $1 OFFSET $2
`

type ListTalkSessionsParams struct {
	Limit     int32
	Offset    int32
	Latitude  sql.NullFloat64
	Longitude sql.NullFloat64
	Status    sql.NullString
	Theme     sql.NullString
	SortKey   sql.NullString
}

type ListTalkSessionsRow struct {
	TalkSessionID    uuid.UUID
	Theme            string
	Description      sql.NullString
	ScheduledEndTime time.Time
	City             sql.NullString
	Prefecture       sql.NullString
	CreatedAt        time.Time
	OpinionCount     int64
	DisplayName      sql.NullString
	DisplayID        sql.NullString
	IconUrl          sql.NullString
	LocationID       uuid.NullUUID
	Latitude         float64
	Longitude        float64
	Distance         interface{}
}

// ListTalkSessions
//
//	SELECT
//	    talk_sessions.talk_session_id,
//	    talk_sessions.theme,
//	    talk_sessions.description,
//	    talk_sessions.scheduled_end_time,
//	    talk_sessions.city AS city,
//	    talk_sessions.prefecture AS prefecture,
//	    talk_sessions.created_at,
//	    COALESCE(oc.opinion_count, 0) AS opinion_count,
//	    users.display_name AS display_name,
//	    users.display_id AS display_id,
//	    users.icon_url AS icon_url,
//	    talk_session_locations.talk_session_id as location_id,
//	    COALESCE(ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS latitude,
//	    COALESCE(ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location))),0)::float AS longitude,
//	    CASE
//	        WHEN $3::float IS NOT NULL AND $4::float IS NOT NULL AND talk_session_locations.location IS NOT NULL
//	            THEN ('SRID=4326;POINT(' ||
//	            ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
//	            ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
//	            <->
//	            ('SRID=4326;POINT(' || $3::float || ' ' || $4::float || ')')::geometry
//	        ELSE NULL
//	    END AS distance
//	FROM talk_sessions
//	LEFT JOIN (
//	    SELECT talk_session_id, COUNT(opinion_id) AS opinion_count
//	    FROM opinions
//	    GROUP BY talk_session_id
//	) oc ON talk_sessions.talk_session_id = oc.talk_session_id
//	LEFT JOIN users
//	    ON talk_sessions.owner_id = users.user_id
//	LEFT JOIN talk_session_locations
//	    ON talk_sessions.talk_session_id = talk_session_locations.talk_session_id
//	WHERE
//	    CASE
//	        WHEN $5::text = 'finished' THEN scheduled_end_time <= now()
//	        WHEN $5::text = 'open' THEN scheduled_end_time > now()
//	        ELSE TRUE
//	    END
//	    AND
//	    (CASE
//	        WHEN $6::text IS NOT NULL
//	        THEN talk_sessions.theme LIKE '%' || $6::text || '%'
//	        ELSE TRUE
//	    END)
//	    AND
//	    (CASE $7::text
//	            WHEN 'nearlest' THEN
//	                $3::float IS NOT NULL AND $4::float IS NOT NULL
//	                AND
//	                ('SRID=4326;POINT(' ||
//	                ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
//	                ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
//	                <->
//	                ('SRID=4326;POINT(' || $3::float || ' ' || $4::float || ')')::geometry <= 100000
//	            ELSE TRUE
//	    END)
//	ORDER BY
//	    CASE $7::text
//	        WHEN 'oldest' THEN (EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM talk_sessions.created_at))*-1
//	        WHEN 'mostReplies' THEN -oc.opinion_count
//	        WHEN 'nearest' THEN (
//	        CASE
//	            WHEN $3::float IS NOT NULL AND $4::float IS NOT NULL AND talk_session_locations.location IS NOT NULL
//	                THEN ('SRID=4326;POINT(' ||
//	                ST_Y(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ' ' ||
//	                ST_X(ST_GeomFromWKB(ST_AsBinary(talk_session_locations.location),4326)) || ')')::geometry
//	                <->
//	                ('SRID=4326;POINT(' || $3::float || ' ' || $4::float || ')')::geometry
//	            ELSE NULL
//	        END
//	        )
//	        ELSE EXTRACT(EPOCH FROM talk_sessions.created_at)*-1
//	    END ASC
//	LIMIT $1 OFFSET $2
func (q *Queries) ListTalkSessions(ctx context.Context, arg ListTalkSessionsParams) ([]ListTalkSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listTalkSessions,
		arg.Limit,
		arg.Offset,
		arg.Latitude,
		arg.Longitude,
		arg.Status,
		arg.Theme,
		arg.SortKey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTalkSessionsRow
	for rows.Next() {
		var i ListTalkSessionsRow
		if err := rows.Scan(
			&i.TalkSessionID,
			&i.Theme,
			&i.Description,
			&i.ScheduledEndTime,
			&i.City,
			&i.Prefecture,
			&i.CreatedAt,
			&i.OpinionCount,
			&i.DisplayName,
			&i.DisplayID,
			&i.IconUrl,
			&i.LocationID,
			&i.Latitude,
			&i.Longitude,
			&i.Distance,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTalkSessionLocation = `-- name: UpdateTalkSessionLocation :exec
UPDATE talk_session_locations SET location = ST_GeographyFromText($2) WHERE talk_session_id = $1
`

type UpdateTalkSessionLocationParams struct {
	TalkSessionID       uuid.UUID
	StGeographyfromtext interface{}
}

// UpdateTalkSessionLocation
//
//	UPDATE talk_session_locations SET location = ST_GeographyFromText($2) WHERE talk_session_id = $1
func (q *Queries) UpdateTalkSessionLocation(ctx context.Context, arg UpdateTalkSessionLocationParams) error {
	_, err := q.db.ExecContext(ctx, updateTalkSessionLocation, arg.TalkSessionID, arg.StGeographyfromtext)
	return err
}
