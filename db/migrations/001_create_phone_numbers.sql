-- +goose Up
CREATE TABLE phone_numbers (
    id            BIGSERIAL PRIMARY KEY,
    phone_number  VARCHAR(20) NOT NULL UNIQUE,
    source        VARCHAR(50) NOT NULL,
    country       VARCHAR(100) NOT NULL,
    region        VARCHAR(100),
    provider      VARCHAR(100),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS phone_numbers;
