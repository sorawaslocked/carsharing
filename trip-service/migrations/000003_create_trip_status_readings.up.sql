CREATE TABLE IF NOT EXISTS trip_status_readings (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id     UUID        NOT NULL REFERENCES trips (id) ON DELETE CASCADE,
    from_status TEXT        NOT NULL,
    to_status   TEXT        NOT NULL,
    actor_type  TEXT        NOT NULL,
    actor_id    UUID,
    reason      TEXT,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trip_status_readings_trip_id   ON trip_status_readings (trip_id);
CREATE INDEX IF NOT EXISTS idx_trip_status_readings_changed_at ON trip_status_readings (trip_id, changed_at);
