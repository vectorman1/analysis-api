CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS analysis;
CREATE SCHEMA IF NOT EXISTS "user";

DO
$do$
BEGIN
        IF NOT EXISTS (
                SELECT FROM pg_catalog.pg_user
                WHERE usename = 'harb') THEN
            CREATE USER harb WITH ENCRYPTED PASSWORD 'HueHue123';
END IF;
END
$do$;

GRANT ALL PRIVILEGES ON DATABASE analysis TO harb;
GRANT ALL PRIVILEGES ON SCHEMA analysis TO harb;
GRANT ALL PRIVILEGES ON SCHEMA "user" TO harb;

CREATE TABLE IF NOT EXISTS analysis.symbols
(
    id SERIAL PRIMARY KEY,
    uuid uuid UNIQUE NOT NULL,
    currency_code TEXT NOT NULL,
    isin TEXT NOT NULL,
    identifier TEXT NOT NULL,
    name TEXT NOT NULL,
    minimum_order_quantity REAL NOT NULL,
    market_name TEXT NOT NULL,
    market_hours_gmt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
    );

CREATE TABLE IF NOT EXISTS "user".users
(
    id SERIAL PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    private_role BIGINT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
    );
