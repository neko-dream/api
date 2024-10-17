// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: analysis.sql

package model

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getGroupInfoByTalkSessionId = `-- name: GetGroupInfoByTalkSessionId :many
SELECT
    user_group_info.pos_x,
    user_group_info.pos_y,
    user_group_info.group_id,
    user_group_info.perimeter_index,
    users.display_id AS display_id,
    user_group_info.user_id
FROM user_group_info
LEFT JOIN users
    ON user_group_info.user_id = users.user_id
WHERE talk_session_id = $1
`

type GetGroupInfoByTalkSessionIdRow struct {
	PosX           float64
	PosY           float64
	GroupID        int32
	PerimeterIndex sql.NullInt32
	DisplayID      sql.NullString
	UserID         uuid.UUID
}

func (q *Queries) GetGroupInfoByTalkSessionId(ctx context.Context, talkSessionID uuid.UUID) ([]GetGroupInfoByTalkSessionIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getGroupInfoByTalkSessionId, talkSessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGroupInfoByTalkSessionIdRow
	for rows.Next() {
		var i GetGroupInfoByTalkSessionIdRow
		if err := rows.Scan(
			&i.PosX,
			&i.PosY,
			&i.GroupID,
			&i.PerimeterIndex,
			&i.DisplayID,
			&i.UserID,
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

const getGroupListByTalkSessionId = `-- name: GetGroupListByTalkSessionId :many
SELECT
    DISTINCT user_group_info.group_id
FROM user_group_info
WHERE talk_session_id = $1
`

func (q *Queries) GetGroupListByTalkSessionId(ctx context.Context, talkSessionID uuid.UUID) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getGroupListByTalkSessionId, talkSessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var group_id int32
		if err := rows.Scan(&group_id); err != nil {
			return nil, err
		}
		items = append(items, group_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRepresentativeOpinionsByTalkSessionId = `-- name: GetRepresentativeOpinionsByTalkSessionId :many
SELECT
    representative_opinions.group_id,
    representative_opinions.rank,
    opinions.opinion_id,
    opinions.talk_session_id,
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
FROM representative_opinions
LEFT JOIN opinions
    ON representative_opinions.opinion_id = opinions.opinion_id
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
WHERE representative_opinions.rank < 2
    AND opinions.talk_session_id = $1
`

type GetRepresentativeOpinionsByTalkSessionIdRow struct {
	GroupID         int32
	Rank            int32
	OpinionID       uuid.NullUUID
	TalkSessionID   uuid.NullUUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         sql.NullString
	ReferenceUrl    sql.NullString
	PictureUrl      sql.NullString
	CreatedAt       sql.NullTime
	DisplayName     sql.NullString
	DisplayID       sql.NullString
	IconUrl         sql.NullString
	ReplyCount      int64
}

func (q *Queries) GetRepresentativeOpinionsByTalkSessionId(ctx context.Context, talkSessionID uuid.UUID) ([]GetRepresentativeOpinionsByTalkSessionIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getRepresentativeOpinionsByTalkSessionId, talkSessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepresentativeOpinionsByTalkSessionIdRow
	for rows.Next() {
		var i GetRepresentativeOpinionsByTalkSessionIdRow
		if err := rows.Scan(
			&i.GroupID,
			&i.Rank,
			&i.OpinionID,
			&i.TalkSessionID,
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
