-- name: FindOrganizationsByIDs :many
SELECT
    sqlc.embed(organizations)
FROM organizations
WHERE organization_id = ANY($1::uuid[]);
