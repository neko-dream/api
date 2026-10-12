-- name: GetSessionCreationPermission :one
SELECT can_create_talk_session FROM users WHERE user_id = $1;

-- name: SetSessionCreationPermission :one
UPDATE users SET can_create_talk_session = $2, updated_at = now()
WHERE users.user_id = $1
  AND (NOT $2::boolean OR (withdrawal_date IS NULL
    AND EXISTS (SELECT 1 FROM user_auths WHERE user_auths.user_id = users.user_id AND provider <> 'GUEST')))
RETURNING can_create_talk_session;

-- name: GetCurrentUser :one
SELECT users.user_id, users.display_id, users.display_name, users.icon_url,
       users.withdrawal_date, users.can_create_talk_session,
       NOT EXISTS (SELECT 1 FROM user_auths WHERE user_auths.user_id = users.user_id AND provider <> 'GUEST') AS is_guest
FROM users WHERE users.user_id = $1;
