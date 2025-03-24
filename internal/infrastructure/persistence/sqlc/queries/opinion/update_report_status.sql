-- name: UpdateReportStatus :exec
UPDATE opinion_reports
SET status = $1
WHERE opinion_report_id = $2
RETURNING *;

