CREATE EXTENSION IF NOT EXISTS "earthdistance" CASCADE;

CREATE TYPE car_status AS ENUM ('available', 'reserved', 'in_use', 'maintenance', 'out_of_service');

CREATE TABLE cars (
      id                UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
      model_id          UUID             NOT NULL REFERENCES car_models (id),
      vin               CHAR(17)         NOT NULL UNIQUE,
      license_plate     VARCHAR(20)      NOT NULL UNIQUE,
      color             VARCHAR(50)      NOT NULL,
      year_manufactured SMALLINT         NOT NULL,
      status            car_status       NOT NULL,
      mileage_km        BIGINT           NOT NULL,
      fuel_level        NUMERIC(5, 2),
      battery_level     NUMERIC(5, 2),
      latitude          DOUBLE PRECISION NOT NULL,
      longitude         DOUBLE PRECISION NOT NULL,
      notes             TEXT,
      image_keys        TEXT[]           NOT NULL DEFAULT '{}',
      last_seen_at      TIMESTAMPTZ      NOT NULL,
      created_at        TIMESTAMPTZ      NOT NULL,
      updated_at        TIMESTAMPTZ      NOT NULL
);

CREATE INDEX idx_cars_model_id ON cars (model_id);
CREATE INDEX idx_cars_status   ON cars (status);
CREATE INDEX idx_cars_location ON cars USING GIST (ll_to_earth(latitude, longitude));
