-- name: CreatePasswordAuth :exec
INSERT INTO password_auth (
  password_auth_id,
  user_id,
  password_hash,
  salt,
  required_password_change,
  last_changed,
  created_at,
  updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
