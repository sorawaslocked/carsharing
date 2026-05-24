package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/adapter/postgres"
	"carsharing/trip-service/internal/model"
)

func TestTripStatusReadingRepo_Create(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	statusRepo := postgres.NewTripStatusReadingRepo(discardLogger(), pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "bbbbbbbb-cccc-4ccc-8ccc-cccccccccccc")
	actorID := pgTestUserID

	r, err := statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     trip.ID,
		FromStatus: model.TripStatus(""),
		ToStatus:   model.TripStatusActive,
		ActorType:  sharedmodel.ActorTypeUser,
		ActorID:    &actorID,
		ChangedAt:  time.Now(),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, trip.ID, r.TripID)
	assert.Equal(t, model.TripStatusActive, r.ToStatus)
}

func TestTripStatusReadingRepo_List(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	statusRepo := postgres.NewTripStatusReadingRepo(discardLogger(), pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "cccccccc-dddd-4ddd-8ddd-dddddddddddd")
	actorID := pgTestUserID
	now := time.Now()

	for _, toStatus := range []model.TripStatus{model.TripStatusActive, model.TripStatusCompleted} {
		_, err := statusRepo.Create(ctx, model.TripStatusReadingCreate{
			TripID:    trip.ID,
			ToStatus:  toStatus,
			ActorType: sharedmodel.ActorTypeUser,
			ActorID:   &actorID,
			ChangedAt: now,
		})
		require.NoError(t, err)
	}

	readings, err := statusRepo.List(ctx, model.TripStatusReadingFilter{
		TripID:     trip.ID,
		Pagination: &sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, readings, 2)
}

func TestTripStatusReadingRepo_List_Pagination(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	statusRepo := postgres.NewTripStatusReadingRepo(discardLogger(), pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "dddddddd-eeee-4eee-8eee-eeeeeeeeeeee")
	actorID := pgTestUserID

	for i := range 3 {
		_, err := statusRepo.Create(ctx, model.TripStatusReadingCreate{
			TripID:    trip.ID,
			ToStatus:  model.TripStatusActive,
			ActorType: sharedmodel.ActorTypeUser,
			ActorID:   &actorID,
			ChangedAt: time.Now().Add(time.Duration(i) * time.Second),
		})
		require.NoError(t, err)
	}

	readings, err := statusRepo.List(ctx, model.TripStatusReadingFilter{
		TripID:     trip.ID,
		Pagination: &sharedmodel.Pagination{Limit: 2, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, readings, 2)
}
