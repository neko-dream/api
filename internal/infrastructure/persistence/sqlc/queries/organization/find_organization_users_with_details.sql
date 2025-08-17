-- name: FindOrganizationUsersWithDetails :many
SELECT
    sqlc.embed(ou),
    sqlc.embed(u),
    sqlc.embed(o)
FROM organizations o
LEFT JOIN organization_users ou ON o.organization_id = ou.organization_id
LEFT JOIN users u ON ou.user_id = u.user_id
WHERE o.organization_id = $1
-- 退会ユーザーは表示しない
AND u.withdrawal_date IS NULL
ORDER BY ou.role ASC, u.user_id ASC;
