//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryReadingRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	t.Run("inserts minimal reading", func(t *testing.T) {
		reading := model.TelemetryReading{
			CarID:      carID,
			ActorType:  sharedmodel.ActorTypeTelemetry,
			RecordedAt: time.Now().UTC(),
		}
		err := newTelemetryReadingRepo().Insert(ctx, reading)
		assert.NoError(t, err)
	})

	t.Run("inserts reading with all fields", func(t *testing.T) {
		mileage := int64(51_000)
		fuel := float32(75.5)
		battery := float32(90.0)
		reading := model.TelemetryReading{
			CarID:        carID,
			Location:     &sharedmodel.Location{Latitude: 51.5074, Longitude: -0.1278},
			FuelPct:      &fuel,
			BatteryLevel: &battery,
			MileageKM:    &mileage,
			ActorType:    sharedmodel.ActorTypeTelemetry,
			RecordedAt:   time.Now().UTC(),
		}
		err := newTelemetryReadingRepo().Insert(ctx, reading)
		assert.NoError(t, err)
	})
}

func TestTelemetryReadingRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	repo := newTelemetryReadingRepo()

	base := time.Now().UTC()

	readings := []model.TelemetryReading{
		{CarID: carID, ActorType: sharedmodel.ActorTypeTelemetry, RecordedAt: base.Add(-2 * time.Hour)},
		{CarID: carID, ActorType: sharedmodel.ActorTypeTelemetry, RecordedAt: base.Add(-1 * time.Hour)},
		{CarID: carID, ActorType: sharedmodel.ActorTypeTelemetry, RecordedAt: base},
	}
	for _, r := range readings {
		require.NoError(t, repo.Insert(ctx, r))
	}

	t.Run("returns all readings for car", func(t *testing.T) {
		results, err := repo.Find(ctx, model.TelemetryReadingFilter{CarID: carID})
		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("filters by time range", func(t *testing.T) {
		from := base.Add(-90 * time.Minute)
		to := base.Add(1 * time.Minute)
		results, err := repo.Find(ctx, model.TelemetryReadingFilter{
			CarID:     carID,
			TimeRange: &sharedmodel.TimeRange{From: from, To: to},
		})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.TelemetryReadingFilter{
			CarID:      carID,
			Pagination: &sharedmodel.Pagination{Limit: 2, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("returns empty for unknown car", func(t *testing.T) {
		results, err := repo.Find(ctx, model.TelemetryReadingFilter{
			CarID: "00000000-0000-4000-8000-000000000000",
		})
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}
