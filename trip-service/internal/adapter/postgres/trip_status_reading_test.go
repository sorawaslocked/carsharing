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

func TestTripStatusReadingRepo_Create_ReturnsReading(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripStatusReadingRepo()

	trip := mustInsertTrip(t, testTripCreate())
	actorID := trip.UserID

	reading, err := r.Create(ctx, model.TripStatusReadingCreate{
		TripID:     trip.ID,
		FromStatus: model.TripStatus(""),
		ToStatus:   model.TripStatusActive,
		ActorType:  sharedmodel.ActorTypeUser,
		ActorID:    &actorID,
		ChangedAt:  time.Now(),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, reading.ID)
	assert.Equal(t, trip.ID, reading.TripID)
	assert.Equal(t, model.TripStatusActive, reading.ToStatus)
}

// --- List ---

func TestTripStatusReadingRepo_List_ReturnsAll(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripStatusReadingRepo()

	trip := mustInsertTrip(t, testTripCreate())
	actorID := trip.UserID
	now := time.Now()

	for _, toStatus := range []model.TripStatus{model.TripStatusActive, model.TripStatusCompleted} {
		_, err := r.Create(ctx, model.TripStatusReadingCreate{
			TripID:    trip.ID,
			ToStatus:  toStatus,
			ActorType: sharedmodel.ActorTypeUser,
			ActorID:   &actorID,
			ChangedAt: now,
		})
		require.NoError(t, err)
	}

	readings, err := r.List(ctx, model.TripStatusReadingFilter{
		TripID:     trip.ID,
		Pagination: &sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, readings, 2)
}

func TestTripStatusReadingRepo_List_Pagination(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripStatusReadingRepo()

	trip := mustInsertTrip(t, testTripCreate())
	actorID := trip.UserID

	for i := range 3 {
		_, err := r.Create(ctx, model.TripStatusReadingCreate{
			TripID:    trip.ID,
			ToStatus:  model.TripStatusActive,
			ActorType: sharedmodel.ActorTypeUser,
			ActorID:   &actorID,
			ChangedAt: time.Now().Add(time.Duration(i) * time.Second),
		})
		require.NoError(t, err)
	}

	readings, err := r.List(ctx, model.TripStatusReadingFilter{
		TripID:     trip.ID,
		Pagination: &sharedmodel.Pagination{Limit: 2, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, readings, 2)
}
