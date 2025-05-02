-- name: FindReportByOpinionIDs :many
SELECT
    sqlc.embed(opinion_reports),
    sqlc.embed(opinions)
FROM
    opinion_reports
LEFT JOIN opinions
    ON opinion_reports.opinion_id = opinions.opinion_id
WHERE
    opinion_reports.opinion_id = ANY(sqlc.arg('opinion_ids')::uuid[])
    AND opinion_reports.status = sqlc.arg('status')
ORDER BY opinion_reports.created_at DESC
;
