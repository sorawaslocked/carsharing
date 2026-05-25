//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"carsharing/trip-service/internal/model"
)

// --- Create ---

func TestTripSummaryRepo_Create_ReturnsSummary(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripSummaryRepo()

	trip := mustInsertTrip(t, testTripCreate())
	now := time.Now().Truncate(time.Millisecond)

	s, err := r.Create(ctx, model.TripSummaryCreate{
		TripID:             trip.ID,
		BookingID:          trip.BookingID,
		StartedAt:          trip.StartedAt,
		EndedAt:            now,
		DurationSeconds:    600,
		DistanceTraveledKM: 10.5,
		PricingSnapshot:    model.PricingSnapshot{RateTenge: 50},
		BaseCostTenge:      500,
		DistanceCostTenge:  105,
		OvertimeCostTenge:  0,
		TotalCostTenge:     605,
	})

	require.NoError(t, err)
	assert.Equal(t, trip.ID, s.TripID)
	assert.Equal(t, int64(600), s.DurationSeconds)
	assert.InDelta(t, 10.5, s.DistanceTraveledKM, 0.001)
	assert.Equal(t, int32(605), s.TotalCostTenge)
}

// --- GetByTripID ---

func TestTripSummaryRepo_GetByTripID_Found(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newTripSummaryRepo()

	trip := mustInsertTrip(t, testTripCreate())
	_, err := r.Create(ctx, model.TripSummaryCreate{
		TripID:             trip.ID,
		BookingID:          trip.BookingID,
		StartedAt:          trip.StartedAt,
		EndedAt:            time.Now(),
		DurationSeconds:    300,
		DistanceTraveledKM: 5,
		PricingSnapshot:    model.PricingSnapshot{RateTenge: 30},
		BaseCostTenge:      300,
		DistanceCostTenge:  50,
		OvertimeCostTenge:  0,
		TotalCostTenge:     350,
	})
	require.NoError(t, err)

	got, err := r.GetByTripID(ctx, trip.ID)

	require.NoError(t, err)
	assert.Equal(t, trip.ID, got.TripID)
	assert.Equal(t, int32(350), got.TotalCostTenge)
}

func TestTripSummaryRepo_GetByTripID_NotFound(t *testing.T) {
	truncate(t)

	_, err := newTripSummaryRepo().GetByTripID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- Transaction ---

func TestTripSummaryRepo_Create_InTx(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	trip := mustInsertTrip(t, testTripCreate())

	err := newTransactor().InTx(ctx, func(ctx context.Context) error {
		_, e := newTripSummaryRepo().Create(ctx, model.TripSummaryCreate{
			TripID:             trip.ID,
			BookingID:          trip.BookingID,
			StartedAt:          trip.StartedAt,
			EndedAt:            time.Now(),
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
	got, err := newTripSummaryRepo().GetByTripID(ctx, trip.ID)
	require.NoError(t, err)
	assert.Equal(t, int32(220), got.TotalCostTenge)
}
