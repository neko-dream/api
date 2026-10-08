-- name: SavePasswordLoginChallenge :one
INSERT INTO password_login_challenges (user_id, token_hash, browser_hash, credential_hash, pin_hash, created_at, expires_at)
VALUES (sqlc.arg(user_id), sqlc.arg(token_hash), sqlc.arg(browser_hash), sqlc.arg(credential_hash), sqlc.arg(pin_hash), sqlc.arg(created_at), sqlc.arg(expires_at))
ON CONFLICT (user_id) DO UPDATE
SET token_hash = EXCLUDED.token_hash,
    browser_hash = EXCLUDED.browser_hash,
    credential_hash = EXCLUDED.credential_hash,
    pin_hash = EXCLUDED.pin_hash,
    pin_attempts = 0,
    created_at = EXCLUDED.created_at,
    expires_at = EXCLUDED.expires_at,
    consumed_at = NULL
WHERE password_login_challenges.created_at <= sqlc.arg(resend_before)::timestamptz
RETURNING user_id;

-- name: ConsumePasswordLoginChallenge :one
UPDATE password_login_challenges
SET consumed_at = sqlc.arg(now)::timestamptz
WHERE token_hash = sqlc.arg(token_hash)
  AND browser_hash = sqlc.arg(browser_hash)
  AND expires_at > sqlc.arg(now)::timestamptz
  AND consumed_at IS NULL
  AND pin_attempts < 5
RETURNING user_id, credential_hash;

-- name: FindPasswordLoginChallengeForPIN :one
SELECT * FROM password_login_challenges
WHERE browser_hash = sqlc.arg(browser_hash)
FOR UPDATE;

-- name: RecordPasswordLoginPINFailure :one
UPDATE password_login_challenges
SET pin_attempts = pin_attempts + 1
WHERE browser_hash = sqlc.arg(browser_hash) AND pin_attempts < 5 AND consumed_at IS NULL
RETURNING pin_attempts;
