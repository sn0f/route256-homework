-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('new', 'awaiting payment', 'failed', 'payed', 'cancelled');

CREATE TABLE IF NOT EXISTS order_statuses (
    id integer NOT NULL PRIMARY KEY,
    name order_status NOT NULL
);

INSERT INTO order_statuses (id, name) VALUES (1, 'new') ON CONFLICT (id) DO NOTHING;
INSERT INTO order_statuses (id, name) VALUES (2, 'awaiting payment') ON CONFLICT (id) DO NOTHING;
INSERT INTO order_statuses (id, name) VALUES (3, 'failed') ON CONFLICT (id) DO NOTHING;
INSERT INTO order_statuses (id, name) VALUES (4, 'payed') ON CONFLICT (id) DO NOTHING;
INSERT INTO order_statuses (id, name) VALUES (5, 'cancelled') ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS orders (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id bigint NOT NULL,
    status_id integer NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id bigint NOT NULL,
    sku bigint NOT NULL,
    count bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS stocks (
    warehouse_id bigint NOT NULL,
    sku bigint NOT NULL,
    count bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (warehouse_id, sku)
);

CREATE TABLE IF NOT EXISTS reserves (
    order_id bigint NOT NULL,
    warehouse_id bigint NOT NULL,
    sku bigint NOT NULL,
    count bigint NOT NULL,
    PRIMARY KEY (order_id, warehouse_id, sku)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_statuses;
DROP TYPE IF EXISTS order_status;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS reserves;
DROP TABLE IF EXISTS stocks;
-- +goose StatementEnd
