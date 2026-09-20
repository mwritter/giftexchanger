-- +goose Up
CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exchange_id UUID NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    giver_id UUID NOT NULL REFERENCES users (id),
    recipient_id UUID NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (exchange_id, giver_id),
    UNIQUE (exchange_id, recipient_id),
    CONSTRAINT assignments_no_self CHECK (giver_id <> recipient_id),
    FOREIGN KEY (exchange_id, giver_id)
        REFERENCES exchange_members (exchange_id, user_id),
    FOREIGN KEY (exchange_id, recipient_id)
        REFERENCES exchange_members (exchange_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS assignments;
