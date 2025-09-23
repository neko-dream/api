ALTER TABLE talk_sessions ADD COLUMN hide_top BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_talk_sessions_hide_top ON talk_sessions(hide_top) WHERE hide_top = FALSE;
