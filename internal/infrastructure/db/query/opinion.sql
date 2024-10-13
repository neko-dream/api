-- name: CreateOpinion :exec
INSERT INTO opinions (
    opinion_id,
    talk_session_id,
    user_id,
    parent_opinion_id,
    title,
    content,
    reference_url,
    picture_url,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetOpinionByID :one
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.reference_url,
    opinions.picture_url,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS vote_type,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 親意見に対するユーザーの投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
-- ユーザーIDが提供された場合、そのユーザーの投票ステータスを一緒に取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) cv ON opinions.user_id = COALESCE(sqlc.narg('user_id'), opinions.user_id)
    AND opinions.opinion_id = cv.opinion_id
WHERE opinions.opinion_id = $1;

-- name: GetOpinionReplies :many
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.reference_url,
    opinions.picture_url,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS vote_type,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 親意見に対する子意見主の投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id
    FROM votes
    WHERE votes.opinion_id = $1
) pv ON opinions.user_id = pv.user_id
-- ユーザーIDが提供された場合、そのユーザーの投票ステータスを一緒に取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) cv ON opinions.user_id = COALESCE(sqlc.narg('user_id'), opinions.user_id)
    AND opinions.opinion_id = cv.opinion_id
WHERE opinions.parent_opinion_id = $1;

-- name: GetRandomOpinions :many
SELECT
    opinions.opinion_id,
    opinions.talk_session_id,
    opinions.user_id,
    opinions.parent_opinion_id,
    opinions.title,
    opinions.content,
    opinions.reference_url,
    opinions.picture_url,
    opinions.created_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS vote_type
    -- 意見に対するリプライ数（再帰）
    -- 0 AS reply_count
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
-- 指定されたユーザーが投票していない意見のみを取得
LEFT JOIN (
    SELECT opinions.opinion_id
    FROM opinions
    LEFT JOIN votes
        ON opinions.opinion_id = votes.opinion_id
        AND votes.user_id = $1
    GROUP BY opinions.opinion_id
    HAVING COUNT(votes.vote_id) = 0
) vote_count ON opinions.opinion_id = vote_count.opinion_id
WHERE opinions.talk_session_id = $2
    AND vote_count.opinion_id = opinions.opinion_id
ORDER BY RANDOM()
LIMIT $3;

