// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_by_opinions.sql

package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const findReportByOpinionIDs = `-- name: FindReportByOpinionIDs :many
SELECT
    opinion_reports.opinion_report_id, opinion_reports.opinion_id, opinion_reports.talk_session_id, opinion_reports.reporter_id, opinion_reports.reason, opinion_reports.status, opinion_reports.created_at, opinion_reports.updated_at, opinion_reports.reason_text,
    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.created_at, opinions.picture_url, opinions.reference_url
FROM
    opinion_reports
LEFT JOIN opinions
    ON opinion_reports.opinion_id = opinions.opinion_id
WHERE
    opinions.opinion_id = ANY($1::uuid[])
ORDER BY opinion_reports.created_at DESC
`

type FindReportByOpinionIDsRow struct {
	OpinionReport OpinionReport
	Opinion       Opinion
}

// FindReportByOpinionIDs
//
//	SELECT
//	    opinion_reports.opinion_report_id, opinion_reports.opinion_id, opinion_reports.talk_session_id, opinion_reports.reporter_id, opinion_reports.reason, opinion_reports.status, opinion_reports.created_at, opinion_reports.updated_at, opinion_reports.reason_text,
//	    opinions.opinion_id, opinions.talk_session_id, opinions.user_id, opinions.parent_opinion_id, opinions.title, opinions.content, opinions.created_at, opinions.picture_url, opinions.reference_url
//	FROM
//	    opinion_reports
//	LEFT JOIN opinions
//	    ON opinion_reports.opinion_id = opinions.opinion_id
//	WHERE
//	    opinions.opinion_id = ANY($1::uuid[])
//	ORDER BY opinion_reports.created_at DESC
func (q *Queries) FindReportByOpinionIDs(ctx context.Context, dollar_1 []uuid.UUID) ([]FindReportByOpinionIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, findReportByOpinionIDs, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindReportByOpinionIDsRow
	for rows.Next() {
		var i FindReportByOpinionIDsRow
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
			&i.Opinion.OpinionID,
			&i.Opinion.TalkSessionID,
			&i.Opinion.UserID,
			&i.Opinion.ParentOpinionID,
			&i.Opinion.Title,
			&i.Opinion.Content,
			&i.Opinion.CreatedAt,
			&i.Opinion.PictureUrl,
			&i.Opinion.ReferenceUrl,
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
