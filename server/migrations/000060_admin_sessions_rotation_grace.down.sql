DROP INDEX IF EXISTS idx_admin_sessions_previous_token_hash;
ALTER TABLE admin_sessions DROP COLUMN IF EXISTS rotated_at;
ALTER TABLE admin_sessions DROP COLUMN IF EXISTS previous_token_hash;
