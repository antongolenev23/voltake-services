CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    aggregate_type TEXT NOT NULL,
    aggregate_id UUID NOT NULL,

    event_type TEXT NOT NULL,

    payload JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_events_unprocessed
ON outbox_events(created_at)
WHERE processed_at IS NULL;