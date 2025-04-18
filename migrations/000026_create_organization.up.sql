CREATE TABLE IF NOT EXISTS organizations (
    organization_id UUID PRIMARY KEY,
    organization_type INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS organization_users (
    organization_user_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    role INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS password_auth (
    password_auth_id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(255),
    required_password_change BOOLEAN NOT NULL DEFAULT TRUE,
    last_changed TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_password_auth_user_id ON password_auth(user_id);
