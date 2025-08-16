-- name: GetUserStats :one
WITH user_activity AS (
    SELECT DISTINCT users.user_id,
        CASE WHEN votes.user_id IS NOT NULL THEN 1 ELSE 0 END as has_voted,
        CASE WHEN opinions.user_id IS NOT NULL THEN 1 ELSE 0 END as has_posted
    FROM users
    LEFT JOIN votes ON users.user_id = votes.user_id
    LEFT JOIN opinions ON users.user_id = opinions.user_id
    WHERE users.user_id != '00000000-0000-0000-0000-000000000001'::uuid
    AND users.display_id IS NOT NULL
)
SELECT
    COUNT(*) as total_users,
    SUM(has_voted) as users_with_votes,
    SUM(has_posted) as users_with_posts,
    SUM(CASE WHEN has_voted = 1 OR has_posted = 1 THEN 1 ELSE 0 END) as active_users
FROM user_activity;

-- name: GetDailyUserStats :many
WITH date_range AS (
    SELECT generate_series(
        CURRENT_DATE - ((sqlc.arg('offset')::integer + sqlc.arg('limit')::integer - 1) * INTERVAL '1 day'),
        CURRENT_DATE - (sqlc.arg('offset')::integer * INTERVAL '1 day'),
        INTERVAL '1 day'
    )::date as activity_date
),
user_activity AS (
    SELECT DISTINCT
        users.user_id,
        DATE_TRUNC('day', COALESCE(votes.created_at, opinions.created_at)) as activity_date,
        CASE WHEN votes.user_id IS NOT NULL THEN 1 ELSE 0 END as has_voted,
        CASE WHEN opinions.user_id IS NOT NULL THEN 1 ELSE 0 END as has_posted
    FROM users
    LEFT JOIN votes ON users.user_id = votes.user_id
    LEFT JOIN opinions ON users.user_id = opinions.user_id
    WHERE users.user_id != '00000000-0000-0000-0000-000000000001'::uuid
    AND COALESCE(votes.created_at, opinions.created_at) IS NOT NULL
    AND users.display_id IS NOT NULL
)
SELECT
    date_range.activity_date,
    COALESCE(COUNT(user_activity.user_id), 0)::integer as total_users,
    COALESCE(SUM(has_voted), 0) as users_with_votes,
    COALESCE(SUM(has_posted), 0) as users_with_posts,
    COALESCE(SUM(CASE WHEN has_voted = 1 OR has_posted = 1 THEN 1 ELSE 0 END), 0)::integer as active_users
FROM date_range
LEFT JOIN user_activity ON date_range.activity_date = user_activity.activity_date::date
GROUP BY date_range.activity_date
ORDER BY date_range.activity_date DESC;

-- name: GetWeeklyUserStats :many
WITH date_range AS (
    SELECT generate_series(
        DATE_TRUNC('week', CURRENT_DATE - ((sqlc.arg('offset')::integer + sqlc.arg('limit')::integer - 1) * INTERVAL '1 week')),
        DATE_TRUNC('week', CURRENT_DATE - (sqlc.arg('offset')::integer * INTERVAL '1 week')),
        INTERVAL '1 week'
    )::date as activity_date
),
user_activity AS (
    SELECT DISTINCT
        users.user_id,
        DATE_TRUNC('week', COALESCE(votes.created_at, opinions.created_at)) as activity_date,
        CASE WHEN votes.user_id IS NOT NULL THEN 1 ELSE 0 END as has_voted,
        CASE WHEN opinions.user_id IS NOT NULL THEN 1 ELSE 0 END as has_posted
    FROM users
    LEFT JOIN votes ON users.user_id = votes.user_id
    LEFT JOIN opinions ON users.user_id = opinions.user_id
    WHERE users.user_id != '00000000-0000-0000-0000-000000000001'::uuid
    AND COALESCE(votes.created_at, opinions.created_at) IS NOT NULL
    AND users.display_id IS NOT NULL
)
SELECT
    date_range.activity_date,
    COALESCE(COUNT(user_activity.user_id), 0)::integer as total_users,
    COALESCE(SUM(has_voted), 0) as users_with_votes,
    COALESCE(SUM(has_posted), 0) as users_with_posts,
    COALESCE(SUM(CASE WHEN has_voted = 1 OR has_posted = 1 THEN 1 ELSE 0 END), 0)::integer as active_users
FROM date_range
LEFT JOIN user_activity ON date_range.activity_date = user_activity.activity_date::date
GROUP BY date_range.activity_date
ORDER BY date_range.activity_date DESC;
