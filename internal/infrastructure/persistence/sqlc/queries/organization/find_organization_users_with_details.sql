-- name: FindOrganizationUsersWithDetails :many
SELECT
    ou.user_id,
    ou.role,
    u.display_id,
    u.display_name,
    u.icon_url
FROM organization_users ou
INNER JOIN users u ON ou.user_id = u.user_id
WHERE ou.organization_id = $1
ORDER BY ou.role ASC, u.display_name ASC;