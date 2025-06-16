-- Remove unique index for primary alias per organization
DROP INDEX IF EXISTS unique_primary_alias_per_org;

-- Remove is_primary column from organization_aliases table
ALTER TABLE organization_aliases DROP COLUMN IF EXISTS is_primary;