-- ユーザー退会機能の追加

-- usersテーブルにwithdrawal_dateカラムを追加
ALTER TABLE users ADD COLUMN withdrawal_date TIMESTAMP WITH TIME ZONE;

-- インデックスの追加（退会ユーザーの検索用）
CREATE INDEX idx_users_withdrawal_date ON users(withdrawal_date) WHERE withdrawal_date IS NOT NULL;

-- ユーザーステータス変更ログテーブルの作成
CREATE TABLE user_status_change_logs (
    user_status_change_logs_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'withdrawn', 'reactivated')),
    reason TEXT,
    changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    changed_by VARCHAR(20) NOT NULL CHECK (changed_by IN ('user', 'admin', 'system')),
    ip_address INET,
    user_agent TEXT,
    additional_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_status_logs_user_id ON user_status_change_logs(user_id);
CREATE INDEX idx_user_status_logs_changed_at ON user_status_change_logs(changed_at);
CREATE INDEX idx_user_status_logs_status ON user_status_change_logs(status);
COMMENT ON COLUMN users.withdrawal_date IS 'ユーザーの退会日時。NULLの場合はアクティブユーザー';
COMMENT ON TABLE user_status_change_logs IS 'ユーザーのステータス変更履歴（退会・復活など）';
