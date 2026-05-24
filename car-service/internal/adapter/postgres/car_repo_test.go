//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCarRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	repo := newCarRepo()

	t.Run("returns generated ID", func(t *testing.T) {
		id, err := repo.Insert(ctx, testCarWithVIN(modelID, "1HGCM82633A000001", "PLT-001"))
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("duplicate VIN returns ErrDuplicateVIN", func(t *testing.T) {
		_, err := repo.Insert(ctx, testCarWithVIN(modelID, "1HGCM82633A000002", "PLT-002"))
		require.NoError(t, err)

		_, err = repo.Insert(ctx, testCarWithVIN(modelID, "1HGCM82633A000002", "PLT-003"))
		assert.ErrorIs(t, err, model.ErrDuplicateVIN)
	})

	t.Run("duplicate license plate returns ErrDuplicateLicensePlate", func(t *testing.T) {
		_, err := repo.Insert(ctx, testCarWithVIN(modelID, "1HGCM82633A000004", "PLT-004"))
		require.NoError(t, err)

		_, err = repo.Insert(ctx, testCarWithVIN(modelID, "1HGCM82633A000005", "PLT-004"))
		assert.ErrorIs(t, err, model.ErrDuplicateLicensePlate)
	})
}

func TestCarRepository_FindByID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)

	t.Run("returns car when found", func(t *testing.T) {
		id := mustInsertCar(t, modelID)
		got, err := newCarRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, modelID, got.ModelID)
		assert.Equal(t, "1HGCM82633A123456", got.VIN)
		assert.Equal(t, model.CarStatusAvailable, got.Status)
		assert.Equal(t, int64(50_000), got.MileageKM)
	})

	t.Run("returns ErrCarNotFound for unknown ID", func(t *testing.T) {
		_, err := newCarRepo().FindByID(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})
}

func TestCarRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	repo := newCarRepo()

	available := testCarWithVIN(modelID, "1HGCM82633A000010", "PLT-010")
	available.Status = model.CarStatusAvailable
	_, err := repo.Insert(ctx, available)
	require.NoError(t, err)

	reserved := testCarWithVIN(modelID, "1HGCM82633A000011", "PLT-011")
	reserved.Status = model.CarStatusReserved
	_, err = repo.Insert(ctx, reserved)
	require.NoError(t, err)

	t.Run("returns all cars when no filter", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarFilter{})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by status", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarFilter{Status: ptr(model.CarStatusReserved)})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, model.CarStatusReserved, results[0].Status)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarFilter{
			Pagination: &sharedmodel.Pagination{Limit: 1, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestCarRepository_Update(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	id := mustInsertCar(t, modelID)

	t.Run("updates mileage and status", func(t *testing.T) {
		err := newCarRepo().Update(ctx, id, model.CarUpdate{
			MileageKM: ptr[int64](60_000),
			Status:    ptr(model.CarStatusReserved),
		})
		require.NoError(t, err)

		got, err := newCarRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, int64(60_000), got.MileageKM)
		assert.Equal(t, model.CarStatusReserved, got.Status)
	})

	t.Run("updates location", func(t *testing.T) {
		loc := sharedmodel.Location{Latitude: 48.8566, Longitude: 2.3522}
		err := newCarRepo().Update(ctx, id, model.CarUpdate{Location: &loc})
		require.NoError(t, err)

		got, err := newCarRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.InDelta(t, loc.Latitude, got.Location.Latitude, 0.0001)
		assert.InDelta(t, loc.Longitude, got.Location.Longitude, 0.0001)
	})

	t.Run("returns ErrCarNotFound for unknown ID", func(t *testing.T) {
		err := newCarRepo().Update(ctx, "00000000-0000-4000-8000-000000000000",
			model.CarUpdate{MileageKM: ptr[int64](1)})
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})
}

func TestCarRepository_Delete(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)

	t.Run("deletes existing car", func(t *testing.T) {
		id := mustInsertCar(t, modelID)
		err := newCarRepo().Delete(ctx, id)
		require.NoError(t, err)

		_, err = newCarRepo().FindByID(ctx, id)
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})

	t.Run("returns ErrCarNotFound for unknown ID", func(t *testing.T) {
		err := newCarRepo().Delete(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})
}
