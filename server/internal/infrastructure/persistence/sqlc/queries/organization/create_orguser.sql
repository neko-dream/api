-- name: CreateOrgUser :exec
INSERT INTO organization_users (
    organization_user_id,
    user_id,
    organization_id,
    role,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6);
