// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: opinion.sql

package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const countOpinions = `-- name: CountOpinions :one
SELECT
    COUNT(opinions.*) AS opinion_count
FROM opinions
WHERE
    CASE
        WHEN $1::uuid IS NOT NULL THEN opinions.user_id = $1::uuid
        ELSE TRUE
    END
    AND
    CASE
        WHEN $2::uuid IS NOT NULL THEN opinions.talk_session_id = $2::uuid
        ELSE TRUE
    END
    AND
    CASE
        WHEN $3::uuid IS NOT NULL THEN opinions.parent_opinion_id = $3::uuid
        ELSE TRUE
    END
`

type CountOpinionsParams struct {
	UserID          uuid.NullUUID
	TalkSessionID   uuid.NullUUID
	ParentOpinionID uuid.NullUUID
}

// CountOpinions
//
//	SELECT
//	    COUNT(opinions.*) AS opinion_count
//	FROM opinions
//	WHERE
//	    CASE
//	        WHEN $1::uuid IS NOT NULL THEN opinions.user_id = $1::uuid
//	        ELSE TRUE
//	    END
//	    AND
//	    CASE
//	        WHEN $2::uuid IS NOT NULL THEN opinions.talk_session_id = $2::uuid
//	        ELSE TRUE
//	    END
//	    AND
//	    CASE
//	        WHEN $3::uuid IS NOT NULL THEN opinions.parent_opinion_id = $3::uuid
//	        ELSE TRUE
//	    END
func (q *Queries) CountOpinions(ctx context.Context, arg CountOpinionsParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countOpinions, arg.UserID, arg.TalkSessionID, arg.ParentOpinionID)
	var opinion_count int64
	err := row.Scan(&opinion_count)
	return opinion_count, err
}

const createOpinion = `-- name: CreateOpinion :exec
INSERT INTO opinions (
    opinion_id,
    talk_session_id,
    user_id,
    parent_opinion_id,
    title,
    content,
    reference_url,
    picture_url,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

type CreateOpinionParams struct {
	OpinionID       uuid.UUID
	TalkSessionID   uuid.UUID
	UserID          uuid.UUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         string
	ReferenceUrl    sql.NullString
	PictureUrl      sql.NullString
	CreatedAt       time.Time
}

// CreateOpinion
//
//	INSERT INTO opinions (
//	    opinion_id,
//	    talk_session_id,
//	    user_id,
//	    parent_opinion_id,
//	    title,
//	    content,
//	    reference_url,
//	    picture_url,
//	    created_at
//	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
func (q *Queries) CreateOpinion(ctx context.Context, arg CreateOpinionParams) error {
	_, err := q.db.ExecContext(ctx, createOpinion,
		arg.OpinionID,
		arg.TalkSessionID,
		arg.UserID,
		arg.ParentOpinionID,
		arg.Title,
		arg.Content,
		arg.ReferenceUrl,
		arg.PictureUrl,
		arg.CreatedAt,
	)
	return err
}

const getOpinionByID = `-- name: GetOpinionByID :one
SELECT
    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.created_at, opinions.picture_url, opinions.reference_url,
    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON pv.opinion_id = opinions.parent_opinion_id
    AND  pv.user_id = opinions.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
    WHERE votes.user_id = $2::uuid
) cv ON opinions.opinion_id = cv.opinion_id
WHERE opinions.opinion_id = $1
`

type GetOpinionByIDParams struct {
	OpinionID uuid.UUID
	UserID    uuid.NullUUID
}

type GetOpinionByIDRow struct {
	Opinion         Opinion
	User            User
	ParentVoteType  int16
	CurrentVoteType int16
}

// 親意見に対するユーザーの投票を取得
// ユーザーIDが提供された場合、そのユーザーの投票ステータスを一緒に取得
//
//	SELECT
//	    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.created_at, opinions.picture_url, opinions.reference_url,
//	    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
//	    COALESCE(pv.vote_type, 0) AS parent_vote_type,
//	    COALESCE(cv.vote_type, 0) AS current_vote_type
//	FROM opinions
//	LEFT JOIN users
//	    ON opinions.user_id = users.user_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.user_id, votes.opinion_id
//	    FROM votes
//	) pv ON pv.opinion_id = opinions.parent_opinion_id
//	    AND  pv.user_id = opinions.user_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.user_id, votes.opinion_id
//	    FROM votes
//	    WHERE votes.user_id = $2::uuid
//	) cv ON opinions.opinion_id = cv.opinion_id
//	WHERE opinions.opinion_id = $1
func (q *Queries) GetOpinionByID(ctx context.Context, arg GetOpinionByIDParams) (GetOpinionByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getOpinionByID, arg.OpinionID, arg.UserID)
	var i GetOpinionByIDRow
	err := row.Scan(
		&i.Opinion.OpinionID,
		&i.Opinion.TalkSessionID,
		&i.Opinion.UserID,
		&i.Opinion.ParentOpinionID,
		&i.Opinion.Title,
		&i.Opinion.Content,
		&i.Opinion.CreatedAt,
		&i.Opinion.PictureUrl,
		&i.Opinion.ReferenceUrl,
		&i.User.UserID,
		&i.User.DisplayID,
		&i.User.DisplayName,
		&i.User.IconUrl,
		&i.User.CreatedAt,
		&i.User.UpdatedAt,
		&i.ParentVoteType,
		&i.CurrentVoteType,
	)
	return i, err
}

const getOpinionReplies = `-- name: GetOpinionReplies :many
SELECT
    DISTINCT opinions.opinion_id,
    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.created_at, opinions.picture_url, opinions.reference_url,
    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
    WHERE votes.opinion_id = $1
) pv ON opinions.user_id = pv.user_id
    AND opinions.opinion_id = pv.opinion_id
LEFT JOIN (
    SELECT votes.vote_type, votes.opinion_id
    FROM votes
    WHERE votes.user_id = $2::uuid
) cv ON opinions.opinion_id = cv.opinion_id
WHERE opinions.parent_opinion_id = $1
GROUP BY opinions.opinion_id, users.display_name, users.display_id, users.icon_url, pv.vote_type, cv.vote_type
ORDER BY opinions.created_at DESC
`

type GetOpinionRepliesParams struct {
	OpinionID uuid.UUID
	UserID    uuid.NullUUID
}

type GetOpinionRepliesRow struct {
	OpinionID       uuid.UUID
	Opinion         Opinion
	User            User
	ParentVoteType  int16
	CurrentVoteType int16
}

// 親意見に対する子意見主の投票を取得
// ユーザーIDが提供された場合、そのユーザーの投票ステータスを取得
//
//	SELECT
//	    DISTINCT opinions.opinion_id,
//	    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.created_at, opinions.picture_url, opinions.reference_url,
//	    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
//	    COALESCE(pv.vote_type, 0) AS parent_vote_type,
//	    COALESCE(cv.vote_type, 0) AS current_vote_type
//	FROM opinions
//	LEFT JOIN users
//	    ON opinions.user_id = users.user_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.user_id, votes.opinion_id
//	    FROM votes
//	    WHERE votes.opinion_id = $1
//	) pv ON opinions.user_id = pv.user_id
//	    AND opinions.opinion_id = pv.opinion_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.opinion_id
//	    FROM votes
//	    WHERE votes.user_id = $2::uuid
//	) cv ON opinions.opinion_id = cv.opinion_id
//	WHERE opinions.parent_opinion_id = $1
//	GROUP BY opinions.opinion_id, users.display_name, users.display_id, users.icon_url, pv.vote_type, cv.vote_type
//	ORDER BY opinions.created_at DESC
func (q *Queries) GetOpinionReplies(ctx context.Context, arg GetOpinionRepliesParams) ([]GetOpinionRepliesRow, error) {
	rows, err := q.db.QueryContext(ctx, getOpinionReplies, arg.OpinionID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOpinionRepliesRow
	for rows.Next() {
		var i GetOpinionRepliesRow
		if err := rows.Scan(
			&i.OpinionID,
			&i.Opinion.OpinionID,
			&i.Opinion.TalkSessionID,
			&i.Opinion.UserID,
			&i.Opinion.ParentOpinionID,
			&i.Opinion.Title,
			&i.Opinion.Content,
			&i.Opinion.CreatedAt,
			&i.Opinion.PictureUrl,
			&i.Opinion.ReferenceUrl,
			&i.User.UserID,
			&i.User.DisplayID,
			&i.User.DisplayName,
			&i.User.IconUrl,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.ParentVoteType,
			&i.CurrentVoteType,
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

const getOpinionsByTalkSessionID = `-- name: GetOpinionsByTalkSessionID :many
WITH unique_opinions AS (
    SELECT DISTINCT ON (opinions.opinion_id)
        opinions.opinion_id,
        opinions.talk_session_id,
        opinions.user_id,
        opinions.parent_opinion_id,
        opinions.title,
        opinions.content,
        opinions.reference_url,
        opinions.picture_url,
        opinions.created_at
    FROM opinions
    WHERE opinions.talk_session_id = $1
)
SELECT
    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.reference_url, opinions.picture_url, opinions.created_at,
    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(rc.reply_count, 0) AS reply_count,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM unique_opinions opinions
LEFT JOIN users ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT DISTINCT ON (opinion_id) vote_type, user_id, opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
LEFT JOIN (
    SELECT DISTINCT ON (opinion_id) vote_type, user_id, opinion_id
    FROM votes
    WHERE user_id = $4::uuid
) cv ON opinions.opinion_id = cv.opinion_id
ORDER BY
    CASE $5::text
        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'mostReplies' THEN COALESCE(rc.reply_count, 0)
    END DESC
LIMIT $2 OFFSET $3
`

type GetOpinionsByTalkSessionIDParams struct {
	TalkSessionID uuid.UUID
	Limit         int32
	Offset        int32
	UserID        uuid.NullUUID
	SortKey       sql.NullString
}

type GetOpinionsByTalkSessionIDRow struct {
	Opinion         Opinion
	User            User
	ParentVoteType  int16
	ReplyCount      int64
	CurrentVoteType int16
}

// GetOpinionsByTalkSessionID
//
//	WITH unique_opinions AS (
//	    SELECT DISTINCT ON (opinions.opinion_id)
//	        opinions.opinion_id,
//	        opinions.talk_session_id,
//	        opinions.user_id,
//	        opinions.parent_opinion_id,
//	        opinions.title,
//	        opinions.content,
//	        opinions.reference_url,
//	        opinions.picture_url,
//	        opinions.created_at
//	    FROM opinions
//	    WHERE opinions.talk_session_id = $1
//	)
//	SELECT
//	    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.reference_url, opinions.picture_url, opinions.created_at,
//	    users.user_id, users.display_id, users.display_name, users.icon_url, users.created_at, users.updated_at,
//	    COALESCE(pv.vote_type, 0) AS parent_vote_type,
//	    COALESCE(rc.reply_count, 0) AS reply_count,
//	    COALESCE(cv.vote_type, 0) AS current_vote_type
//	FROM unique_opinions opinions
//	LEFT JOIN users ON opinions.user_id = users.user_id
//	LEFT JOIN (
//	    SELECT DISTINCT ON (opinion_id) vote_type, user_id, opinion_id
//	    FROM votes
//	) pv ON opinions.parent_opinion_id = pv.opinion_id
//	    AND opinions.user_id = pv.user_id
//	LEFT JOIN (
//	    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
//	    FROM opinions
//	    GROUP BY parent_opinion_id
//	) rc ON opinions.opinion_id = rc.parent_opinion_id
//	LEFT JOIN (
//	    SELECT DISTINCT ON (opinion_id) vote_type, user_id, opinion_id
//	    FROM votes
//	    WHERE user_id = $4::uuid
//	) cv ON opinions.opinion_id = cv.opinion_id
//	ORDER BY
//	    CASE $5::text
//	        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
//	        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
//	        WHEN 'mostReplies' THEN COALESCE(rc.reply_count, 0)
//	    END DESC
//	LIMIT $2 OFFSET $3
func (q *Queries) GetOpinionsByTalkSessionID(ctx context.Context, arg GetOpinionsByTalkSessionIDParams) ([]GetOpinionsByTalkSessionIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getOpinionsByTalkSessionID,
		arg.TalkSessionID,
		arg.Limit,
		arg.Offset,
		arg.UserID,
		arg.SortKey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOpinionsByTalkSessionIDRow
	for rows.Next() {
		var i GetOpinionsByTalkSessionIDRow
		if err := rows.Scan(
			&i.Opinion.OpinionID,
			&i.Opinion.TalkSessionID,
			&i.Opinion.UserID,
			&i.Opinion.ParentOpinionID,
			&i.Opinion.Title,
			&i.Opinion.Content,
			&i.Opinion.CreatedAt,
			&i.Opinion.PictureUrl,
			&i.Opinion.ReferenceUrl,
			&i.User.UserID,
			&i.User.DisplayID,
			&i.User.DisplayName,
			&i.User.IconUrl,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.ParentVoteType,
			&i.ReplyCount,
			&i.CurrentVoteType,
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

const getOpinionsByUserID = `-- name: GetOpinionsByUserID :many
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.reference_url,
    opinions.picture_url,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    -- 意見に対するリプライ数（再帰）
    COALESCE(rc.reply_count, 0) AS reply_count
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
WHERE opinions.user_id = $1
ORDER BY
    CASE $4::text
        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'mostReplies' THEN reply_count
    END DESC
LIMIT $2 OFFSET $3
`

type GetOpinionsByUserIDParams struct {
	UserID  uuid.UUID
	Limit   int32
	Offset  int32
	SortKey sql.NullString
}

type GetOpinionsByUserIDRow struct {
	OpinionID       uuid.UUID
	TalkSessionID   uuid.UUID
	UserID          uuid.UUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         string
	ReferenceUrl    sql.NullString
	PictureUrl      sql.NullString
	CreatedAt       time.Time
	DisplayName     sql.NullString
	DisplayID       sql.NullString
	IconUrl         sql.NullString
	ParentVoteType  int16
	ReplyCount      int64
}

// latest, mostReply, oldestでソート
//
//	SELECT
//	    opinions.opinion_id,
//	    opinions.talk_session_id,
//	    opinions.user_id,
//	    opinions.parent_opinion_id,
//	    opinions.title,
//	    opinions.content,
//	    opinions.reference_url,
//	    opinions.picture_url,
//	    opinions.created_at,
//	    users.display_name AS display_name,
//	    users.display_id AS display_id,
//	    users.icon_url AS icon_url,
//	    COALESCE(pv.vote_type, 0) AS parent_vote_type,
//	    -- 意見に対するリプライ数（再帰）
//	    COALESCE(rc.reply_count, 0) AS reply_count
//	FROM opinions
//	LEFT JOIN users
//	    ON opinions.user_id = users.user_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.user_id, votes.opinion_id
//	    FROM votes
//	) pv ON opinions.parent_opinion_id = pv.opinion_id
//	    AND opinions.user_id = pv.user_id
//	LEFT JOIN (
//	    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
//	    FROM opinions
//	    GROUP BY parent_opinion_id
//	) rc ON opinions.opinion_id = rc.parent_opinion_id
//	WHERE opinions.user_id = $1
//	ORDER BY
//	    CASE $4::text
//	        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
//	        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
//	        WHEN 'mostReplies' THEN reply_count
//	    END DESC
//	LIMIT $2 OFFSET $3
func (q *Queries) GetOpinionsByUserID(ctx context.Context, arg GetOpinionsByUserIDParams) ([]GetOpinionsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getOpinionsByUserID,
		arg.UserID,
		arg.Limit,
		arg.Offset,
		arg.SortKey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOpinionsByUserIDRow
	for rows.Next() {
		var i GetOpinionsByUserIDRow
		if err := rows.Scan(
			&i.OpinionID,
			&i.TalkSessionID,
			&i.UserID,
			&i.ParentOpinionID,
			&i.Title,
			&i.Content,
			&i.ReferenceUrl,
			&i.PictureUrl,
			&i.CreatedAt,
			&i.DisplayName,
			&i.DisplayID,
			&i.IconUrl,
			&i.ParentVoteType,
			&i.ReplyCount,
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

const getParentOpinions = `-- name: GetParentOpinions :many
WITH RECURSIVE opinion_tree AS (
    -- ベースケース：指定された意見から開始
    SELECT
        o.opinion_id,
        o.parent_opinion_id,
        1 as level
    FROM opinions o
    WHERE o.opinion_id = $1

    UNION ALL

    SELECT
        p.opinion_id,
        p.parent_opinion_id,
        t.level + 1
    FROM opinions p
    INNER JOIN opinion_tree t ON t.parent_opinion_id = p.opinion_id
)
SELECT
    o.opinion_id,
    o.talk_session_id,
    o.user_id,
    o.parent_opinion_id,
    o.title,
    o.content,
    o.reference_url,
    o.picture_url,
    o.created_at,
    u.display_name AS display_name,
    u.display_id AS display_id,
    u.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(rc.reply_count, 0) AS reply_count,
    COALESCE(cv.vote_type, 0) AS current_vote_type,
    ot.level
FROM opinion_tree ot
JOIN opinions o ON ot.opinion_id = o.opinion_id
LEFT JOIN users u ON o.user_id = u.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON o.parent_opinion_id = pv.opinion_id
    AND o.user_id = pv.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) cv ON opinions.user_id = $2::uuid
    AND opinions.opinion_id = cv.opinion_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON o.opinion_id = rc.parent_opinion_id
ORDER BY ot.level DESC
`

type GetParentOpinionsParams struct {
	OpinionID uuid.UUID
	UserID    uuid.NullUUID
}

type GetParentOpinionsRow struct {
	OpinionID       uuid.UUID
	TalkSessionID   uuid.UUID
	UserID          uuid.UUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         string
	ReferenceUrl    sql.NullString
	PictureUrl      sql.NullString
	CreatedAt       time.Time
	DisplayName     sql.NullString
	DisplayID       sql.NullString
	IconUrl         sql.NullString
	ParentVoteType  int16
	ReplyCount      int64
	CurrentVoteType int16
	Level           int32
}

// GetParentOpinions
//
//	WITH RECURSIVE opinion_tree AS (
//	    -- ベースケース：指定された意見から開始
//	    SELECT
//	        o.opinion_id,
//	        o.parent_opinion_id,
//	        1 as level
//	    FROM opinions o
//	    WHERE o.opinion_id = $1
//
//	    UNION ALL
//
//	    SELECT
//	        p.opinion_id,
//	        p.parent_opinion_id,
//	        t.level + 1
//	    FROM opinions p
//	    INNER JOIN opinion_tree t ON t.parent_opinion_id = p.opinion_id
//	)
//	SELECT
//	    o.opinion_id,
//	    o.talk_session_id,
//	    o.user_id,
//	    o.parent_opinion_id,
//	    o.title,
//	    o.content,
//	    o.reference_url,
//	    o.picture_url,
//	    o.created_at,
//	    u.display_name AS display_name,
//	    u.display_id AS display_id,
//	    u.icon_url AS icon_url,
//	    COALESCE(pv.vote_type, 0) AS parent_vote_type,
//	    COALESCE(rc.reply_count, 0) AS reply_count,
//	    COALESCE(cv.vote_type, 0) AS current_vote_type,
//	    ot.level
//	FROM opinion_tree ot
//	JOIN opinions o ON ot.opinion_id = o.opinion_id
//	LEFT JOIN users u ON o.user_id = u.user_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.user_id, votes.opinion_id
//	    FROM votes
//	) pv ON o.parent_opinion_id = pv.opinion_id
//	    AND o.user_id = pv.user_id
//	LEFT JOIN (
//	    SELECT votes.vote_type, votes.user_id, votes.opinion_id
//	    FROM votes
//	) cv ON opinions.user_id = $2::uuid
//	    AND opinions.opinion_id = cv.opinion_id
//	LEFT JOIN (
//	    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
//	    FROM opinions
//	    GROUP BY parent_opinion_id
//	) rc ON o.opinion_id = rc.parent_opinion_id
//	ORDER BY ot.level DESC
func (q *Queries) GetParentOpinions(ctx context.Context, arg GetParentOpinionsParams) ([]GetParentOpinionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getParentOpinions, arg.OpinionID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetParentOpinionsRow
	for rows.Next() {
		var i GetParentOpinionsRow
		if err := rows.Scan(
			&i.OpinionID,
			&i.TalkSessionID,
			&i.UserID,
			&i.ParentOpinionID,
			&i.Title,
			&i.Content,
			&i.ReferenceUrl,
			&i.PictureUrl,
			&i.CreatedAt,
			&i.DisplayName,
			&i.DisplayID,
			&i.IconUrl,
			&i.ParentVoteType,
			&i.ReplyCount,
			&i.CurrentVoteType,
			&i.Level,
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

const getRandomOpinions = `-- name: GetRandomOpinions :many
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.reference_url,
    opinions.picture_url,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(rc.reply_count, 0) AS reply_count
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT opinions.opinion_id
    FROM opinions
    LEFT JOIN votes
        ON opinions.opinion_id = votes.opinion_id
        AND votes.user_id = $1
    GROUP BY opinions.opinion_id
    HAVING COUNT(votes.vote_id) = 0
) vote_count ON opinions.opinion_id = vote_count.opinion_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON rc.opinion_id = opinions.opinion_id
LEFT JOIN (
    SELECT rank, opinion_id
    FROM representative_opinions
) ro ON opinions.opinion_id = ro.opinion_id
WHERE opinions.talk_session_id = $2
    AND vote_count.opinion_id = opinions.opinion_id
    AND opinions.parent_opinion_id IS NULL
ORDER BY
    CASE $4::text
        WHEN 'top' THEN COALESCE(ro.rank, 0)
        ELSE RANDOM()
    END ASC
LIMIT $3
`

type GetRandomOpinionsParams struct {
	UserID        uuid.UUID
	TalkSessionID uuid.UUID
	Limit         int32
	SortKey       sql.NullString
}

type GetRandomOpinionsRow struct {
	OpinionID       uuid.UUID
	TalkSessionID   uuid.UUID
	UserID          uuid.UUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         string
	ReferenceUrl    sql.NullString
	PictureUrl      sql.NullString
	CreatedAt       time.Time
	DisplayName     sql.NullString
	DisplayID       sql.NullString
	IconUrl         sql.NullString
	ReplyCount      int64
}

// 指定されたユーザーが投票していない意見のみを取得
// この意見に対するリプライ数
// グループ内のランクを取得
// トークセッションに紐づく意見のみを取得
//
//	SELECT
//	    opinions.opinion_id,
//	    opinions.talk_session_id,
//	    opinions.user_id,
//	    opinions.parent_opinion_id,
//	    opinions.title,
//	    opinions.content,
//	    opinions.reference_url,
//	    opinions.picture_url,
//	    opinions.created_at,
//	    users.display_name AS display_name,
//	    users.display_id AS display_id,
//	    users.icon_url AS icon_url,
//	    COALESCE(rc.reply_count, 0) AS reply_count
//	FROM opinions
//	LEFT JOIN users
//	    ON opinions.user_id = users.user_id
//	LEFT JOIN (
//	    SELECT opinions.opinion_id
//	    FROM opinions
//	    LEFT JOIN votes
//	        ON opinions.opinion_id = votes.opinion_id
//	        AND votes.user_id = $1
//	    GROUP BY opinions.opinion_id
//	    HAVING COUNT(votes.vote_id) = 0
//	) vote_count ON opinions.opinion_id = vote_count.opinion_id
//	LEFT JOIN (
//	    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
//	    FROM opinions
//	    GROUP BY parent_opinion_id
//	) rc ON rc.opinion_id = opinions.opinion_id
//	LEFT JOIN (
//	    SELECT rank, opinion_id
//	    FROM representative_opinions
//	) ro ON opinions.opinion_id = ro.opinion_id
//	WHERE opinions.talk_session_id = $2
//	    AND vote_count.opinion_id = opinions.opinion_id
//	    AND opinions.parent_opinion_id IS NULL
//	ORDER BY
//	    CASE $4::text
//	        WHEN 'top' THEN COALESCE(ro.rank, 0)
//	        ELSE RANDOM()
//	    END ASC
//	LIMIT $3
func (q *Queries) GetRandomOpinions(ctx context.Context, arg GetRandomOpinionsParams) ([]GetRandomOpinionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRandomOpinions,
		arg.UserID,
		arg.TalkSessionID,
		arg.Limit,
		arg.SortKey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRandomOpinionsRow
	for rows.Next() {
		var i GetRandomOpinionsRow
		if err := rows.Scan(
			&i.OpinionID,
			&i.TalkSessionID,
			&i.UserID,
			&i.ParentOpinionID,
			&i.Title,
			&i.Content,
			&i.ReferenceUrl,
			&i.PictureUrl,
			&i.CreatedAt,
			&i.DisplayName,
			&i.DisplayID,
			&i.IconUrl,
			&i.ReplyCount,
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
