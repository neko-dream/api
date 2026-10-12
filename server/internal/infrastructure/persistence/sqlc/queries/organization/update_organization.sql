-- name: UpdateOrganization :exec
UPDATE organizations SET
    name = $2,
    icon_url = $3
WHERE organization_id = $1;

