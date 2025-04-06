-- name: FindReportByOpinionID :many
SELECT
    sqlc.embed(opinion_reports)
FROM
    opinion_reports
WHERE
    opinion_id = sqlc.narg('opinion_id')::uuid
;
