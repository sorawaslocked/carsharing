CREATE TABLE IF NOT EXISTS trip_summaries (
    trip_id              UUID             PRIMARY KEY REFERENCES trips (id) ON DELETE CASCADE,
    booking_id           UUID             NOT NULL,
    started_at           TIMESTAMPTZ      NOT NULL,
    ended_at             TIMESTAMPTZ      NOT NULL,
    duration_seconds     BIGINT           NOT NULL,
    distance_traveled_km DOUBLE PRECISION NOT NULL,
    pricing_snapshot     JSONB            NOT NULL,
    base_cost_tenge             INTEGER          NOT NULL,
    distance_cost_tenge         INTEGER          NOT NULL,
    overtime_cost_tenge         INTEGER          NOT NULL,
    zone_fee_adjustment_tenge   INTEGER          NOT NULL DEFAULT 0,
    total_cost_tenge            INTEGER          NOT NULL
);
