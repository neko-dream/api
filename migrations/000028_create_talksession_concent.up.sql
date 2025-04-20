CREATE TABLE IF NOT EXISTS talksession_consents (
  talksession_id UUID NOT NULL,
  user_id UUID NOT NULL,
  -- 同意した項目
  restrictions JSONB,
  consented_at TIMESTAMP NOT NULL,
  PRIMARY KEY (talksession_id, user_id)
);
