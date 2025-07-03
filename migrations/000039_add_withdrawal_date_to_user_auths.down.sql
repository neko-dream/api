-- Remove withdrawal_date from user_auths table
DROP INDEX IF EXISTS idx_user_auths_withdrawal_date;
ALTER TABLE user_auths DROP COLUMN IF EXISTS withdrawal_date;