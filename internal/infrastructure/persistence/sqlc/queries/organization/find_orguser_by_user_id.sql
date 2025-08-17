-- name: FindOrgUserByUserID :many
SELECT
    sqlc.embed(organization_users)
FROM organization_users
WHERE organization_users.user_id = $1;

-- name: FindOrgUserByUserIDWithOrganization :many
SELECT
    sqlc.embed(users),
    sqlc.embed(organization_users),
    sqlc.embed(organizations)
FROM organization_users
LEFT JOIN organizations ON organization_users.organization_id = organizations.organization_id
LEFT JOIN users ON organization_users.user_id = users.user_id
WHERE organization_users.user_id = $1
    AND users.withdrawal_date IS NULL;

