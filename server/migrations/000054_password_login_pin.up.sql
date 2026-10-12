ALTER TABLE password_login_challenges
    ADD COLUMN pin_hash BYTEA CHECK (pin_hash IS NULL OR octet_length(pin_hash) = 32),
    ADD COLUMN pin_attempts INTEGER NOT NULL DEFAULT 0 CHECK (pin_attempts BETWEEN 0 AND 5);
