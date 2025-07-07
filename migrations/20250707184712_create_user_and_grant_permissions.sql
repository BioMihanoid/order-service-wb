-- +goose Up
-- +goose StatementBegin
DO $$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
            CREATE ROLE order_service_user WITH LOGIN PASSWORD 'postgres';
        END IF;
    END
$$;

GRANT CONNECT ON DATABASE order_service_wb_db TO order_service_user;
GRANT USAGE ON SCHEMA public TO order_service_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO order_service_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO order_service_user;
GRANT CREATE ON SCHEMA public TO order_service_user; -- Если нужно создавать таблицы
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM order_service_user;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM order_service_user;
REVOKE USAGE ON SCHEMA public FROM order_service_user;
REVOKE CONNECT ON DATABASE order_service_wb_db FROM order_service_user;

DROP ROLE IF EXISTS order_service_user;
-- +goose StatementEnd