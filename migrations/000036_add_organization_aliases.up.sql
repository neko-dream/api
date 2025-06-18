CREATE TABLE organization_aliases (
    alias_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(organization_id) ON DELETE CASCADE,
    alias_name VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL REFERENCES users(user_id),
    deactivated_at TIMESTAMP,
    deactivated_by UUID REFERENCES users(user_id),
    CONSTRAINT unique_org_alias_name UNIQUE (organization_id, alias_name)
);

CREATE INDEX idx_organization_aliases_org_id ON organization_aliases(organization_id);
CREATE INDEX idx_organization_aliases_active ON organization_aliases(organization_id, deactivated_at);

CREATE UNIQUE INDEX unique_primary_alias_per_org ON organization_aliases(organization_id)
WHERE is_primary = TRUE AND deactivated_at IS NULL;

ALTER TABLE talk_sessions
ADD COLUMN organization_id UUID REFERENCES organizations(organization_id),
ADD COLUMN organization_alias_id UUID REFERENCES organization_aliases(alias_id);

CREATE INDEX idx_talk_sessions_organization_id ON talk_sessions(organization_id);

INSERT INTO organization_aliases (organization_id, alias_name, is_primary, created_by)
SELECT
    o.organization_id,
    o.name,
    TRUE,
    o.owner_id
FROM organizations o
WHERE NOT EXISTS (
    SELECT 1 FROM organization_aliases oa
    WHERE oa.organization_id = o.organization_id AND oa.is_primary = TRUE
);
