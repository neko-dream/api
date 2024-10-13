-- name: CreateOpinion :exec
INSERT INTO opinions (
    opinion_id,
    talk_session_id,
    user_id,
    parent_opinion_id,
    title,
    content,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetOpinionByID :one
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 親意見に対するユーザーの投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = vote.user_id
WHERE opinions.opinion_id = $1;

-- name: GetOpinionReplies :many
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(parent_vote.vote_type, 0) AS vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 親意見に対する子意見主の投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id
    FROM votes
    WHERE votes.opinion_id = $1
) parent_vote ON opinions.user_id = parent_vote.user_id
WHERE opinions.parent_opinion_id = $1;
