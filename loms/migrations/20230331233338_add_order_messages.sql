-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_messages (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id bigint NOT NULL,
    status_id integer NOT NULL,
    is_processed boolean NOT NULL DEFAULT FALSE,
    error text NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_messages;
-- +goose StatementEnd
