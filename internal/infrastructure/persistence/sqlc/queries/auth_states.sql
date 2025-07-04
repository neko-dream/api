-- name: CreateAuthState :one
INSERT INTO auth_states (
    state,
    provider,
    redirect_url,
    expires_at,
    registration_url,
    organization_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetAuthState :one
SELECT * FROM auth_states
WHERE state = $1 AND expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: DeleteAuthState :exec
DELETE FROM auth_states
WHERE state = $1;

-- name: DeleteExpiredAuthStates :exec
DELETE FROM auth_states
WHERE expires_at <= CURRENT_TIMESTAMP;
