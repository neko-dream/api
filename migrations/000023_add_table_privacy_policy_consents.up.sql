CREATE TABLE policy_consents (
  policy_consent_id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  policy_version VARCHAR(20) NOT NULL,
  consented_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  ip_address VARCHAR(45) NOT NULL,
  user_agent TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX idx_user_policy (user_id, policy_version)
);

CREATE TABLE policy_versions (
  version VARCHAR(20) PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  INDEX idx_version_created_at (version, created_at)
);
