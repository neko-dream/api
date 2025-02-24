-- talk_sessions テーブルの finished_at カラムを削除
ALTER TABLE talk_sessions DROP COLUMN IF EXISTS finished_at CASCADE;
