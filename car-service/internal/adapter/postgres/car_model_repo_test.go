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

func TestCarModelRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("returns generated ID", func(t *testing.T) {
		id, err := newCarModelRepo().Insert(ctx, testCarModel())
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("data round-trips correctly", func(t *testing.T) {
		cm := testCarModel()
		cm.Brand = "Honda"
		cm.FuelType = model.CarFuelTypeElectric
		cm.EngineVolume = func() *float32 { v := float32(2.0); return &v }()

		id, err := newCarModelRepo().Insert(ctx, cm)
		require.NoError(t, err)

		got, err := newCarModelRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, cm.Brand, got.Brand)
		assert.Equal(t, cm.Model, got.Model)
		assert.Equal(t, cm.FuelType, got.FuelType)
		assert.Equal(t, cm.Class, got.Class)
		assert.Equal(t, cm.Features, got.Features)
		assert.NotNil(t, got.EngineVolume)
		assert.InDelta(t, *cm.EngineVolume, *got.EngineVolume, 0.01)
	})
}

func TestCarModelRepository_FindByID(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("returns model when found", func(t *testing.T) {
		id := mustInsertCarModel(t)
		got, err := newCarModelRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "Toyota", got.Brand)
	})

	t.Run("returns ErrCarModelNotFound for unknown ID", func(t *testing.T) {
		_, err := newCarModelRepo().FindByID(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarModelNotFound)
	})
}

func TestCarModelRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	repo := newCarModelRepo()

	_, err := repo.Insert(ctx, testCarModelWithBrand("Toyota"))
	require.NoError(t, err)
	_, err = repo.Insert(ctx, testCarModelWithBrand("Honda"))
	require.NoError(t, err)

	t.Run("returns all models when no filter", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarModelFilter{})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by brand", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarModelFilter{Brand: ptr("Toyota")})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Toyota", results[0].Brand)
	})

	t.Run("filters by class", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarModelFilter{Class: ptr(model.CarClassComfort)})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarModelFilter{
			Pagination: &sharedmodel.Pagination{Limit: 1, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestCarModelRepository_Update(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("updates brand field", func(t *testing.T) {
		id := mustInsertCarModel(t)
		err := newCarModelRepo().Update(ctx, id, model.CarModelUpdate{Brand: ptr("Nissan")})
		require.NoError(t, err)

		got, err := newCarModelRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Nissan", got.Brand)
	})

	t.Run("returns ErrCarModelNotFound for unknown ID", func(t *testing.T) {
		err := newCarModelRepo().Update(ctx, "00000000-0000-4000-8000-000000000000",
			model.CarModelUpdate{Brand: ptr("X")})
		assert.ErrorIs(t, err, model.ErrCarModelNotFound)
	})
}

func TestCarModelRepository_Delete(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("deletes existing model", func(t *testing.T) {
		id := mustInsertCarModel(t)
		err := newCarModelRepo().Delete(ctx, id)
		require.NoError(t, err)

		_, err = newCarModelRepo().FindByID(ctx, id)
		assert.ErrorIs(t, err, model.ErrCarModelNotFound)
	})

	t.Run("returns ErrCarModelNotFound for unknown ID", func(t *testing.T) {
		err := newCarModelRepo().Delete(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarModelNotFound)
	})
}
