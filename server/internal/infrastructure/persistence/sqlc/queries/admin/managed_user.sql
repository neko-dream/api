-- name: ListManagedUsers :many
SELECT sqlc.embed(users),
  EXISTS(SELECT 1 FROM password_setup_tokens WHERE password_setup_tokens.user_id = users.user_id)
    AND NOT EXISTS(SELECT 1 FROM password_auth WHERE password_auth.user_id = users.user_id) AS password_setup_required,
  ARRAY(SELECT provider FROM user_auths WHERE user_auths.user_id = users.user_id
        UNION SELECT 'PASSWORD'::varchar FROM password_auth WHERE password_auth.user_id = users.user_id ORDER BY provider)::text[] AS providers
FROM users
WHERE sqlc.arg(search)::text = ''
 OR strpos(lower(coalesce(display_name, '')), lower(sqlc.arg(search)::text)) > 0
 OR strpos(lower(coalesce(display_id, '')), lower(sqlc.arg(search)::text)) > 0
 OR strpos(user_id::text, lower(sqlc.arg(search)::text)) > 0
ORDER BY created_at DESC, user_id DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountManagedUsers :one
SELECT count(*) FROM users
WHERE sqlc.arg(search)::text = ''
 OR strpos(lower(coalesce(display_name, '')), lower(sqlc.arg(search)::text)) > 0
 OR strpos(lower(coalesce(display_id, '')), lower(sqlc.arg(search)::text)) > 0
 OR strpos(user_id::text, lower(sqlc.arg(search)::text)) > 0;

-- name: GetManagedUser :one
SELECT sqlc.embed(users),
  EXISTS(SELECT 1 FROM password_setup_tokens WHERE password_setup_tokens.user_id = users.user_id)
    AND NOT EXISTS(SELECT 1 FROM password_auth WHERE password_auth.user_id = users.user_id) AS password_setup_required,
  ARRAY(SELECT provider FROM user_auths WHERE user_auths.user_id = users.user_id
        UNION SELECT 'PASSWORD'::varchar FROM password_auth WHERE password_auth.user_id = users.user_id ORDER BY provider)::text[] AS providers
FROM users WHERE users.user_id = $1;
