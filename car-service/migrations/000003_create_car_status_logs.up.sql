CREATE TABLE car_status_readings (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id      UUID        NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    from_status car_status  NOT NULL,
    to_status   car_status  NOT NULL,
    actor_type  VARCHAR(50) NOT NULL,
    actor_id    VARCHAR(100),
    reason      TEXT,
    metadata    JSONB,
    recorded_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_car_status_readings_car_id     ON car_status_readings (car_id);
CREATE INDEX idx_car_status_readings_recorded_at ON car_status_readings (recorded_at DESC);
