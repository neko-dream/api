-- name: FindOrgUserByOrganizationID :many
SELECT
    sqlc.embed(organization_users)
FROM organization_users
WHERE organization_id = $1;
