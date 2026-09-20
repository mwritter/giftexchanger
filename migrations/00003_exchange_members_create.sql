-- +goose Up
CREATE TABLE exchange_members (
    exchange_id UUID NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (exchange_id, user_id)
);

CREATE INDEX exchange_members_user_id_idx ON exchange_members (user_id);

-- +goose Down
DROP TABLE IF EXISTS exchange_members;