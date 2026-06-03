CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE zone_type AS ENUM ('operating', 'no_drop', 'parking_hub', 'surcharge');

CREATE TABLE zones (
    id                UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(200) NOT NULL,
    type              zone_type NOT NULL,
    boundary_geo_json TEXT      NOT NULL,
    fee_adjustment    INT       NOT NULL DEFAULT 0,
    is_active         BOOLEAN   NOT NULL DEFAULT true,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_zones_type      ON zones (type);
CREATE INDEX idx_zones_is_active ON zones (is_active);
