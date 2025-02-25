-- talksessionsテーブルにthumbnail_urlカラムを削除
ALTER TABLE talk_sessions DROP COLUMN  IF EXISTS thumbnail_url CASCADE;
