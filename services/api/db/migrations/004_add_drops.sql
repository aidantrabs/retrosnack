-- +goose Up
-- +goose StatementBegin

CREATE TABLE drops (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          TEXT        NOT NULL,
    slug          TEXT        NOT NULL UNIQUE,
    description   TEXT        NOT NULL DEFAULT '',
    instagram_url TEXT        NOT NULL DEFAULT '',
    released_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE products
    ADD COLUMN drop_id UUID REFERENCES drops(id),
    ADD COLUMN notes   TEXT NOT NULL DEFAULT '';

ALTER TABLE products DROP CONSTRAINT products_condition_check;
ALTER TABLE products ADD CONSTRAINT products_condition_check
    CHECK (condition IN ('new', 'excellent', 'good', 'fair'));

CREATE INDEX idx_products_drop ON products(drop_id);
CREATE INDEX idx_drops_slug ON drops(slug);
CREATE INDEX idx_drops_released ON drops(released_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_drops_released;
DROP INDEX IF EXISTS idx_drops_slug;
DROP INDEX IF EXISTS idx_products_drop;

ALTER TABLE products DROP CONSTRAINT products_condition_check;
ALTER TABLE products ADD CONSTRAINT products_condition_check
    CHECK (condition IN ('excellent', 'good', 'fair'));

ALTER TABLE products
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS drop_id;

DROP TABLE IF EXISTS drops;

-- +goose StatementEnd
