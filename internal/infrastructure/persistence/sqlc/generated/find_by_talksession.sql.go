// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_by_talksession.sql

package model

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const findReportsByTalkSession = `-- name: FindReportsByTalkSession :many
SELECT
  opinion_reports.opinion_report_id, opinion_reports.opinion_id, opinion_reports.talk_session_id, opinion_reports.reporter_id, opinion_reports.reason, opinion_reports.status, opinion_reports.created_at, opinion_reports.updated_at, opinion_reports.reason_text
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

type FindReportsByTalkSessionParams struct {
	TalkSessionID uuid.NullUUID
	Status        sql.NullString
}

type FindReportsByTalkSessionRow struct {
	OpinionReport OpinionReport
}

// FindReportsByTalkSession
//
//	SELECT
//	  opinion_reports.opinion_report_id, opinion_reports.opinion_id, opinion_reports.talk_session_id, opinion_reports.reporter_id, opinion_reports.reason, opinion_reports.status, opinion_reports.created_at, opinion_reports.updated_at, opinion_reports.reason_text
//	FROM
//	  opinion_reports
//	WHERE
//	  talk_session_id = $1::uuid
//	  AND
//	  CASE
//	    WHEN $2::text IS NOT NULL THEN opinion_reports.status = $2::text
//	    ELSE TRUE
//	  END
func (q *Queries) FindReportsByTalkSession(ctx context.Context, arg FindReportsByTalkSessionParams) ([]FindReportsByTalkSessionRow, error) {
	rows, err := q.db.QueryContext(ctx, findReportsByTalkSession, arg.TalkSessionID, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindReportsByTalkSessionRow
	for rows.Next() {
		var i FindReportsByTalkSessionRow
		if err := rows.Scan(
			&i.OpinionReport.OpinionReportID,
			&i.OpinionReport.OpinionID,
			&i.OpinionReport.TalkSessionID,
			&i.OpinionReport.ReporterID,
			&i.OpinionReport.Reason,
			&i.OpinionReport.Status,
			&i.OpinionReport.CreatedAt,
			&i.OpinionReport.UpdatedAt,
			&i.OpinionReport.ReasonText,
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
