-- +goose Up
CREATE TABLE magic_link_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- token_hash UNIQUE already indexes "find token from link".
CREATE INDEX magic_link_tokens_email_lower_idx ON magic_link_tokens (lower(email));
CREATE INDEX magic_link_tokens_expires_at_idx ON magic_link_tokens (expires_at);

-- +goose Down
DROP TABLE IF EXISTS magic_link_tokens;
