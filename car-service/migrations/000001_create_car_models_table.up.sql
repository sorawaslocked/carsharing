CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE car_fuel_type AS ENUM ('petrol', 'diesel', 'electric', 'hybrid');
CREATE TYPE car_transmission AS ENUM ('manual', 'auto');
CREATE TYPE car_body_type AS ENUM ('sedan', 'hatchback', 'suv', 'crossover', 'minivan', 'coupe', 'convertible', 'pickup');
CREATE TYPE car_class AS ENUM ('economy', 'compact', 'comfort', 'business', 'luxury');

CREATE TABLE car_models (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand         VARCHAR(100)     NOT NULL,
    model         VARCHAR(100)     NOT NULL,
    year          SMALLINT         NOT NULL,
    fuel_type     car_fuel_type    NOT NULL,
    transmission  car_transmission NOT NULL,
    body_type     car_body_type    NOT NULL,
    class         car_class        NOT NULL,
    seats         SMALLINT         NOT NULL,
    engine_volume NUMERIC(4, 1),
    range_km      INTEGER          NOT NULL,
    features      TEXT[]           NOT NULL,
    image_keys    TEXT[]           NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ      NOT NULL,
    updated_at    TIMESTAMPTZ      NOT NULL
);

CREATE INDEX idx_car_models_brand       ON car_models (brand);
CREATE INDEX idx_car_models_fuel_type   ON car_models (fuel_type);
CREATE INDEX idx_car_models_body_type   ON car_models (body_type);
CREATE INDEX idx_car_models_class       ON car_models (class);
CREATE INDEX idx_car_models_seats       ON car_models (seats);
