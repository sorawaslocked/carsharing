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

func TestCarStatusReadingRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	t.Run("inserts successfully", func(t *testing.T) {
		entry := model.CarStatusReading{
			CarID:      carID,
			FromStatus: model.CarStatusAvailable,
			ToStatus:   model.CarStatusReserved,
			ActorType:  sharedmodel.ActorTypeUser,
			RecordedAt: time.Now().UTC(),
		}
		err := newCarStatusReadingRepo().Insert(ctx, entry)
		assert.NoError(t, err)
	})

	t.Run("inserts with metadata and actor ID", func(t *testing.T) {
		entry := model.CarStatusReading{
			CarID:      carID,
			FromStatus: model.CarStatusReserved,
			ToStatus:   model.CarStatusInUse,
			ActorType:  sharedmodel.ActorTypeUser,
			ActorID:    ptr("user-001"),
			Reason:     ptr("trip started"),
			Metadata:   map[string]any{"booking_id": "b-001"},
			RecordedAt: time.Now().UTC(),
		}
		err := newCarStatusReadingRepo().Insert(ctx, entry)
		assert.NoError(t, err)
	})
}

func TestCarStatusReadingRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	repo := newCarStatusReadingRepo()

	base := time.Now().UTC()

	entries := []model.CarStatusReading{
		{CarID: carID, FromStatus: model.CarStatusAvailable, ToStatus: model.CarStatusReserved,
			ActorType: sharedmodel.ActorTypeUser, RecordedAt: base.Add(-2 * time.Hour)},
		{CarID: carID, FromStatus: model.CarStatusReserved, ToStatus: model.CarStatusInUse,
			ActorType: sharedmodel.ActorTypeUser, RecordedAt: base.Add(-1 * time.Hour)},
		{CarID: carID, FromStatus: model.CarStatusInUse, ToStatus: model.CarStatusAvailable,
			ActorType: sharedmodel.ActorTypeSystem, RecordedAt: base},
	}
	for _, e := range entries {
		require.NoError(t, repo.Insert(ctx, e))
	}

	t.Run("returns all readings for car", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarStatusReadingFilter{CarID: carID})
		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("filters by fromStatus", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarStatusReadingFilter{
			CarID:      carID,
			FromStatus: ptr(model.CarStatusReserved),
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, model.CarStatusReserved, results[0].FromStatus)
	})

	t.Run("filters by toStatus", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarStatusReadingFilter{
			CarID:    carID,
			ToStatus: ptr(model.CarStatusAvailable),
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, model.CarStatusAvailable, results[0].ToStatus)
	})

	t.Run("filters by time range", func(t *testing.T) {
		from := base.Add(-90 * time.Minute)
		to := base.Add(1 * time.Minute)
		results, err := repo.Find(ctx, model.CarStatusReadingFilter{
			CarID:     carID,
			TimeRange: &sharedmodel.TimeRange{From: from, To: to},
		})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarStatusReadingFilter{
			CarID:      carID,
			Pagination: &sharedmodel.Pagination{Limit: 2, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})
}
