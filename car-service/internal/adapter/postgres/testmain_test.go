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

	pgadapter "carsharing/car-service/internal/adapter/postgres"
	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"

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

	pgc, err := tcpostgres.Run(ctx, "postgis/postgis:16-3.4",
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

// truncate clears all tables before each test.
func truncate(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		"TRUNCATE car_models, zones, car_maintenance_templates CASCADE")
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func discardLog() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func ptr[T any](v T) *T { return &v }

// --- repo constructors ---

func newCarModelRepo() *pgadapter.CarModelRepository {
	return pgadapter.NewCarModelRepository(discardLog(), testPool)
}

func newZoneRepo() *pgadapter.ZoneRepository {
	return pgadapter.NewZoneRepository(discardLog(), testPool)
}

func newCarRepo() *pgadapter.CarRepository {
	return pgadapter.NewCarRepository(discardLog(), testPool)
}

func newCarInsuranceRepo() *pgadapter.CarInsuranceRepository {
	return pgadapter.NewCarInsuranceRepository(discardLog(), testPool)
}

func newCarMaintenanceTemplateRepo() *pgadapter.CarMaintenanceTemplateRepository {
	return pgadapter.NewCarMaintenanceTemplateRepository(discardLog(), testPool)
}

func newCarMaintenanceRecordRepo() *pgadapter.CarMaintenanceRecordRepository {
	return pgadapter.NewCarMaintenanceRecordRepository(discardLog(), testPool)
}

func newCarServiceStateRepo() *pgadapter.CarServiceStateRepository {
	return pgadapter.NewCarServiceStateRepository(discardLog(), testPool)
}

func newCarStatusReadingRepo() *pgadapter.CarStatusReadingRepository {
	return pgadapter.NewCarStatusReadingRepository(discardLog(), testPool)
}

func newTelemetryReadingRepo() *pgadapter.TelemetryReadingRepository {
	return pgadapter.NewTelemetryReadingRepository(discardLog(), testPool)
}

// --- test factories ---

func testCarModel() model.CarModel {
	now := time.Now().UTC()
	return model.CarModel{
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2022,
		FuelType:     model.CarFuelTypePetrol,
		Transmission: model.CarTransmissionAuto,
		BodyType:     model.CarBodyTypeSedan,
		Class:        model.CarClassComfort,
		Seats:        5,
		RangeKM:      600,
		Features:     []string{"AC", "GPS"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func testCarModelWithBrand(brand string) model.CarModel {
	cm := testCarModel()
	cm.Brand = brand
	return cm
}

func testZone() model.Zone {
	now := time.Now().UTC()
	return model.Zone{
		Name:            "Downtown",
		Type:            model.ZoneTypeOperating,
		BoundaryGeoJSON: `{"type":"Polygon","coordinates":[[[0,0],[1,0],[1,1],[0,1],[0,0]]]}`,
		FeeAdjustment:   0,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func testCar(modelID string) model.Car {
	now := time.Now().UTC()
	return model.Car{
		ModelID:          modelID,
		VIN:              "1HGCM82633A123456",
		LicensePlate:     "ABC-001",
		Color:            "red",
		YearManufactured: 2022,
		Status:           model.CarStatusAvailable,
		MileageKM:        50_000,
		Location:         sharedmodel.Location{Latitude: 51.5074, Longitude: -0.1278},
		LastSeenAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func testCarWithVIN(modelID, vin, plate string) model.Car {
	c := testCar(modelID)
	c.VIN = vin
	c.LicensePlate = plate
	return c
}

func testInsurance(carID string) model.CarInsurance {
	now := time.Now().UTC()
	return model.CarInsurance{
		CarID:     carID,
		Type:      model.InsuranceTypeOSAGO,
		Provider:  "Acme Insurance",
		PolicyNum: "POL-001",
		StartsAt:  now,
		ExpiresAt: now.Add(365 * 24 * time.Hour),
		CostTenge: 50_000,
		Status:    model.InsuranceStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func testTemplate() model.CarMaintenanceTemplate {
	now := time.Now().UTC()
	return model.CarMaintenanceTemplate{
		Name:        "Oil Change",
		KmInterval:  ptr[int32](10_000),
		DayInterval: ptr[int32](180),
		IsMandatory: true,
		WarnPct:     80,
		PullPct:     100,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func testRecord(carID, templateID string) model.CarMaintenanceRecord {
	now := time.Now().UTC()
	return model.CarMaintenanceRecord{
		CarID:      carID,
		TemplateID: templateID,
		Status:     model.MaintenanceRecordStatusPending,
		OdometerAt: 50_000,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// --- must-insert helpers ---

func mustInsertCarModel(t *testing.T) string {
	t.Helper()
	id, err := newCarModelRepo().Insert(context.Background(), testCarModel())
	if err != nil {
		t.Fatalf("mustInsertCarModel: %v", err)
	}
	return id
}

func mustInsertZone(t *testing.T) string {
	t.Helper()
	id, err := newZoneRepo().Insert(context.Background(), testZone())
	if err != nil {
		t.Fatalf("mustInsertZone: %v", err)
	}
	return id
}

func mustInsertCar(t *testing.T, modelID string) string {
	t.Helper()
	id, err := newCarRepo().Insert(context.Background(), testCar(modelID))
	if err != nil {
		t.Fatalf("mustInsertCar: %v", err)
	}
	return id
}

func mustInsertTemplate(t *testing.T) string {
	t.Helper()
	id, err := newCarMaintenanceTemplateRepo().Insert(context.Background(), testTemplate())
	if err != nil {
		t.Fatalf("mustInsertTemplate: %v", err)
	}
	return id
}
