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
WITH user_activity AS (
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
    activity_date::date,
    COUNT(*) as total_users,
    SUM(has_voted) as users_with_votes,
    SUM(has_posted) as users_with_posts,
    SUM(CASE WHEN has_voted = 1 OR has_posted = 1 THEN 1 ELSE 0 END) as active_users
FROM user_activity
GROUP BY activity_date
ORDER BY activity_date DESC
LIMIT $1 OFFSET $2;

-- name: GetWeeklyUserStats :many
WITH user_activity AS (
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
    activity_date::date,
    COUNT(*) as total_users,
    SUM(has_voted) as users_with_votes,
    SUM(has_posted) as users_with_posts,
    SUM(CASE WHEN has_voted = 1 OR has_posted = 1 THEN 1 ELSE 0 END) as active_users
FROM user_activity
GROUP BY activity_date
ORDER BY activity_date DESC
LIMIT $1 OFFSET $2;
