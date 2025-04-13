-- name: FindOrgUserByUserID :one
SELECT
    sqlc.embed(organization_users)
FROM organization_users
WHERE user_id = $1;
