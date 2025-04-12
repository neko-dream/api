-- name: FindReportByOpinionIDs :many
SELECT
    sqlc.embed(opinion_reports),
    sqlc.embed(opinions)
FROM
    opinion_reports
LEFT JOIN opinions
    ON opinion_reports.opinion_id = opinions.opinion_id
WHERE
    opinions.opinion_id = ANY($1::uuid[])
ORDER BY opinion_reports.created_at DESC
;
