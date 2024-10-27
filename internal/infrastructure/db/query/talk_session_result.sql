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

-- name: CreateActionItem :exec
INSERT INTO action_items (
    action_item_id,
    talk_session_id,
    parent_action_item_id,
    content,
    status
) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateActionItem :exec
UPDATE action_items
SET
    content = $2,
    status = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE action_item_id = $1;

-- name: GetActionItemsByTalkSessionID :many
-- トークセッションに紐づくアクションアイテムを再起的に取得
WITH RECURSIVE action_item_tree AS (
    SELECT
        action_items.action_item_id,
        action_items.talk_session_id,
        action_items.parent_action_item_id,
        action_items.content,
        action_items.status,
        action_items.created_at,
        action_items.updated_at,
        1 AS depth
    FROM action_items
    WHERE action_items.talk_session_id = $1
    AND action_items.parent_action_item_id IS NULL

    UNION ALL

    SELECT
        action_items.action_item_id,
        action_items.talk_session_id,
        action_items.parent_action_item_id,
        action_items.content,
        action_items.status,
        action_items.created_at,
        action_items.updated_at,
        action_item_tree.depth + 1
    FROM action_items
    JOIN action_item_tree
        ON action_items.parent_action_item_id = action_item_tree.action_item_id
)
SELECT
    action_item_tree.action_item_id,
    action_item_tree.talk_session_id,
    action_item_tree.parent_action_item_id,
    action_item_tree.content,
    action_item_tree.status,
    action_item_tree.created_at,
    action_item_tree.updated_at,
    action_item_tree.depth,
    users.display_name AS display_name,
    users.display_id AS display,
    users.icon_url AS icon_url
FROM action_item_tree
LEFT JOIN users
    ON action_item_tree.created_by = users.user_id
ORDER BY action_item_tree.created_at ASC;
