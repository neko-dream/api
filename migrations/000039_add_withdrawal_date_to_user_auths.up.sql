-- Add withdrawal_date to user_auths table for tracking user withdrawals
ALTER TABLE user_auths ADD COLUMN withdrawal_date TIMESTAMP;

-- Add index for efficient querying of withdrawn users
CREATE INDEX idx_user_auths_withdrawal_date ON user_auths(withdrawal_date) WHERE withdrawal_date IS NOT NULL;