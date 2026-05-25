//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/model"
)

// --- Create ---

func TestTripRepo_Create_ReturnsTrip(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	data := testTripCreate()

	trip, err := newTripRepo().Create(ctx, data)

	require.NoError(t, err)
	assert.Equal(t, data.ID, trip.ID)
	assert.Equal(t, data.BookingID, trip.BookingID)
	assert.Equal(t, model.TripStatusActive, trip.Status)
	assert.False(t, trip.CreatedAt.IsZero())
}

// --- GetByID ---

func TestTripRepo_GetByID_Found(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripRepo()

	created := mustInsertTrip(t, testTripCreate())

	got, err := r.GetByID(ctx, created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, created.UserID, got.UserID)
}

func TestTripRepo_GetByID_NotFound(t *testing.T) {
	truncate(t)

	_, err := newTripRepo().GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- List ---

func TestTripRepo_List_FilterByUserID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripRepo()

	mustInsertTrip(t, testTripCreate())
	second := testTripCreate()
	second.ID = "00000000-0000-4000-8000-000000000011"
	mustInsertTrip(t, second)

	other := testTripCreate()
	other.ID = "00000000-0000-4000-8000-000000000099"
	other.UserID = "99999999-9999-4999-8999-999999999999"
	mustInsertTrip(t, other)

	userID := testTripCreate().UserID
	trips, err := r.List(ctx, model.TripFilter{
		UserID:     &userID,
		Pagination: &sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, trips, 2)
	for _, trip := range trips {
		assert.Equal(t, userID, trip.UserID)
	}
}

func TestTripRepo_List_FilterByTimeRange(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripRepo()

	base := time.Now().Truncate(time.Millisecond)

	old := testTripCreate()
	old.StartedAt = base.Add(-2 * time.Hour)
	mustInsertTrip(t, old)

	recent := testTripCreate()
	recent.ID = "00000000-0000-4000-8000-000000000011"
	recent.StartedAt = base
	mustInsertTrip(t, recent)

	from := base.Add(-time.Hour)
	to := base.Add(time.Hour)
	trips, err := r.List(ctx, model.TripFilter{
		TimeRange:  &sharedmodel.TimeRange{From: from, To: to},
		Pagination: &sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, trips, 1)
	assert.Equal(t, recent.ID, trips[0].ID)
}

// --- Update ---

func TestTripRepo_Update_UpdatesStatus(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripRepo()

	created := mustInsertTrip(t, testTripCreate())

	updated, err := r.Update(ctx, created.ID, model.TripUpdate{
		Status:    ptr(model.TripStatusCompleted),
		UpdatedAt: time.Now(),
	})

	require.NoError(t, err)
	assert.Equal(t, model.TripStatusCompleted, updated.Status)
}

func TestTripRepo_Update_Conflict(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripRepo()

	created := mustInsertTrip(t, testTripCreate())
	stale := created.UpdatedAt.Add(-time.Second)

	_, err := r.Update(ctx, created.ID, model.TripUpdate{
		Status:            ptr(model.TripStatusCompleted),
		UpdatedAt:         time.Now(),
		ExpectedUpdatedAt: &stale,
	})

	assert.ErrorIs(t, err, model.ErrConflict)
}

func TestTripRepo_Update_OptimisticLock_Success(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripRepo()

	created := mustInsertTrip(t, testTripCreate())

	_, err := r.Update(ctx, created.ID, model.TripUpdate{
		Status:            ptr(model.TripStatusCompleted),
		UpdatedAt:         time.Now(),
		ExpectedUpdatedAt: &created.UpdatedAt,
	})

	require.NoError(t, err)
}

// --- Transaction ---

func TestTripRepo_Transaction_Rollback(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	data := testTripCreate()
	actorID := data.UserID

	_ = newTransactor().InTx(ctx, func(ctx context.Context) error {
		_, err := newTripRepo().Create(ctx, data)
		if err != nil {
			return err
		}
		_, err = newTripStatusReadingRepo().Create(ctx, model.TripStatusReadingCreate{
			TripID:     data.ID,
			FromStatus: model.TripStatus(""),
			ToStatus:   model.TripStatusActive,
			ActorType:  sharedmodel.ActorTypeUser,
			ActorID:    &actorID,
			ChangedAt:  time.Now(),
		})
		if err != nil {
			return err
		}
		return context.Canceled
	})

	_, err := newTripRepo().GetByID(ctx, data.ID)
	assert.ErrorIs(t, err, model.ErrNotFound, "rolled-back trip must not exist")
}
