CREATE TABLE privacy_policy_consents (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id UUID NOT NULL,
  policy_version VARCHAR(20) NOT NULL,
  consented_at DATETIME NOT NULL,
  ip_address VARCHAR(45),
  user_agent TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX idx_user_policy (user_id, policy_version)
);
