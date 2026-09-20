-- +goose Up
CREATE TABLE exchanges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id UUID NOT NULL REFERENCES users (id),
    name TEXT NOT NULL,
    exchange_date DATE,
    state TEXT NOT NULL DEFAULT 'draft',
    invite_token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT exchanges_state_check CHECK (
        state IN ('draft', 'open', 'active', 'completed')
    )
);

CREATE INDEX exchanges_organizer_id_idx ON exchanges (organizer_id);

-- +goose Down
DROP TABLE IF EXISTS exchanges;