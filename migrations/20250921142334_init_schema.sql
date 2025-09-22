-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(50),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(100),
    shardkey VARCHAR(10),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE delivery (
    delivery_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    phone VARCHAR(50),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE payment (
    payment_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(100),
    amount DECIMAL(10,2),
    payment_dt BIGINT,
    bank VARCHAR(100),
    delivery_cost DECIMAL(10,2),
    goods_total DECIMAL(10,2),
    custom_fee DECIMAL(10,2)
);

CREATE TABLE items (
    item_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    chrt_id INTEGER,
    track_number VARCHAR(255),
    price DECIMAL(10,2),
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INTEGER,
    size VARCHAR(50),
    total_price DECIMAL(10,2),
    nm_id INTEGER,
    brand VARCHAR(100),
    status INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
