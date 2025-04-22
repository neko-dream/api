-- name: GetRandomSeedOpinions :many
SELECT
    sqlc.embed(opinions),
    sqlc.embed(users),
    COALESCE(rc.reply_count, 0) AS reply_count
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 指定されたユーザーが投票していない意見のみを取得
LEFT JOIN (
    SELECT opinions.opinion_id
    FROM opinions
    LEFT JOIN votes
        ON opinions.opinion_id = votes.opinion_id
        AND votes.user_id = $2
    GROUP BY opinions.opinion_id
    HAVING COUNT(votes.vote_id) = 0
) vote_count ON opinions.opinion_id = vote_count.opinion_id
-- この意見に対するリプライ数
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON rc.opinion_id = opinions.opinion_id
-- トークセッションに紐づく意見のみを取得
WHERE opinions.talk_session_id = $1
    AND vote_count.opinion_id = opinions.opinion_id
    AND opinions.parent_opinion_id IS NULL
    AND opinions.user_id = '00000000-0000-0000-0000-000000000001'::uuid
LIMIT $3;

