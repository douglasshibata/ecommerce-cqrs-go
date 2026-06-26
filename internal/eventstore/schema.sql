CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (aggregate_id, version)
);

CREATE INDEX IF NOT EXISTS idx_events_aggregate_id ON events (aggregate_id);
