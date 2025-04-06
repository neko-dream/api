-- name: FindReportByOpinionID :many
SELECT
    sqlc.embed(opinion_reports)
FROM
    opinion_reports
WHERE
    opinion_id = sqlc.narg('opinion_id')::uuid
    AND
    CASE
        WHEN sqlc.narg('status')::text IS NOT NULL THEN opinion_reports.status = sqlc.narg('status')::text
        ELSE TRUE
    END
;
