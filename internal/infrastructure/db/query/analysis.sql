-- name: GetGroupInfoByTalkSessionId :many
SELECT
    user_group_info.pos_x,
    user_group_info.pos_y,
    user_group_info.group_id,
    users.display_id AS display_id,
    user_group_info.user_id
FROM user_group_info
LEFT JOIN users
    ON user_group_info.user_id = users.user_id
WHERE talk_session_id = $1;

-- name: GetRepresentativeOpinionsByTalkSessionId :many
SELECT
    representative_opinions.group_id,
    representative_opinions.rank,
    opinions.opinion_id,
    opinions.talk_session_id,
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
FROM representative_opinions
LEFT JOIN opinions
    ON representative_opinions.opinion_id = opinions.opinion_id
LEFT JOIN users
    ON opinions.user_id = users.user_id
LEFT JOIN (
    SELECT COUNT(opinion_id) AS reply_count, parent_opinion_id
    FROM opinions
    GROUP BY parent_opinion_id
) rc ON opinions.opinion_id = rc.parent_opinion_id
WHERE representative_opinions.rank < 2
    AND opinions.talk_session_id = $1;

-- name: GetGroupListByTalkSessionId :many
SELECT
    DISTINCT user_group_info.group_id
FROM user_group_info
WHERE talk_session_id = $1;
