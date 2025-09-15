DROP INDEX IF EXISTS idx_talk_sessions_hide_top;

ALTER TABLE talk_sessions DROP COLUMN hide_top;
