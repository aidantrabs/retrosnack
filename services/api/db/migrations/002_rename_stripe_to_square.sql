-- +goose Up
-- +goose StatementBegin

ALTER TABLE orders RENAME COLUMN stripe_session_id TO checkout_session_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE orders RENAME COLUMN checkout_session_id TO stripe_session_id;

-- +goose StatementEnd
