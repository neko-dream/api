CREATE TABLE policy_consents (
  policy_consent_id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  policy_version VARCHAR(20) NOT NULL,
  consented_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  ip_address VARCHAR(45) NOT NULL,
  user_agent TEXT NOT NULL
);

CREATE INDEX idx_user_policy ON policy_consents (user_id, policy_version);

CREATE TABLE policy_versions (
  version VARCHAR(20) PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_version_created_at ON policy_versions (version, created_at);
