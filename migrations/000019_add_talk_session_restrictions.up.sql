ALTER TABLE talk_sessions ADD COLUMN restrictions JSONB;
ALTER TABLE talk_sessions ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
