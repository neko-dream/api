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
    sqlc.embed(opinions),
    sqlc.embed(users),
    COALESCE(cv.vote_type, 0) AS current_vote_type,
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(rc.reply_count, 0) AS reply_count
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 親意見に対するユーザーの投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON pv.opinion_id = opinions.parent_opinion_id
    AND  pv.user_id = opinions.user_id
-- この意見に対するリプライ数
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON rc.opinion_id = opinions.opinion_id
-- ユーザーIDが提供された場合、そのユーザーの投票ステータスを一緒に取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
    WHERE votes.user_id = sqlc.narg('user_id')::uuid
) cv ON opinions.opinion_id = cv.opinion_id
WHERE opinions.opinion_id = $1;

-- name: GetOpinionReplies :many
SELECT
    DISTINCT opinions.opinion_id,
    sqlc.embed(opinions),
    sqlc.embed(users),
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
-- 親意見に対する子意見主の投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
    WHERE votes.opinion_id = $1
) pv ON opinions.user_id = pv.user_id
    AND opinions.opinion_id = pv.opinion_id
-- ユーザーIDが提供された場合、そのユーザーの投票ステータスを取得
LEFT JOIN (
    SELECT votes.vote_type, votes.opinion_id
    FROM votes
    WHERE votes.user_id = sqlc.narg('user_id')::uuid
) cv ON opinions.opinion_id = cv.opinion_id
WHERE opinions.parent_opinion_id = $1
GROUP BY opinions.opinion_id,users.user_id, users.display_name, users.display_id, users.icon_url, pv.vote_type, cv.vote_type
ORDER BY opinions.created_at DESC;

-- name: GetRandomOpinions :many
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
-- 通報された意見を除外
LEFT JOIN (
	SELECT DISTINCT opinion_reports.opinion_id, opinion_reports.status
	FROM opinion_reports
	WHERE opinion_reports.talk_session_id = $2
) opr ON opr.opinion_id = opinions.opinion_id
-- トークセッションに紐づく意見のみを取得
WHERE opinions.talk_session_id = $2
    AND vote_count.opinion_id = opinions.opinion_id
    -- exclude_opinion_idsが空でない場合、除外する意見を指定
    AND (
        CASE
            WHEN sqlc.arg('excludes_len')::int = 0 THEN TRUE
            ELSE opinions.opinion_id != ANY(sqlc.narg('exclude_opinion_ids')::uuid[])
        END
    )
    -- 親意見がないものを取得
    AND opinions.parent_opinion_id IS NULL
    -- 削除されたものはスワイプ意見から除外
    AND (opr.opinion_id IS NULL OR opr.status != 'deleted')
ORDER BY RANDOM()
LIMIT $3;

-- name: GetOpinionsByRank :many
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
        AND votes.user_id = sqlc.arg('user_id')::uuid
    GROUP BY opinions.opinion_id
    HAVING COUNT(votes.vote_id) = 0
) vote_count ON opinions.opinion_id = vote_count.opinion_id
-- 親意見に対するユーザーの投票を取得
LEFT JOIN (
    SELECT votes.vote_type, votes.user_id, votes.opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND  pv.user_id = opinions.user_id
-- この意見に対するリプライ数
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id as opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON rc.opinion_id = opinions.opinion_id
-- 通報された意見を除外
LEFT JOIN (
	SELECT DISTINCT opinion_reports.opinion_id, opinion_reports.status
	FROM opinion_reports
	WHERE opinion_reports.talk_session_id = sqlc.arg('talk_session_id')::uuid
) opr ON opr.opinion_id = opinions.opinion_id
LEFT JOIN representative_opinions ON opinions.opinion_id = representative_opinions.opinion_id
WHERE opinions.talk_session_id = sqlc.arg('talk_session_id')::uuid
    AND vote_count.opinion_id = opinions.opinion_id
    AND opinions.parent_opinion_id IS NULL
    -- 削除されたものはスワイプ意見から除外
    AND (opr.opinion_id IS NULL OR opr.status != 'deleted')
    AND representative_opinions.rank = sqlc.arg('rank')::int
LIMIT sqlc.arg('limit')::int
;

-- name: CountSwipeableOpinions :one
SELECT COUNT(vote_count.opinion_id) AS random_opinion_count
FROM opinions
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
-- 通報された意見を除外
LEFT JOIN (
	SELECT DISTINCT opinion_reports.opinion_id, opinion_reports.status
	FROM opinion_reports
	WHERE opinion_reports.talk_session_id = $2
) opr ON opr.opinion_id = opinions.opinion_id
-- トークセッションに紐づく意見のみを取得
WHERE opinions.talk_session_id = $2
    AND vote_count.opinion_id = opinions.opinion_id
    AND opinions.parent_opinion_id IS NULL
    AND (opr.opinion_id IS NULL OR opr.status != 'deleted');

-- name: GetOpinionsByUserID :many
SELECT
    sqlc.embed(opinions),
    sqlc.embed(users),
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
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
        WHEN 'mostReplies' THEN reply_count
    END DESC
LIMIT $2 OFFSET $3
;

-- name: GetOpinionsByTalkSessionID :many
WITH unique_opinions AS (
    SELECT DISTINCT ON (opinions.opinion_id)
        opinions.*
    FROM opinions
    WHERE opinions.talk_session_id = $1
)
SELECT
    sqlc.embed(opinions),
    sqlc.embed(users),
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
    COALESCE(rc.reply_count, 0) AS reply_count,
    COALESCE(cv.vote_type, 0) AS current_vote_type
FROM unique_opinions opinions
LEFT JOIN users ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT DISTINCT ON (opinion_id) vote_type, user_id, opinion_id
    FROM votes
) pv ON opinions.parent_opinion_id = pv.opinion_id
    AND opinions.user_id = pv.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
LEFT JOIN (
    SELECT DISTINCT ON (opinion_id) vote_type, user_id, opinion_id
    FROM votes
    WHERE user_id = sqlc.narg('user_id')::uuid
) cv ON opinions.opinion_id = cv.opinion_id
WHERE opinions.parent_opinion_id IS NULL
    -- IsSeedがtrueの場合、ユーザーIDが00000000-0000-0000-0000-000000000001の意見のみを取得
    AND (
        CASE
            WHEN sqlc.narg('is_seed')::boolean IS TRUE THEN opinions.user_id = '00000000-0000-0000-0000-000000000001'::uuid
            ELSE TRUE
        END
    )
ORDER BY
    CASE sqlc.narg('sort_key')::text
        WHEN 'latest' THEN EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'oldest' THEN EXTRACT(EPOCH FROM TIMESTAMP '2199-12-31 23:59:59') - EXTRACT(EPOCH FROM opinions.created_at)
        WHEN 'mostReplies' THEN COALESCE(rc.reply_count, 0)
    END DESC
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
    sqlc.embed(o),
    sqlc.embed(u),
    COALESCE(pv.vote_type, 0) AS parent_vote_type,
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
) cv ON o.opinion_id = cv.opinion_id
    AND cv.user_id = sqlc.narg('user_id')::uuid
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON o.opinion_id = rc.parent_opinion_id
ORDER BY ot.level DESC;

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
