CREATE TABLE IF NOT EXISTS bookings (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID        NOT NULL,
    car_id            UUID        NOT NULL,
    committed_periods INT,
    status            TEXT        NOT NULL DEFAULT 'created',
    pricing_rule_id   UUID        NOT NULL REFERENCES pricing_rules(id),
    pricing_snapshot  JSONB       NOT NULL,
    expires_at        TIMESTAMPTZ NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS booking_status_history (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id  UUID        NOT NULL REFERENCES bookings(id),
    from_status TEXT        NOT NULL,
    to_status   TEXT        NOT NULL,
    actor_type  TEXT        NOT NULL,
    actor_id    UUID,
    reason      TEXT,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id       ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_car_id        ON bookings(car_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status        ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_expires_at    ON bookings(expires_at) WHERE status = 'created';
CREATE INDEX IF NOT EXISTS idx_status_history_booking ON booking_status_history(booking_id);
