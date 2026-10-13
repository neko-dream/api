CREATE TABLE password_setup_tokens (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE CHECK (octet_length(token_hash) = 32),
    email_subject TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    CHECK (expires_at > created_at)
);
