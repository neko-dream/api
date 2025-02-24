-- votes に talk_session_id カラムを追加
ALTER TABLE votes ADD COLUMN talk_session_id uuid NOT NULL;
