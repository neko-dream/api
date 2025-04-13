-- name: UpdatePasswordAuth :exec
UPDATE password_auth
SET
  password_hash = $2,
  salt = $3,
  last_changed = $4,
  updated_at = $5
WHERE user_id = $1;
