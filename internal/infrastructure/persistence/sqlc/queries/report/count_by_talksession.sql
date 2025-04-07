-- name: CountReportsByTalkSession :one
SELECT
  COUNT(*) AS count
FROM
  opinion_reports
WHERE
  talk_session_id = sqlc.narg('talk_session_id')::uuid
  AND
  CASE
    WHEN sqlc.narg('status')::text IS NOT NULL THEN opinion_reports.status = sqlc.narg('status')::text
    ELSE TRUE
  END
;
