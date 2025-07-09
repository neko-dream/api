-- 既存のdevicesテーブルにカラムを追加
ALTER TABLE devices
ADD COLUMN IF NOT EXISTS device_name TEXT,
ADD COLUMN IF NOT EXISTS app_version TEXT,
ADD COLUMN IF NOT EXISTS os_version TEXT,
ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMP;

-- platformの値を更新
ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_platform_check;
ALTER TABLE devices ADD CONSTRAINT devices_platform_check
    CHECK (platform IN ('APNS', 'GCM', 'ios', 'android', 'web'));

-- device_tokenにユニーク制約を追加
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_token ON devices(device_token);

-- 通知設定テーブルを作成
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    push_notification_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

CREATE TABLE IF NOT EXISTS notification_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    notification_type TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    data JSONB,
    status TEXT NOT NULL CHECK (status IN ('sent', 'failed', 'pending')),
    failure_reason TEXT,
    read BOOLEAN NOT NULL DEFAULT false,
    read_at TIMESTAMP,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notification_history_user_id ON notification_history(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_history_sent_at ON notification_history(sent_at DESC);
