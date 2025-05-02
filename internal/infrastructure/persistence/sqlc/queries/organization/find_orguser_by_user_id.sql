-- name: FindOrgUserByUserID :many
SELECT
    sqlc.embed(organization_users)
FROM organization_users
WHERE organization_users.user_id = $1;
