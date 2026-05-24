package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		os.Exit(m.Run())
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic("connecting to test db: " + err.Error())
	}
	defer pool.Close()

	if _, err = pool.Exec(context.Background(), testSchema); err != nil {
		panic("applying schema: " + err.Error())
	}

	testPool = pool
	os.Exit(m.Run())
}

func requireDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		t.Skip("TEST_DATABASE_URL not set")
	}
	_, err := testPool.Exec(context.Background(),
		`TRUNCATE trip_status_readings, trip_summaries, trips CASCADE`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return testPool
}

const testSchema = `
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

CREATE TABLE IF NOT EXISTS trip_summaries (
	trip_id              UUID             PRIMARY KEY REFERENCES trips (id) ON DELETE CASCADE,
	booking_id           UUID             NOT NULL,
	started_at           TIMESTAMPTZ      NOT NULL,
	ended_at             TIMESTAMPTZ      NOT NULL,
	duration_seconds     BIGINT           NOT NULL,
	distance_traveled_km DOUBLE PRECISION NOT NULL,
	pricing_snapshot     JSONB            NOT NULL,
	base_cost_tenge      INTEGER          NOT NULL,
	distance_cost_tenge  INTEGER          NOT NULL,
	overtime_cost_tenge  INTEGER          NOT NULL,
	total_cost_tenge     INTEGER          NOT NULL
);

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
`
