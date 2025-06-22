-- Add organization_id to sessions table for organization-based login
ALTER TABLE sessions ADD COLUMN organization_id UUID;

-- Add index for performance
CREATE INDEX idx_sessions_organization_id ON sessions(organization_id);