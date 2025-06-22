-- name: FindOrganizationByCode :one
SELECT
    sqlc.embed(organizations)
FROM organizations
WHERE code = $1;