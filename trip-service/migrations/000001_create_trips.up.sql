CREATE TABLE IF NOT EXISTS trips (
    id                   UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id           UUID             NOT NULL,
    user_id              UUID             NOT NULL,
    car_id               UUID             NOT NULL,
    status               TEXT             NOT NULL DEFAULT 'active',

    started_at           TIMESTAMPTZ      NOT NULL,
    start_latitude       DOUBLE PRECISION NOT NULL,
    start_longitude      DOUBLE PRECISION NOT NULL,
    start_mileage_km     BIGINT           NOT NULL,
    start_fuel_level     REAL,

    ended_at             TIMESTAMPTZ,
    end_latitude         DOUBLE PRECISION,
    end_longitude        DOUBLE PRECISION,
    end_mileage_km       BIGINT,
    end_fuel_level       REAL,

    distance_traveled_km DOUBLE PRECISION,
    duration_seconds     BIGINT,
    final_cost_tenge     INTEGER,
    cancel_reason        TEXT,

    created_at           TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trips_booking_id ON trips (booking_id);
CREATE INDEX IF NOT EXISTS idx_trips_user_id    ON trips (user_id);
CREATE INDEX IF NOT EXISTS idx_trips_car_id     ON trips (car_id);
CREATE INDEX IF NOT EXISTS idx_trips_status     ON trips (status);
CREATE INDEX IF NOT EXISTS idx_trips_started_at ON trips (started_at);
