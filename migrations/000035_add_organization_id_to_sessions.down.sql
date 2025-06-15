-- Remove organization_id from sessions table
DROP INDEX IF EXISTS idx_sessions_organization_id;
ALTER TABLE sessions DROP COLUMN organization_id;