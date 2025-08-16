-- name: GetNotificationPreference :one
SELECT * FROM notification_preferences
WHERE user_id = $1;

-- name: CreateNotificationPreference :one
INSERT INTO notification_preferences (
    user_id,
    push_notification_enabled
) VALUES (
    $1, $2
) RETURNING *;

-- name: UpdateNotificationPreference :one
UPDATE notification_preferences SET
    push_notification_enabled = $2
WHERE user_id = $1
RETURNING *;

-- name: UpsertNotificationPreference :one
INSERT INTO notification_preferences (
    user_id,
    push_notification_enabled
) VALUES (
    $1, $2
) ON CONFLICT (user_id) DO UPDATE SET
    push_notification_enabled = EXCLUDED.push_notification_enabled
RETURNING *;

-- name: GetNotificationPreferencesByUserIDs :many
SELECT * FROM notification_preferences
WHERE user_id = ANY($1::uuid[]);
