CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    entity_id uuid NOT NULL,
    actor_id uuid NOT NULL,
    actor_role VARCHAR(255) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP NOT NULL,
    context JSONB
)