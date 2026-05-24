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

// --- CarMaintenanceTemplateRepository ---

func TestCarMaintenanceTemplateRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("returns generated ID", func(t *testing.T) {
		id, err := newCarMaintenanceTemplateRepo().Insert(ctx, testTemplate())
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("data round-trips correctly", func(t *testing.T) {
		tmpl := testTemplate()
		tmpl.Name = "Tire Rotation"
		tmpl.KmInterval = ptr[int32](20_000)
		tmpl.DayInterval = nil
		tmpl.IsMandatory = false

		id, err := newCarMaintenanceTemplateRepo().Insert(ctx, tmpl)
		require.NoError(t, err)

		got, err := newCarMaintenanceTemplateRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Tire Rotation", got.Name)
		assert.Equal(t, int32(20_000), *got.KmInterval)
		assert.Nil(t, got.DayInterval)
		assert.False(t, got.IsMandatory)
	})
}

func TestCarMaintenanceTemplateRepository_FindByID(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("returns template when found", func(t *testing.T) {
		id := mustInsertTemplate(t)
		got, err := newCarMaintenanceTemplateRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "Oil Change", got.Name)
	})

	t.Run("returns ErrCarMaintenanceTemplateNotFound for unknown ID", func(t *testing.T) {
		_, err := newCarMaintenanceTemplateRepo().FindByID(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

func TestCarMaintenanceTemplateRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	repo := newCarMaintenanceTemplateRepo()

	mandatory := testTemplate()
	mandatory.Name = "Mandatory Service"
	mandatory.IsMandatory = true
	_, err := repo.Insert(ctx, mandatory)
	require.NoError(t, err)

	optional := testTemplate()
	optional.Name = "Optional Checkup"
	optional.IsMandatory = false
	_, err = repo.Insert(ctx, optional)
	require.NoError(t, err)

	t.Run("returns all templates when no filter", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarMaintenanceTemplateFilter{})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by isMandatory", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarMaintenanceTemplateFilter{IsMandatory: ptr(true)})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Mandatory Service", results[0].Name)
	})

	t.Run("pagination limits results", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarMaintenanceTemplateFilter{
			Pagination: &sharedmodel.Pagination{Limit: 1, Offset: 0},
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestCarMaintenanceTemplateRepository_Update(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("updates name and km interval", func(t *testing.T) {
		id := mustInsertTemplate(t)
		err := newCarMaintenanceTemplateRepo().Update(ctx, id, model.CarMaintenanceTemplateUpdate{
			Name:       ptr("Updated Service"),
			KmInterval: ptr[int32](15_000),
		})
		require.NoError(t, err)

		got, err := newCarMaintenanceTemplateRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Service", got.Name)
		assert.Equal(t, int32(15_000), *got.KmInterval)
	})

	t.Run("returns ErrCarMaintenanceTemplateNotFound for unknown ID", func(t *testing.T) {
		err := newCarMaintenanceTemplateRepo().Update(ctx, "00000000-0000-4000-8000-000000000000",
			model.CarMaintenanceTemplateUpdate{Name: ptr("X")})
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

func TestCarMaintenanceTemplateRepository_Delete(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	t.Run("deletes existing template", func(t *testing.T) {
		id := mustInsertTemplate(t)
		err := newCarMaintenanceTemplateRepo().Delete(ctx, id)
		require.NoError(t, err)

		_, err = newCarMaintenanceTemplateRepo().FindByID(ctx, id)
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})

	t.Run("returns ErrCarMaintenanceTemplateNotFound for unknown ID", func(t *testing.T) {
		err := newCarMaintenanceTemplateRepo().Delete(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

// --- CarMaintenanceRecordRepository ---

func TestCarMaintenanceRecordRepository_Insert(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	templateID := mustInsertTemplate(t)

	t.Run("returns generated ID", func(t *testing.T) {
		id, err := newCarMaintenanceRecordRepo().Insert(ctx, testRecord(carID, templateID))
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("data round-trips correctly", func(t *testing.T) {
		rec := testRecord(carID, templateID)
		rec.AssignedTo = ptr("mechanic-01")
		rec.Notes = ptr("Scheduled service")
		due := time.Now().UTC().Add(7 * 24 * time.Hour)
		rec.DueBy = &due

		id, err := newCarMaintenanceRecordRepo().Insert(ctx, rec)
		require.NoError(t, err)

		got, err := newCarMaintenanceRecordRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, carID, got.CarID)
		assert.Equal(t, templateID, got.TemplateID)
		assert.Equal(t, model.MaintenanceRecordStatusPending, got.Status)
		assert.Equal(t, "mechanic-01", *got.AssignedTo)
		assert.Equal(t, "Scheduled service", *got.Notes)
		assert.NotNil(t, got.DueBy)
	})
}

func TestCarMaintenanceRecordRepository_FindByID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	templateID := mustInsertTemplate(t)

	t.Run("returns record when found", func(t *testing.T) {
		id, err := newCarMaintenanceRecordRepo().Insert(ctx, testRecord(carID, templateID))
		require.NoError(t, err)

		got, err := newCarMaintenanceRecordRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, carID, got.CarID)
	})

	t.Run("returns ErrCarMaintenanceRecordNotFound for unknown ID", func(t *testing.T) {
		_, err := newCarMaintenanceRecordRepo().FindByID(ctx, "00000000-0000-4000-8000-000000000000")
		assert.ErrorIs(t, err, model.ErrCarMaintenanceRecordNotFound)
	})
}

func TestCarMaintenanceRecordRepository_Find(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	templateID := mustInsertTemplate(t)
	repo := newCarMaintenanceRecordRepo()

	pending := testRecord(carID, templateID)
	pending.Status = model.MaintenanceRecordStatusPending
	_, err := repo.Insert(ctx, pending)
	require.NoError(t, err)

	inProgress := testRecord(carID, templateID)
	inProgress.Status = model.MaintenanceRecordStatusInProgress
	_, err = repo.Insert(ctx, inProgress)
	require.NoError(t, err)

	t.Run("returns all for car", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarMaintenanceRecordFilter{CarID: &carID})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by status", func(t *testing.T) {
		status := model.MaintenanceRecordStatusInProgress
		results, err := repo.Find(ctx, model.CarMaintenanceRecordFilter{
			CarID:  &carID,
			Status: &status,
		})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, model.MaintenanceRecordStatusInProgress, results[0].Status)
	})

	t.Run("filters by templateID", func(t *testing.T) {
		results, err := repo.Find(ctx, model.CarMaintenanceRecordFilter{TemplateID: &templateID})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})
}

func TestCarMaintenanceRecordRepository_Update(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	templateID := mustInsertTemplate(t)

	t.Run("updates status and completedKM", func(t *testing.T) {
		id, err := newCarMaintenanceRecordRepo().Insert(ctx, testRecord(carID, templateID))
		require.NoError(t, err)

		now := time.Now().UTC()
		err = newCarMaintenanceRecordRepo().Update(ctx, id, model.CarMaintenanceRecordUpdate{
			Status:      ptr(model.MaintenanceRecordStatusCompleted),
			CompletedKM: ptr[int32](55_000),
			CompletedAt: &now,
		})
		require.NoError(t, err)

		got, err := newCarMaintenanceRecordRepo().FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.MaintenanceRecordStatusCompleted, got.Status)
		assert.Equal(t, int32(55_000), *got.CompletedKM)
	})

	t.Run("returns ErrCarMaintenanceRecordNotFound for unknown ID", func(t *testing.T) {
		err := newCarMaintenanceRecordRepo().Update(ctx, "00000000-0000-4000-8000-000000000000",
			model.CarMaintenanceRecordUpdate{Status: ptr(model.MaintenanceRecordStatusCompleted)})
		assert.ErrorIs(t, err, model.ErrCarMaintenanceRecordNotFound)
	})
}

// --- CarServiceStateRepository ---

func TestCarServiceStateRepository_Upsert(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)
	templateID := mustInsertTemplate(t)
	repo := newCarServiceStateRepo()

	t.Run("inserts new state", func(t *testing.T) {
		state := model.CarServiceState{
			CarID:      carID,
			TemplateID: templateID,
			LastKM:     50_000,
			NextDueKM:  ptr[int32](60_000),
		}
		err := repo.Upsert(ctx, state)
		require.NoError(t, err)

		results, err := repo.FindAll(ctx, model.CarServiceStateFilter{CarID: &carID})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, int32(50_000), results[0].LastKM)
		assert.Equal(t, int32(60_000), *results[0].NextDueKM)
	})

	t.Run("updates existing state on conflict", func(t *testing.T) {
		updated := model.CarServiceState{
			CarID:      carID,
			TemplateID: templateID,
			LastKM:     55_000,
			NextDueKM:  ptr[int32](65_000),
		}
		err := repo.Upsert(ctx, updated)
		require.NoError(t, err)

		results, err := repo.FindAll(ctx, model.CarServiceStateFilter{CarID: &carID})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, int32(55_000), results[0].LastKM)
		assert.Equal(t, int32(65_000), *results[0].NextDueKM)
	})
}

func TestCarServiceStateRepository_FindAll(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	modelID := mustInsertCarModel(t)
	carID := mustInsertCar(t, modelID)

	tmpl1ID := mustInsertTemplate(t)

	tmpl2 := testTemplate()
	tmpl2.Name = "Brake Check"
	tmpl2ID, err := newCarMaintenanceTemplateRepo().Insert(ctx, tmpl2)
	require.NoError(t, err)

	repo := newCarServiceStateRepo()
	err = repo.Upsert(ctx, model.CarServiceState{CarID: carID, TemplateID: tmpl1ID, LastKM: 50_000})
	require.NoError(t, err)
	err = repo.Upsert(ctx, model.CarServiceState{CarID: carID, TemplateID: tmpl2ID, LastKM: 50_000})
	require.NoError(t, err)

	t.Run("returns all states for car", func(t *testing.T) {
		results, err := repo.FindAll(ctx, model.CarServiceStateFilter{CarID: &carID})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filters by templateID", func(t *testing.T) {
		results, err := repo.FindAll(ctx, model.CarServiceStateFilter{TemplateID: &tmpl1ID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, tmpl1ID, results[0].TemplateID)
	})
}
