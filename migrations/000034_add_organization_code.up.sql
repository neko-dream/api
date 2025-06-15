ALTER TABLE organizations ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
CREATE INDEX idx_organizations_code ON organizations(code);
ALTER TABLE auth_states ADD COLUMN organization_id UUID;
