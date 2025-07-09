-- Drop indexes
DROP INDEX IF EXISTS idx_domain_events_aggregate;
DROP INDEX IF EXISTS idx_domain_events_status_event_type_occurred;

-- Drop domain_events table
DROP TABLE IF EXISTS domain_events;