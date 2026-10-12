-- 管理画面用の管理者アカウント許可リスト。
-- ログインは Google 認証。email(Googleアカウント)で照合する。既存の users/組織とは無関係。
CREATE TABLE IF NOT EXISTS admin_users (
    admin_user_id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admin_users_email ON admin_users(email);
