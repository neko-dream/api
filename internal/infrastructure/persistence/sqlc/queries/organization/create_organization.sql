-- name: CreateOrganization :exec
INSERT INTO organizations (
    organization_id,
    organization_type,
    name,
    owner_id,
    code,
    icon_url
) VALUES ($1, $2, $3, $4, $5, $6);
