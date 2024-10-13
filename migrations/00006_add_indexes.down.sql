-- opinions テーブルのインデックスを削除
DROP INDEX IF EXISTS idx_opinions_talk_session_id ON opinions;
DROP INDEX IF EXISTS idx_opinions_user_id ON opinions;
DROP INDEX IF EXISTS idx_opinions_parent_opinion_id ON opinions;
DROP INDEX IF EXISTS idx_opinions_opinion_id_parent_opinion_id ON opinions;

-- users テーブルのインデックスを削除
DROP INDEX IF EXISTS idx_users_user_id ON users;

-- votes テーブルのインデックスを削除
DROP INDEX IF EXISTS idx_votes_opinion_id_user_id ON votes;
DROP INDEX IF EXISTS idx_votes_user_id_opinion_id ON votes;

-- sessions テーブルのインデックスを削除
DROP INDEX IF EXISTS idx_session_id_user_id ON `sessions`;

-- user_auths テーブルのインデックスを削除
DROP INDEX IF EXISTS idx_user_id_user_subject ON user_auths;

-- user_demographics テーブルのインデックスを削除
DROP INDEX IF EXISTS idx_user_demographics_user_id ON user_demographics;
