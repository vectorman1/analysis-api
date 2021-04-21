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
    currencyCode TEXT NOT NULL,
    isin TEXT NOT NULL,
    identifier TEXT NOT NULL,
    name TEXT NOT NULL,
    minimumOrderQuantity REAL NOT NULL,
    marketName TEXT NOT NULL,
    marketHoursGmt TEXT NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT now(),
    updatedAt TIMESTAMPTZ NOT NULL DEFAULT now(),
    deletedAt TIMESTAMPTZ NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "user".users
(
    id SERIAL PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    privateRole BIGINT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT now(),
    updatedAt TIMESTAMPTZ NOT NULL DEFAULT now(),
    deletedAt TIMESTAMPTZ NULL DEFAULT NULL
    );

-- Insert admin account
INSERT INTO "user".users VALUES (1, uuid_generate_v4(), 1, 'admin', '$2a$10$DT3TWK7tRrfdGhxY0KS9hux3PutpaU.7z1UQmn6eitfroNzwaUDMe');
