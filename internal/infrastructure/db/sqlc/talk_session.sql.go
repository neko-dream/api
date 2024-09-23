// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: talk_session.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createTalkSession = `-- name: CreateTalkSession :exec
INSERT INTO talk_sessions (talk_session_id, theme, created_at) VALUES ($1, $2, $3)
`

type CreateTalkSessionParams struct {
	TalkSessionID uuid.UUID
	Theme         string
	CreatedAt     time.Time
}

func (q *Queries) CreateTalkSession(ctx context.Context, arg CreateTalkSessionParams) error {
	_, err := q.db.ExecContext(ctx, createTalkSession, arg.TalkSessionID, arg.Theme, arg.CreatedAt)
	return err
}

const editTalkSession = `-- name: EditTalkSession :exec
UPDATE talk_sessions SET theme = $2, finished_at = $3 WHERE talk_session_id = $1
`

type EditTalkSessionParams struct {
	TalkSessionID uuid.UUID
	Theme         string
	FinishedAt    sql.NullTime
}

func (q *Queries) EditTalkSession(ctx context.Context, arg EditTalkSessionParams) error {
	_, err := q.db.ExecContext(ctx, editTalkSession, arg.TalkSessionID, arg.Theme, arg.FinishedAt)
	return err
}

const getTalkSessionByID = `-- name: GetTalkSessionByID :one
SELECT talk_session_id, theme, finished_at, created_at FROM talk_sessions WHERE talk_session_id = $1
`

func (q *Queries) GetTalkSessionByID(ctx context.Context, talkSessionID uuid.UUID) (TalkSession, error) {
	row := q.db.QueryRowContext(ctx, getTalkSessionByID, talkSessionID)
	var i TalkSession
	err := row.Scan(
		&i.TalkSessionID,
		&i.Theme,
		&i.FinishedAt,
		&i.CreatedAt,
	)
	return i, err
}