-- name: DeletePasswordAuth :exec
DELETE FROM password_auth
WHERE user_id = $1;
