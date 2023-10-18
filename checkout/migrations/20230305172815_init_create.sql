-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items (
    user_id bigint NOT NULL,
    sku int4 NOT NULL,
    count int4 NOT NULL DEFAULT 0,
    PRIMARY KEY(user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
