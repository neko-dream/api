-- name: CreateDomainEvent :one
INSERT INTO domain_events (
    id,
    event_type,
    event_data,
    aggregate_id,
    aggregate_type,
    status,
    occurred_at,
    retry_count
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) RETURNING *;

-- name: GetUnprocessedEvents :many
SELECT * FROM domain_events
WHERE status IN ('pending', 'failed')
  AND (status != 'failed' OR retry_count < $1)
  AND ($2::text[] IS NULL OR event_type = ANY($2::text[]))
ORDER BY occurred_at ASC
LIMIT $3
FOR UPDATE SKIP LOCKED;

-- name: MarkEventAsProcessed :one
UPDATE domain_events
SET status = 'processed',
    processed_at = NOW()
WHERE id = $1
  AND status IN ('pending', 'processing', 'failed')
RETURNING *;

-- name: MarkEventAsFailed :one
UPDATE domain_events
SET status = 'failed',
    failed_at = NOW(),
    failure_reason = $2,
    retry_count = retry_count + 1
WHERE id = $1
  AND status IN ('pending', 'processing')
RETURNING *;

-- name: GetEventsByAggregateID :many
SELECT * FROM domain_events
WHERE aggregate_id = $1
  AND aggregate_type = $2
ORDER BY occurred_at ASC;