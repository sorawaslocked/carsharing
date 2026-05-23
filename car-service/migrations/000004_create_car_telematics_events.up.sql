CREATE TABLE car_telematics_events (
    id            UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id        UUID             NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    latitude      DOUBLE PRECISION NOT NULL DEFAULT 0,
    longitude     DOUBLE PRECISION NOT NULL DEFAULT 0,
    fuel_level    NUMERIC(5, 2),
    battery_level NUMERIC(5, 2),
    odometer_km   BIGINT           NOT NULL DEFAULT 0,
    actor_id      VARCHAR(100),
    actor_type    VARCHAR(50)      NOT NULL DEFAULT 'telemetry',
    metadata      JSONB,
    recorded_at   TIMESTAMPTZ      NOT NULL,
    received_at   TIMESTAMPTZ      NOT NULL
);

CREATE INDEX idx_car_telematics_events_car_id      ON car_telematics_events (car_id);
CREATE INDEX idx_car_telematics_events_recorded_at ON car_telematics_events (recorded_at DESC);
CREATE INDEX idx_car_telematics_events_car_id_recorded_at ON car_telematics_events (car_id, recorded_at DESC);
