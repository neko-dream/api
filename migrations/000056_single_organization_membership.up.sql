-- Never silently choose/delete existing memberships: duplicates block migration.
DO $$
BEGIN
  IF EXISTS (SELECT user_id FROM organization_users GROUP BY user_id HAVING COUNT(*) > 1) THEN
    RAISE EXCEPTION 'Cannot enforce one organization per user: organization_users contains duplicate user_id memberships; resolve explicitly before migrating';
  END IF;
END $$;
ALTER TABLE organization_users ADD CONSTRAINT organization_users_one_per_user UNIQUE (user_id);
