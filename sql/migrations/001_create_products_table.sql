-- +goose Up

CREATE TABLE IF NOT EXISTS products (
    id bigserial PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL,
    price integer NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS products;

-- something in here