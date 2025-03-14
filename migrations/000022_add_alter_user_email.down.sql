-- userからemailを削除
ALTER TABLE user DROP COLUMN IF EXISTS email;
ALTER TABLE user DROP COLUMN IF EXISTS email_verified;
