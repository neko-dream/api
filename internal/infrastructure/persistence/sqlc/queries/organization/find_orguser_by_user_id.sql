-- name: FindOrgUserByUserID :many
SELECT
    sqlc.embed(organization_users)
FROM organization_users
WHERE organization_users.user_id = $1;

-- name: FindOrgUserByUserIDWithOrganization :many
SELECT
    sqlc.embed(organization_users),
    sqlc.embed(organizations)
FROM organization_users
LEFT JOIN organizations ON organization_users.organization_id = organizations.organization_id
WHERE organization_users.user_id = $1;

