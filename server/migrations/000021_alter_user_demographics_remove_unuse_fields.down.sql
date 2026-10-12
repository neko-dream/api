-- talk_sessionsからhousehold_sizeとoccupationを追加
ALTER TABLE user_demographics ADD COLUMN IF NOT EXISTS household_size SMALLINT;
ALTER TABLE user_demographics ADD COLUMN IF NOT EXISTS occupation SMALLINT;
