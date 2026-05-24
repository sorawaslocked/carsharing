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

func TestZoneRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("returns generated ID", func(t *testing.T) {
		id, err := newZoneRepo().Insert(ctx, testZone())
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("data round-trips correctly", func(t *testing.T) {
		z := testZone()
		z.Name = "Parking Hub North"
		z.Type = model.ZoneParkingHub
		z.FeeAdjustment = 500
		z.IsActive = false

		id, err := newZoneRepo().Insert(ctx, z)
		require.NoError(t, err)

		got, err := newZoneRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, z.Name, got.Name)
		assert.Equal(t, z.Type, got.Type)
		assert.Equal(t, z.FeeAdjustment, got.FeeAdjustment)
		assert.Equal(t, z.IsActive, got.IsActive)
		assert.Equal(t, z.BoundaryGeoJSON, got.BoundaryGeoJSON)
	})
}

func TestZoneRepository_FindByID(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("returns zone when found", func(t *testing.T) {
		id := mustInsertZone(t)
		got, err := newZoneRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "Downtown", got.Name)
	})

	t.Run("returns ErrZoneNotFound for unknown ID", func(t *testing.T) {
		_, err := newZoneRepo().FindByID(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrZoneNotFound)
	})
}

func TestZoneRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	repo := newZoneRepo()

	z1 := testZone()
	z1.Name = "Operating Zone"
	z1.Type = model.ZoneTypeOperating
	z1.IsActive = true
	_, err := repo.Insert(ctx, z1)
	require.NoError(t, err)

	z2 := testZone()
	z2.Name = "No Drop Zone"
	z2.Type = model.ZoneTypeNoDrop
	z2.IsActive = false
	_, err = repo.Insert(ctx, z2)
	require.NoError(t, err)

	t.Run("returns all zones when no filter", func(t *testing.T) {
		results, err := repo.Find(ctx, model.ZoneFilter{})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by type", func(t *testing.T) {
		results, err := repo.Find(ctx, model.ZoneFilter{Type: ptr(model.ZoneTypeNoDrop)})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "No Drop Zone", results[0].Name)
	})

	t.Run("filters by isActive", func(t *testing.T) {
		results, err := repo.Find(ctx, model.ZoneFilter{IsActive: ptr(true)})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.True(t, results[0].IsActive)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.ZoneFilter{
			Pagination: &sharedmodel.Pagination{Limit: 1, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestZoneRepository_Update(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("updates name and isActive", func(t *testing.T) {
		id := mustInsertZone(t)
		err := newZoneRepo().Update(ctx, id, model.ZoneUpdate{
			Name:     ptr("Updated Zone"),
			IsActive: ptr(false),
		})
		require.NoError(t, err)

		got, err := newZoneRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Zone", got.Name)
		assert.False(t, got.IsActive)
	})

	t.Run("returns ErrZoneNotFound for unknown ID", func(t *testing.T) {
		err := newZoneRepo().Update(ctx, "00000000-0000-4000-8000-000000000000",
			model.ZoneUpdate{Name: ptr("X")})
		assert.ErrorIs(t, err, model.ErrZoneNotFound)
	})
}

func TestZoneRepository_Delete(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("deletes existing zone", func(t *testing.T) {
		id := mustInsertZone(t)
		err := newZoneRepo().Delete(ctx, id)
		require.NoError(t, err)

		_, err = newZoneRepo().FindByID(ctx, id)
		assert.ErrorIs(t, err, model.ErrZoneNotFound)
	})

	t.Run("returns ErrZoneNotFound for unknown ID", func(t *testing.T) {
		err := newZoneRepo().Delete(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrZoneNotFound)
	})
}
