-- Create domain_events table for event sourcing
CREATE TABLE domain_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    event_data JSONB NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'processing', 'processed', 'failed')),
    occurred_at TIMESTAMP NOT NULL,
    processed_at TIMESTAMP,
    failed_at TIMESTAMP,
    failure_reason TEXT,
    retry_count INT NOT NULL DEFAULT 0
);

-- Create indexes for efficient querying
CREATE INDEX idx_domain_events_status_event_type_occurred ON domain_events(status, event_type, occurred_at);
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_id, aggregate_type);
