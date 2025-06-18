-- Add is_primary column back to organization_aliases table
ALTER TABLE organization_aliases ADD COLUMN IF NOT EXISTS is_primary BOOLEAN DEFAULT FALSE;

-- Recreate unique index for primary alias per organization
CREATE UNIQUE INDEX unique_primary_alias_per_org ON organization_aliases(organization_id)
WHERE is_primary = TRUE AND deactivated_at IS NULL;

-- Set the first active alias for each organization as primary
UPDATE organization_aliases oa
SET is_primary = TRUE
FROM (
    SELECT DISTINCT ON (organization_id) alias_id
    FROM organization_aliases
    WHERE deactivated_at IS NULL
    ORDER BY organization_id, created_at ASC
) first_alias
WHERE oa.alias_id = first_alias.alias_id;