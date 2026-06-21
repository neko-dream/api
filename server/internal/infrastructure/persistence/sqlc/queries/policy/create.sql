-- name: CreatePolicy :exec
INSERT INTO policy_versions (
    version,
    created_at
) VALUES ($1, $2);
