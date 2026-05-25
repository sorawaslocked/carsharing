//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	pgadapter "carsharing/booking-service/internal/adapter/postgres"
	"carsharing/booking-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	pgc, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres container: %v\n", err)
		return 1
	}
	defer pgc.Terminate(ctx) //nolint:errcheck

	connStr, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get connection string: %v\n", err)
		return 1
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pool: %v\n", err)
		return 1
	}
	defer testPool.Close()

	for i := range 10 {
		if err = testPool.Ping(ctx); err == nil {
			break
		}
		if i == 9 {
			fmt.Fprintf(os.Stderr, "ping postgres: %v\n", err)
			return 1
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err := applyMigrations(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "apply migrations: %v\n", err)
		return 1
	}

	return m.Run()
}

func applyMigrations(ctx context.Context) error {
	entries, err := os.ReadDir("../../../migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, filepath.Join("../../../migrations", e.Name()))
		}
	}
	sort.Strings(upFiles)

	for _, path := range upFiles {
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := testPool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("exec %s: %w", path, err)
		}
	}

	return nil
}

// truncate clears all tables between tests. TRUNCATE pricing_rules CASCADE covers bookings and booking_status_history.
func truncate(t *testing.T) {
	t.Helper()
	if _, err := testPool.Exec(context.Background(), "TRUNCATE pricing_rules CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func discardLog() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newBookingRepo() *pgadapter.BookingRepository {
	return pgadapter.NewBookingRepository(discardLog(), testPool)
}

func newPricingRuleRepo() *pgadapter.PricingRuleRepository {
	return pgadapter.NewPricingRuleRepository(discardLog(), testPool)
}

func ptr[T any](v T) *T { return &v }

// testPricingRuleCreate returns a minimal valid PricingRuleCreate (by_minute, 100 tenge/min).
func testPricingRuleCreate() model.PricingRuleCreate {
	return model.PricingRuleCreate{
		Type:      string(model.PricingRuleTypeByMinute),
		RateTenge: 100,
	}
}

// mustInsertPricingRule inserts a pricing rule and returns its generated ID.
func mustInsertPricingRule(t *testing.T, data model.PricingRuleCreate) string {
	t.Helper()
	id, err := newPricingRuleRepo().Create(context.Background(), data)
	if err != nil {
		t.Fatalf("mustInsertPricingRule: %v", err)
	}
	return id
}

// testBookingCreate returns a BookingCreate referencing the given pricing rule.
func testBookingCreate(pricingRuleID string) model.BookingCreate {
	return model.BookingCreate{
		UserID:        "00000000-0000-4000-8000-000000000001",
		CarID:         "00000000-0000-4000-8000-000000000002",
		PricingRuleID: pricingRuleID,
	}
}

// mustInsertBooking inserts a booking expiring 15 minutes from now and returns its generated ID.
func mustInsertBooking(t *testing.T, data model.BookingCreate) string {
	t.Helper()
	id, err := newBookingRepo().Create(context.Background(), data, time.Now().Add(15*time.Minute))
	if err != nil {
		t.Fatalf("mustInsertBooking: %v", err)
	}
	return id
}
