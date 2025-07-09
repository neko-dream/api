-- 通知履歴テーブルを削除
DROP TABLE IF EXISTS notification_history;

-- 通知設定テーブルを削除
DROP TABLE IF EXISTS notification_preferences;

-- devicesテーブルから追加したカラムを削除
ALTER TABLE devices 
DROP COLUMN IF EXISTS device_name,
DROP COLUMN IF EXISTS app_version,
DROP COLUMN IF EXISTS os_version,
DROP COLUMN IF EXISTS last_active_at;

-- platformの制約を元に戻す
ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_platform_check;
ALTER TABLE devices ADD CONSTRAINT devices_platform_check 
    CHECK (platform IN ('APNS', 'GCM'));

-- 追加したインデックスを削除
DROP INDEX IF EXISTS idx_devices_token;