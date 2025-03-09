
-- talk_sessionsからhousehold_sizeとoccupationを削除
ALTER TABLE user_demographics DROP COLUMN IF EXISTS household_size;
ALTER TABLE user_demographics DROP COLUMN IF EXISTS occupation;
