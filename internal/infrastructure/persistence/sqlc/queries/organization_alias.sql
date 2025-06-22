-- name: CreateOrganizationAlias :one
INSERT INTO organization_aliases (
    alias_id,
    organization_id,
    alias_name,
    created_by
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetActiveOrganizationAliases :many
SELECT * FROM organization_aliases
WHERE organization_id = $1 AND deactivated_at IS NULL
ORDER BY created_at ASC;

-- name: GetOrganizationAliasById :one
SELECT * FROM organization_aliases
WHERE alias_id = $1;

-- name: DeactivateOrganizationAlias :exec
UPDATE organization_aliases
SET deactivated_at = CURRENT_TIMESTAMP,
    deactivated_by = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE alias_id = $1 AND deactivated_at IS NULL;

-- name: CountActiveAliasesByOrganization :one
SELECT COUNT(*) FROM organization_aliases
WHERE organization_id = $1 AND deactivated_at IS NULL;

-- name: CheckAliasNameExists :one
SELECT EXISTS(
    SELECT 1 FROM organization_aliases
    WHERE organization_id = $1 AND alias_name = $2 AND deactivated_at IS NULL
);
