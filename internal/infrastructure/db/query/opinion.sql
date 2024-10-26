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
) cv ON opinions.user_id = sqlc.narg('user_id')::uuid
    AND opinions.opinion_id = cv.opinion_id
WHERE opinions.opinion_id = $1;

-- name: GetOpinionReplies :many
SELECT
    DISTINCT opinions.opinion_id,
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
) cv ON opinions.user_id = sqlc.narg('user_id')::uuid
    AND opinions.opinion_id = cv.opinion_id
WHERE opinions.parent_opinion_id = $1
GROUP BY opinions.opinion_id, users.display_name, users.display_id, users.icon_url, pv.vote_type, cv.vote_type;

-- name: GetRandomOpinions :many
-- SELECT
--     opinions.opinion_id,
--     opinions.talk_session_id,
--     opinions.user_id,
--     opinions.parent_opinion_id,
--     opinions.title,
--     opinions.content,
--     opinions.reference_url,
--     opinions.picture_url,
--     opinions.created_at,
--     users.display_name AS display_name,
--     users.display_id AS display_id,
--     users.icon_url AS icon_url,
--     COALESCE(pv.vote_type, 0) AS vote_type,
--     COALESCE(rc.reply_count, 0) AS reply_count
-- FROM opinions
-- LEFT JOIN users
--     ON opinions.user_id = users.user_id
-- -- 親意見に対する意見投稿ユーザーの意思を取得
-- LEFT JOIN (
--     SELECT votes.vote_type, votes.user_id, votes.opinion_id
--     FROM votes
-- ) pv ON pv.opinion_id = opinions.parent_opinion_id
--     AND pv.user_id = opinions.user_id
-- -- 指定されたユーザーが投票していない意見のみを取得
-- LEFT JOIN (
--     SELECT opinions.opinion_id
--     FROM opinions
--     LEFT JOIN votes
--         ON opinions.opinion_id = votes.opinion_id
--         AND votes.user_id = $1
--     GROUP BY opinions.opinion_id
--     HAVING COUNT(votes.vote_id) = 0
-- ) vote_count ON opinions.opinion_id = vote_count.opinion_id
-- -- この意見に対するリプライ数
-- LEFT JOIN (
--     SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
--     FROM opinions
--     GROUP BY parent_opinion_id
-- ) rc ON rc.opinion_id = opinions.opinion_id
-- -- グループ内のランクを取得
-- LEFT JOIN (
--     SELECT rank, opinion_id
--     FROM representative_opinions
-- ) ro ON opinions.opinion_id = ro.opinion_id
-- -- トークセッションに紐づく意見のみを取得
-- WHERE opinions.talk_session_id = $2
--     AND vote_count.opinion_id = opinions.opinion_id
-- ORDER BY
--     CASE sqlc.narg('sort_key')::text
--         WHEN 'top' THEN COALESCE(ro.rank, 0)
--         ELSE RANDOM()
--     END ASC
-- LIMIT $3;
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
        AND votes.user_id = $1
    GROUP BY opinions.opinion_id
    HAVING COUNT(votes.vote_id) = 0
) vote_count ON opinions.opinion_id = vote_count.opinion_id
-- この意見に対するリプライ数
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON rc.opinion_id = opinions.opinion_id
-- グループ内のランクを取得
LEFT JOIN (
    SELECT rank, opinion_id
    FROM representative_opinions
) ro ON opinions.opinion_id = ro.opinion_id
-- トークセッションに紐づく意見のみを取得
WHERE opinions.talk_session_id = $2
    AND vote_count.opinion_id = opinions.opinion_id
    AND opinions.parent_opinion_id IS NULL
ORDER BY
    CASE sqlc.narg('sort_key')::text
        WHEN 'top' THEN COALESCE(ro.rank, 0)
        ELSE RANDOM()
    END ASC
LIMIT $3;


-- name: GetOpinionsByUserID :many
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
    -- 意見に対するリプライ数（再帰）
    COALESCE(rc.reply_count, 0) AS reply_count
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
WHERE opinions.user_id = $1
-- latest, mostReply, oldestでソート
ORDER BY
    CASE sqlc.narg('sort_key')::text
        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'mostReply' THEN reply_count
    END ASC
LIMIT $2 OFFSET $3
;


-- name: GetOpinionsByTalkSessionID :many
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
    -- 意見に対するリプライ数（再帰）
    COALESCE(rc.reply_count, 0) AS reply_count,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) cv ON opinions.user_id = sqlc.narg('user_id')::uuid
    AND opinions.opinion_id = cv.opinion_id
WHERE opinions.talk_session_id = $1
-- latest, mostReply, oldestでソート
ORDER BY
    CASE sqlc.narg('sort_key')::text
        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'mostReply' THEN reply_count
    END ASC
LIMIT $2 OFFSET $3;

-- name: GetParentOpinions :many
WITH RECURSIVE opinion_tree AS (
    -- ベースケース：指定された意見から開始
    SELECT
        o.opinion_id,
        o.parent_opinion_id,
        1 as level
    FROM opinions o
    WHERE o.opinion_id = $1

    UNION ALL

    SELECT
        p.opinion_id,
        p.parent_opinion_id,
        t.level + 1
    FROM opinions p
    INNER JOIN opinion_tree t ON t.parent_opinion_id = p.opinion_id
)
SELECT
    o.opinion_id,
    o.talk_session_id,
    o.user_id,
    o.parent_opinion_id,
    o.title,
    o.content,
    o.reference_url,
    o.picture_url,
    o.created_at,
    u.display_name AS display_name,
    u.display_id AS display_id,
    u.icon_url AS icon_url,
    COALESCE(pv.vote_type, 0) AS vote_type,
    COALESCE(rc.reply_count, 0) AS reply_count,
    COALESCE(cv.vote_type, 0) AS current_vote_type,
    ot.level
FROM opinion_tree ot
JOIN opinions o ON ot.opinion_id = o.opinion_id
LEFT JOIN users u ON o.user_id = u.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON o.parent_opinion_id = pv.opinion_id
    AND o.user_id = pv.user_id
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) cv ON opinions.user_id = sqlc.narg('user_id')::uuid
    AND opinions.opinion_id = cv.opinion_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON o.opinion_id = rc.parent_opinion_id;

-- name: CountOpinions :one
SELECT
    COUNT(opinions.*) AS opinion_count
FROM opinions
WHERE
    CASE
        WHEN sqlc.narg('user_id')::uuid IS NOT NULL THEN opinions.user_id = sqlc.narg('user_id')::uuid
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('talk_session_id')::uuid IS NOT NULL THEN opinions.talk_session_id = sqlc.narg('talk_session_id')::uuid
        ELSE TRUE
    END
    AND
    CASE
        WHEN sqlc.narg('parent_opinion_id')::uuid IS NOT NULL THEN opinions.parent_opinion_id = sqlc.narg('parent_opinion_id')::uuid
        ELSE TRUE
    END
;
