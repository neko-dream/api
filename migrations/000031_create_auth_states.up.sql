CREATE TABLE auth_states (
    id SERIAL PRIMARY KEY,
    state VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    redirect_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(state)
);

CREATE INDEX idx_auth_states_state ON auth_states(state);
CREATE INDEX idx_auth_states_expires_at ON auth_states(expires_at);
