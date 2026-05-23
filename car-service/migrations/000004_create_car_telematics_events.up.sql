CREATE TABLE car_telemetry_readings (
    id            UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id        UUID             NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    latitude      DOUBLE PRECISION,
    longitude     DOUBLE PRECISION,
    fuel_pct      NUMERIC(5, 2),
    fuel_raw_pct  NUMERIC(5, 2),
    battery_level NUMERIC(5, 2),
    mileage_km    BIGINT,
    actor_id      VARCHAR(100),
    actor_type    VARCHAR(50)      NOT NULL DEFAULT 'telemetry',
    reason        TEXT,
    metadata      JSONB,
    recorded_at   TIMESTAMPTZ      NOT NULL
);

CREATE INDEX idx_car_telemetry_readings_car_id      ON car_telemetry_readings (car_id);
CREATE INDEX idx_car_telemetry_readings_recorded_at ON car_telemetry_readings (recorded_at DESC);
CREATE INDEX idx_car_telemetry_readings_car_id_recorded_at ON car_telemetry_readings (car_id, recorded_at DESC);
