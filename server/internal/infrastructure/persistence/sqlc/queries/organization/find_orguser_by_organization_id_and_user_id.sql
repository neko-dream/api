-- name: FindOrgUserByOrganizationIDAndUserID :one
SELECT
    sqlc.embed(organization_users)
FROM organization_users
WHERE organization_id = $1
  AND user_id = $2;
