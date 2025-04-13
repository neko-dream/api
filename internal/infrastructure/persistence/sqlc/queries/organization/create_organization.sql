-- name: CreateOrganization :exec
INSERT INTO organizations (
    organization_id,
    organization_type,
    name,
    owner_id
) VALUES ($1, $2, $3, $4)
