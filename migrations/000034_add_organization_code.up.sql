ALTER TABLE organizations ADD COLUMN code VARCHAR(50);

UPDATE organizations
SET code = CONCAT('ORG-', LPAD(CAST(row_num AS TEXT), 4, '0'))
FROM (
    SELECT organization_id, ROW_NUMBER() OVER (ORDER BY organization_id) as row_num
    FROM organizations
) ranked
WHERE organizations.organization_id = ranked.organization_id;

ALTER TABLE organizations ALTER COLUMN code SET NOT NULL;
CREATE INDEX idx_organizations_code ON organizations(code);
ALTER TABLE auth_states ADD COLUMN organization_id UUID;
