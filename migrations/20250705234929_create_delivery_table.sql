-- +goose Up
-- +goose StatementBegin
CREATE TABLE delivery (
    order_uid VARCHAR PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR,
    phone VARCHAR,
    zip VARCHAR,
    city VARCHAR,
    address VARCHAR,
    region VARCHAR,
    email VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS delivery;
-- +goose StatementEnd
