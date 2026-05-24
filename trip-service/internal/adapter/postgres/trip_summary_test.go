package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"carsharing/trip-service/internal/adapter/postgres"
	"carsharing/trip-service/internal/model"
)

func TestTripSummaryRepo_Create(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	summaryRepo := postgres.NewTripSummaryRepo(discardLogger(), pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "77777777-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	now := time.Now().Truncate(time.Millisecond)

	s, err := summaryRepo.Create(ctx, model.TripSummaryCreate{
		TripID: trip.ID, BookingID: pgTestBookingID,
		StartedAt: trip.StartedAt, EndedAt: now,
		DurationSeconds: 600, DistanceTraveledKM: 10.5,
		PricingSnapshot:   model.PricingSnapshot{RateTenge: 50},
		BaseCostTenge:     500,
		DistanceCostTenge: 105,
		OvertimeCostTenge: 0,
		TotalCostTenge:    605,
	})

	require.NoError(t, err)
	assert.Equal(t, trip.ID, s.TripID)
	assert.Equal(t, int64(600), s.DurationSeconds)
	assert.InDelta(t, 10.5, s.DistanceTraveledKM, 0.001)
	assert.Equal(t, int32(605), s.TotalCostTenge)
}

func TestTripSummaryRepo_GetByTripID(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	summaryRepo := postgres.NewTripSummaryRepo(discardLogger(), pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "88888888-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	_, err := summaryRepo.Create(ctx, model.TripSummaryCreate{
		TripID: trip.ID, BookingID: pgTestBookingID,
		StartedAt: trip.StartedAt, EndedAt: time.Now(),
		DurationSeconds: 300, DistanceTraveledKM: 5,
		PricingSnapshot:   model.PricingSnapshot{RateTenge: 30},
		BaseCostTenge:     300,
		DistanceCostTenge: 50,
		OvertimeCostTenge: 0,
		TotalCostTenge:    350,
	})
	require.NoError(t, err)

	got, err := summaryRepo.GetByTripID(ctx, trip.ID)

	require.NoError(t, err)
	assert.Equal(t, trip.ID, got.TripID)
	assert.Equal(t, int32(350), got.TotalCostTenge)
}

func TestTripSummaryRepo_GetByTripID_NotFound(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	summaryRepo := postgres.NewTripSummaryRepo(discardLogger(), pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "99999999-aaaa-4aaa-8aaa-aaaaaaaaaaaa")

	_, err := summaryRepo.GetByTripID(ctx, trip.ID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripSummaryRepo_Create_InTransaction(t *testing.T) {
	pool := requireDB(t)
	tripRepo := postgres.NewTripRepo(discardLogger(), pool)
	summaryRepo := postgres.NewTripSummaryRepo(discardLogger(), pool)
	transactor := postgres.NewTransactor(pool)
	ctx := context.Background()

	trip := insertTrip(t, tripRepo, "aaaaaaaa-bbbb-4bbb-8bbb-bbbbbbbbbbbb")

	err := transactor.InTx(ctx, func(ctx context.Context) error {
		_, e := summaryRepo.Create(ctx, model.TripSummaryCreate{
			TripID: trip.ID, BookingID: pgTestBookingID,
			StartedAt: trip.StartedAt, EndedAt: time.Now(),
			DurationSeconds:    120,
			DistanceTraveledKM: 2,
			PricingSnapshot:    model.PricingSnapshot{RateTenge: 20},
			BaseCostTenge:      200,
			DistanceCostTenge:  20,
			OvertimeCostTenge:  0,
			TotalCostTenge:     220,
		})
		return e
	})

	require.NoError(t, err)
	got, err := summaryRepo.GetByTripID(ctx, trip.ID)
	require.NoError(t, err)
	assert.Equal(t, int32(220), got.TotalCostTenge)
}
