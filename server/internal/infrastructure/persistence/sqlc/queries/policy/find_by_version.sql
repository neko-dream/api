-- name: FindPolicyByVersion :one
SELECT
    sqlc.embed(policy_versions)
FROM
    policy_versions
WHERE
    version = $1
LIMIT 1;
