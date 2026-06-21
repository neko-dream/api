-- name: UpdatePasswordAuth :exec
UPDATE password_auth
SET
  password_hash = $2,
  salt = $3,
  required_password_change = $4,
  last_changed = $5,
  updated_at = $6
WHERE user_id = $1;
