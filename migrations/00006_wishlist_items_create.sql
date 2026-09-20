-- +goose Up
CREATE TABLE wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT,
    notes TEXT,
    price_cents INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT wishlist_items_price_check CHECK (
        price_cents IS NULL OR price_cents >= 0
    )
);

CREATE INDEX wishlist_items_user_id_idx ON wishlist_items (user_id);

-- +goose Down
DROP TABLE IF EXISTS wishlist_items;