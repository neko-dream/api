-- name: CreateReport :exec
INSERT INTO opinion_reports (
    opinion_report_id,
    opinion_id,
    talk_session_id,
    reporter_id,
    reason,
    status,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    NOW(),
    NOW()
) RETURNING *;
