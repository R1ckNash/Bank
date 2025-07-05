CREATE TABLE IF NOT EXISTS auth_outbox
(
    id             SERIAL PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id   UUID        NOT NULL,
    type           VARCHAR(50) NOT NULL,
    payload        JSONB       NOT NULL,
    created_at     TIMESTAMP   NOT NULL DEFAULT NOW(),
    sent_at        TIMESTAMP
);
