-- name: FindOrganizationByID :one
SELECT
    sqlc.embed(organizations)
FROM organizations
WHERE organization_id = $1;
