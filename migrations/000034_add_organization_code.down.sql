ALTER TABLE auth_states DROP COLUMN organization_id;
DROP INDEX IF EXISTS idx_organizations_code;
ALTER TABLE organizations DROP COLUMN code;
