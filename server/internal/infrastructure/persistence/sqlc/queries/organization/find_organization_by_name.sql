-- name: FindOrganizationByName :one
SELECT
    sqlc.embed(organizations)
FROM organizations
WHERE name = $1;
