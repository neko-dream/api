-- admin のリフレッシュトークンを保持する(revocable)セッション。
-- token_hash は HMAC(refresh_token, pepper)。生トークンは保存しない。
CREATE TABLE IF NOT EXISTS admin_sessions (
    admin_session_id UUID PRIMARY KEY,
    admin_user_id UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admin_sessions_token_hash ON admin_sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_admin_sessions_admin_user_id ON admin_sessions(admin_user_id);
