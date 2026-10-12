-- name: GetGroupInfoByTalkSessionId :many
SELECT
    user_group_info.pos_x,
    user_group_info.pos_y,
    user_group_info.group_id,
    user_group_info.perimeter_index,
    users.display_id AS display_id,
    users.display_name AS display_name,
    users.icon_url AS icon_url,
    user_group_info.user_id
FROM user_group_info
LEFT JOIN users
    ON user_group_info.user_id = users.user_id
WHERE talk_session_id = $1;

-- name: GetRepresentativeOpinionsByTalkSessionId :many
SELECT
    sqlc.embed(representative_opinions),
    sqlc.embed(opinions),
    sqlc.embed(users),
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
WHERE representative_opinions.rank < 4
    AND opinions.talk_session_id = $1
ORDER BY representative_opinions.rank;

-- name: GetGroupListByTalkSessionId :many
SELECT
    DISTINCT user_group_info.group_id
FROM user_group_info
WHERE talk_session_id = $1;

-- name: GetReportByTalkSessionId :one
SELECT
    -- talk_session_report_history_id as analysis_report_history_id,
    talk_session_id,
    report,
    created_at
FROM talk_session_reports
WHERE talk_session_id = $1;

-- name: FindReportByID :one
SELECT
    -- talk_session_report_history_id as analysis_report_history_id,
    talk_session_id,
    report,
    created_at
FROM talk_session_reports
WHERE talk_session_id = $1;

-- name: GetFeedbackByReportHistoryID :many
SELECT
    report_feedback_id,
    user_id,
    feedback_type,
    created_at
FROM report_feedback
WHERE talk_session_report_history_id = $1;

-- name: SaveReportFeedback :exec
INSERT INTO report_feedback (
    report_feedback_id,
    talk_session_report_history_id,
    user_id,
    feedback_type,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) ON CONFLICT (user_id, talk_session_report_history_id) DO NOTHING;

-- name: AddGeneratedImages :exec
INSERT INTO talk_session_generated_images (talk_session_id, wordmap_url, tsnc_url) VALUES ($1, $2, $3)
ON CONFLICT (talk_session_id) DO UPDATE SET wordmap_url = $2, tsnc_url = $3, updated_at = NOW();

-- name: GetGeneratedImages :one
SELECT
    talk_session_id,
    wordmap_url,
    tsnc_url,
    created_at,
    updated_at
FROM talk_session_generated_images
WHERE talk_session_id = $1::uuid;
