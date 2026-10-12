ALTER TABLE users ADD COLUMN can_create_talk_session boolean NOT NULL DEFAULT false;

-- Preserve existing creation access once. Future membership changes do not grant access.
UPDATE users u SET can_create_talk_session = true
WHERE u.withdrawal_date IS NULL
  AND EXISTS (SELECT 1 FROM organization_users ou WHERE ou.user_id = u.user_id)
  AND EXISTS (SELECT 1 FROM user_auths ua WHERE ua.user_id = u.user_id AND ua.provider <> 'GUEST');
