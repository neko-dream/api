-- Update organization user roles to new values
-- SuperAdmin: 4 -> 10
-- Owner: 3 -> 20
-- Admin: 2 -> 30
-- Member: 1 -> 40

-- First, create a temporary column to avoid conflicts during update
ALTER TABLE organization_users ADD COLUMN role_new INT;

-- Update roles to new values
UPDATE organization_users SET role_new =
    CASE
        WHEN role = 4 THEN 10  -- SuperAdmin
        WHEN role = 3 THEN 20  -- Owner
        WHEN role = 2 THEN 30  -- Admin
        WHEN role = 1 THEN 40  -- Member
        WHEN role = 0 THEN 40  -- Default (Member)
        ELSE role
    END;

-- Drop old column and rename new column
ALTER TABLE organization_users DROP COLUMN role;
ALTER TABLE organization_users RENAME COLUMN role_new TO role;

-- Add NOT NULL constraint
ALTER TABLE organization_users ALTER COLUMN role SET NOT NULL;

-- Update default value to 40 (Member)
ALTER TABLE organization_users ALTER COLUMN role SET DEFAULT 40;
