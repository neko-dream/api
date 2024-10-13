
-- opinions テーブルのインデックスを追加
CREATE INDEX idx_opinions_talk_session_id ON opinions(talk_session_id);
CREATE INDEX idx_opinions_user_id ON opinions(user_id);
CREATE INDEX idx_opinions_parent_opinion_id ON opinions(parent_opinion_id);
CREATE INDEX idx_opinions_opinion_id_parent_opinion_id ON opinions(opinion_id, parent_opinion_id);

-- users テーブルのインデックスを追加
CREATE INDEX idx_users_user_id ON users(user_id);

-- votes テーブルのインデックスを追加
CREATE INDEX idx_votes_opinion_id_user_id ON votes(opinion_id, user_id);
CREATE INDEX idx_votes_user_id_opinion_id ON votes(user_id, opinion_id);
CREATE INDEX idx_votes_vote_id_opinion_id ON votes(vote_id, opinion_id);

-- sessions テーブルのインデックスを追加
CREATE INDEX idx_session_id_user_id ON sessions(user_id, session_id);

-- user_auths テーブルのインデックスを追加
CREATE INDEX idx_user_id_user_subject ON user_auths(user_id, subject);

-- user_demographics テーブルのインデックスを追加
CREATE INDEX idx_user_demographics_user_id ON user_demographics(user_id);
