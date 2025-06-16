ALTER TABLE talk_sessions
DROP COLUMN IF EXISTS organization_alias_id,
DROP COLUMN IF EXISTS organization_id;

DROP TABLE IF EXISTS organization_aliases;
