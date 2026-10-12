-- name: GetLatestPolicyVersion :one
SELECT
    sqlc.embed(policy_versions)
FROM
    policy_versions
ORDER BY created_at DESC
LIMIT 1;

