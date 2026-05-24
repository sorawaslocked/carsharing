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

func TestCarInsuranceRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	t.Run("returns generated ID", func(t *testing.T) {
		id, err := newCarInsuranceRepo().Insert(ctx, testInsurance(carID))
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("data round-trips correctly", func(t *testing.T) {
		ins := testInsurance(carID)
		ins.Type = model.InsuranceTypeKASKO
		ins.PolicyNum = "POL-KASKO-001"
		ins.CostTenge = 100_000

		id, err := newCarInsuranceRepo().Insert(ctx, ins)
		require.NoError(t, err)

		got, err := newCarInsuranceRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, carID, got.CarID)
		assert.Equal(t, model.InsuranceTypeKASKO, got.Type)
		assert.Equal(t, "POL-KASKO-001", got.PolicyNum)
		assert.Equal(t, int32(100_000), got.CostTenge)
		assert.Equal(t, model.InsuranceStatusActive, got.Status)
	})
}

func TestCarInsuranceRepository_FindByID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	t.Run("returns insurance when found", func(t *testing.T) {
		id, err := newCarInsuranceRepo().Insert(ctx, testInsurance(carID))
		require.NoError(t, err)

		got, err := newCarInsuranceRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, carID, got.CarID)
	})

	t.Run("returns ErrCarInsuranceNotFound for unknown ID", func(t *testing.T) {
		_, err := newCarInsuranceRepo().FindByID(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarInsuranceNotFound)
	})
}

func TestCarInsuranceRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	repo := newCarInsuranceRepo()

	active := testInsurance(carID)
	active.Status = model.InsuranceStatusActive
	_, err := repo.Insert(ctx, active)
	require.NoError(t, err)

	expired := testInsurance(carID)
	expired.PolicyNum = "POL-002"
	expired.Status = model.InsuranceStatusExpired
	_, err = repo.Insert(ctx, expired)
	require.NoError(t, err)

	t.Run("returns all for car", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarInsuranceFilter{CarID: &carID})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by status", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarInsuranceFilter{
			CarID:  &carID,
			Status: ptr(model.InsuranceStatusExpired),
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, model.InsuranceStatusExpired, results[0].Status)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarInsuranceFilter{
			CarID:      &carID,
			Pagination: &sharedmodel.Pagination{Limit: 1, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestCarInsuranceRepository_Update(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	t.Run("updates status and cost", func(t *testing.T) {
		id, err := newCarInsuranceRepo().Insert(ctx, testInsurance(carID))
		require.NoError(t, err)

		err = newCarInsuranceRepo().Update(ctx, id, model.CarInsuranceUpdate{
			Status:    ptr(model.InsuranceStatusCancelled),
			CostTenge: ptr[int32](75_000),
		})
		require.NoError(t, err)

		got, err := newCarInsuranceRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.InsuranceStatusCancelled, got.Status)
		assert.Equal(t, int32(75_000), got.CostTenge)
	})

	t.Run("updates expiry date", func(t *testing.T) {
		id, err := newCarInsuranceRepo().Insert(ctx, testInsurance(carID))
		require.NoError(t, err)

		newExpiry := time.Now().UTC().Add(730 * 24 * time.Hour)
		err = newCarInsuranceRepo().Update(ctx, id, model.CarInsuranceUpdate{ExpiresAt: &newExpiry})
		require.NoError(t, err)

		got, err := newCarInsuranceRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.WithinDuration(t, newExpiry, got.ExpiresAt, time.Second)
	})

	t.Run("returns ErrCarInsuranceNotFound for unknown ID", func(t *testing.T) {
		err := newCarInsuranceRepo().Update(ctx, "00000000-0000-4000-8000-000000000000",
			model.CarInsuranceUpdate{Status: ptr(model.InsuranceStatusExpired)})
		assert.ErrorIs(t, err, model.ErrCarInsuranceNotFound)
	})
}

func TestCarInsuranceRepository_Delete(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	t.Run("deletes existing insurance", func(t *testing.T) {
		id, err := newCarInsuranceRepo().Insert(ctx, testInsurance(carID))
		require.NoError(t, err)

		err = newCarInsuranceRepo().Delete(ctx, id)
		require.NoError(t, err)

		_, err = newCarInsuranceRepo().FindByID(ctx, id)
		assert.ErrorIs(t, err, model.ErrCarInsuranceNotFound)
	})

	t.Run("returns ErrCarInsuranceNotFound for unknown ID", func(t *testing.T) {
		err := newCarInsuranceRepo().Delete(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarInsuranceNotFound)
	})
}
