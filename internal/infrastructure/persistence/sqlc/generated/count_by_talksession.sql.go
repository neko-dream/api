// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: count_by_talksession.sql

package model

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const countReportsByTalkSession = `-- name: CountReportsByTalkSession :one
SELECT
  COUNT(*) AS count
FROM
  opinion_reports
WHERE
  talk_session_id = $1::uuid
  AND
  CASE
    WHEN $2::text IS NOT NULL THEN opinion_reports.status = $2::text
    ELSE TRUE
  END
`

type CountReportsByTalkSessionParams struct {
	TalkSessionID uuid.NullUUID
	Status        sql.NullString
}

// CountReportsByTalkSession
//
//	SELECT
//	  COUNT(*) AS count
//	FROM
//	  opinion_reports
//	WHERE
//	  talk_session_id = $1::uuid
//	  AND
//	  CASE
//	    WHEN $2::text IS NOT NULL THEN opinion_reports.status = $2::text
//	    ELSE TRUE
//	  END
func (q *Queries) CountReportsByTalkSession(ctx context.Context, arg CountReportsByTalkSessionParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countReportsByTalkSession, arg.TalkSessionID, arg.Status)
	var count int64
	err := row.Scan(&count)
	return count, err
}
