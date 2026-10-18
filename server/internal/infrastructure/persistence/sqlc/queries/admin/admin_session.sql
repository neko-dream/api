-- name: CreateAdminSession :exec
INSERT INTO admin_sessions (admin_session_id, admin_user_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetAdminSessionByTokenHash :one
-- 現在のトークンか、ローテーション直前のトークンで探す（猶予期間の判定は呼び出し側で行う）
SELECT * FROM admin_sessions
WHERE token_hash = sqlc.arg(token_hash) OR previous_token_hash = sqlc.arg(token_hash)
LIMIT 1;

-- name: RotateAdminSession :execrows
UPDATE admin_sessions
SET previous_token_hash = token_hash,
    rotated_at = sqlc.arg(rotated_at),
    token_hash = sqlc.arg(new_token_hash),
    expires_at = sqlc.arg(expires_at),
    updated_at = CURRENT_TIMESTAMP
WHERE admin_session_id = sqlc.arg(admin_session_id)
  AND token_hash = sqlc.arg(old_token_hash)
  AND expires_at > CURRENT_TIMESTAMP;

-- name: DeleteAdminSessionByTokenHash :exec
-- ローテーション直後に古いトークンでログアウトされても、セッションを消せるようにする
DELETE FROM admin_sessions WHERE token_hash = $1 OR previous_token_hash = $1;
