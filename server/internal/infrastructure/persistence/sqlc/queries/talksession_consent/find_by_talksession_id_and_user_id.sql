-- name: FindTSConsentByTalksessionIdAndUserId :one
SELECT
    sqlc.embed(talksession_consents)
FROM
    talksession_consents
WHERE
    talksession_id = $1
    AND user_id = $2;
