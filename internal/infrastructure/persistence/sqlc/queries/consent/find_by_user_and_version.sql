-- name: FindConsentByUserAndVersion :one
SELECT
    sqlc.embed(policy_consents)
FROM
    policy_consents
WHERE
    user_id = $1
    AND policy_version = $2;
