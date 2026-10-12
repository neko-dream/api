-- name: CreateDevice :one
INSERT INTO devices (
    device_id,
    user_id,
    device_token,
    platform,
    device_name,
    app_version,
    os_version,
    enabled,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
) RETURNING *;

-- name: UpdateDevice :one
UPDATE devices SET
    device_token = $2,
    platform = $3,
    device_name = COALESCE($4, device_name),
    app_version = COALESCE($5, app_version),
    os_version = COALESCE($6, os_version),
    enabled = $7,
    last_active_at = CURRENT_TIMESTAMP,
    updated_at = $8
WHERE device_id = $1
RETURNING *;

-- name: GetDeviceByID :one
SELECT * FROM devices
WHERE device_id = $1;

-- name: GetDevicesByUserID :many
SELECT * FROM devices
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetActiveDevicesByUserID :many
SELECT * FROM devices
WHERE user_id = $1 AND enabled = true
ORDER BY created_at DESC;

-- name: GetActiveDevicesByUserIDs :many
SELECT * FROM devices
WHERE user_id = ANY($1::uuid[]) AND enabled = true
ORDER BY user_id, created_at DESC;

-- name: InvalidateDevice :exec
UPDATE devices SET
    enabled = false,
    updated_at = CURRENT_TIMESTAMP
WHERE device_id = $1;

-- name: InvalidateDeviceByToken :exec
UPDATE devices SET
    enabled = false,
    updated_at = CURRENT_TIMESTAMP
WHERE device_token = $1 AND platform = $2;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE device_id = $1;

-- name: DeleteDeviceByUserID :exec
DELETE FROM devices
WHERE device_id = $1 AND user_id = $2;

-- name: GetDeviceByToken :one
SELECT * FROM devices
WHERE device_token = $1;

-- name: UpdateDeviceActivity :exec
UPDATE devices SET
    last_active_at = CURRENT_TIMESTAMP
WHERE device_id = $1;

-- name: DeleteDevicesByUserID :exec
DELETE FROM devices
WHERE user_id = $1;

-- name: UpsertDevice :one
INSERT INTO devices (
    device_id,
    user_id,
    device_token,
    platform,
    device_name,
    app_version,
    os_version,
    enabled,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
) ON CONFLICT (user_id, device_token, platform) DO UPDATE SET
    device_name = COALESCE(EXCLUDED.device_name, devices.device_name),
    app_version = COALESCE(EXCLUDED.app_version, devices.app_version),
    os_version = COALESCE(EXCLUDED.os_version, devices.os_version),
    enabled = EXCLUDED.enabled,
    last_active_at = CURRENT_TIMESTAMP,
    updated_at = EXCLUDED.updated_at
RETURNING *;


-- name: GetAllActiveDevices :many
SELECT * FROM devices
WHERE enabled = true
ORDER BY user_id, created_at DESC;

-- name: ExistsByDeviceTokenAndPlatform :one
SELECT EXISTS (
    SELECT 1 FROM devices
    WHERE device_token = $1 AND platform = $2
) AS exists;
