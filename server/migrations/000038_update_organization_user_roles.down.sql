-- Revert organization user roles to old values
-- SuperAdmin: 10 -> 4
-- Owner: 20 -> 3
-- Admin: 30 -> 2
-- Member: 40 -> 1

-- First, create a temporary column to avoid conflicts during update
ALTER TABLE organization_users ADD COLUMN role_old INT;

-- Update roles to old values
UPDATE organization_users SET role_old = 
    CASE 
        WHEN role = 10 THEN 4  -- SuperAdmin
        WHEN role = 20 THEN 3  -- Owner
        WHEN role = 30 THEN 2  -- Admin
        WHEN role = 40 THEN 1  -- Member
        ELSE role
    END;

-- Drop new column and rename old column
ALTER TABLE organization_users DROP COLUMN role;
ALTER TABLE organization_users RENAME COLUMN role_old TO role;

-- Add NOT NULL constraint
ALTER TABLE organization_users ALTER COLUMN role SET NOT NULL;

-- Update default value to 0
ALTER TABLE organization_users ALTER COLUMN role SET DEFAULT 0;