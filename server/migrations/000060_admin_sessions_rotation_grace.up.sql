-- リフレッシュトークンのローテーション直後に、直前のトークンを短時間だけ受け付けるための列。
-- 複数タブが同時にリフレッシュしたとき、後から来た側が失敗してログイン画面に戻されるのを防ぐ。
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS previous_token_hash VARCHAR(255);
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS rotated_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_admin_sessions_previous_token_hash ON admin_sessions(previous_token_hash);
