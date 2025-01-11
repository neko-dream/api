-- name: CreateTalkSessionConclusion :exec
INSERT INTO talk_session_conclusions (
    talk_session_id,
    content,
    created_by
) VALUES ($1, $2, $3);

-- name: UpdateTalkSessionConclusion :exec
UPDATE talk_session_conclusions
SET
    content = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE talk_session_id = $1;

-- name: GetTalkSessionConclusionByID :one
SELECT
    talk_session_conclusions.talk_session_id,
    talk_session_conclusions.content,
    talk_session_conclusions.created_by,
    talk_session_conclusions.created_at,
    talk_session_conclusions.updated_at,
    users.user_id AS user_id,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url
FROM talk_session_conclusions
LEFT JOIN users
    ON talk_session_conclusions.created_by = users.user_id
WHERE talk_session_id = $1;

-- name: GetConclusionByID :one
SELECT
    sqlc.embed(conclusion),
    sqlc.embed(users)
FROM talk_session_conclusions as conclusion
LEFT JOIN users
    ON conclusion.created_by = users.user_id
WHERE talk_session_id = $1;

-- name: CreateActionItem :exec
INSERT INTO action_items (
    action_item_id,
    talk_session_id,
    sequence,
    content,
    status,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateActionItem :exec
UPDATE action_items
SET
    content = $2,
    status = $3,
    sequence = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE action_item_id = $1;

-- name: GetActionItemsByTalkSessionID :many
SELECT
    action_items.action_item_id,
    action_items.talk_session_id,
    action_items.sequence,
    action_items.content,
    action_items.status,
    action_items.created_at,
    action_items.updated_at,
    users.display_name AS display_name,
    users.display_id AS display,
    users.icon_url AS icon_url
FROM action_items
LEFT JOIN talk_sessions
    ON talk_sessions.talk_session_id = action_items.talk_session_id
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
WHERE action_items.talk_session_id = $1
ORDER BY action_items.sequence;


-- name: GetActionItemByID :one
SELECT
    action_items.action_item_id,
    action_items.talk_session_id,
    action_items.sequence,
    action_items.content,
    action_items.status,
    action_items.created_at,
    action_items.updated_at,
    users.display_name AS display_name,
    users.display_id AS display_id,
    users.icon_url AS icon_url
FROM action_items
LEFT JOIN talk_sessions
    ON talk_sessions.talk_session_id = action_items.talk_session_id
LEFT JOIN users
    ON talk_sessions.owner_id = users.user_id
WHERE action_item_id = $1
ORDER BY action_items.sequence;


-- name: UpdateSequencesByActionItemID :exec
-- 指定したActionItemいよりSequenceが大きいものをすべて+1する
UPDATE action_items
SET
    sequence = sequence + 1
WHERE talk_session_id = $1
    AND sequence >= $2;
