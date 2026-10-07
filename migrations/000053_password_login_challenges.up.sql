CREATE TABLE password_login_challenges (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE CHECK (octet_length(token_hash) = 32),
    browser_hash BYTEA NOT NULL CHECK (octet_length(browser_hash) = 32),
    credential_hash BYTEA NOT NULL CHECK (octet_length(credential_hash) = 32),
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    CHECK (expires_at > created_at)
);
