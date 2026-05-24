package postgres_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/adapter/postgres"
	"carsharing/trip-service/internal/model"
)

const (
	pgTestUserID    = "11111111-1111-4111-8111-111111111111"
	pgTestBookingID = "22222222-2222-4222-8222-222222222222"
	pgTestCarID     = "44444444-4444-4444-8444-444444444444"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func ptrOf[T any](v T) *T { return &v }

func insertTrip(t *testing.T, repo *postgres.TripRepo, id string) model.Trip {
	t.Helper()
	trip, err := repo.Create(context.Background(), model.TripCreate{
		ID:             id,
		BookingID:      pgTestBookingID,
		UserID:         pgTestUserID,
		CarID:          pgTestCarID,
		Status:         model.TripStatusActive,
		StartedAt:      time.Now().Truncate(time.Millisecond),
		StartLocation:  sharedmodel.Location{Latitude: 51.5, Longitude: -0.1},
		StartMileageKM: 1000,
	})
	require.NoError(t, err)
	return trip
}

func TestTripRepo_Create(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	trip, err := repo.Create(ctx, model.TripCreate{
		ID:             "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		BookingID:      pgTestBookingID,
		UserID:         pgTestUserID,
		CarID:          pgTestCarID,
		Status:         model.TripStatusActive,
		StartedAt:      time.Now(),
		StartLocation:  sharedmodel.Location{Latitude: 51.5, Longitude: -0.1},
		StartMileageKM: 1000,
	})

	require.NoError(t, err)
	assert.Equal(t, "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", trip.ID)
	assert.Equal(t, pgTestBookingID, trip.BookingID)
	assert.Equal(t, model.TripStatusActive, trip.Status)
	assert.False(t, trip.CreatedAt.IsZero())
}

func TestTripRepo_GetByID(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	created := insertTrip(t, repo, "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb")

	got, err := repo.GetByID(ctx, created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, created.UserID, got.UserID)
}

func TestTripRepo_GetByID_NotFound(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "cccccccc-cccc-4ccc-8ccc-cccccccccccc")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripRepo_List_FilterByUserID(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	insertTrip(t, repo, "dddddddd-dddd-4ddd-8ddd-dddddddddddd")
	insertTrip(t, repo, "eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee")

	otherUserID := "99999999-9999-4999-8999-999999999999"
	_, err := repo.Create(ctx, model.TripCreate{
		ID: "ffffffff-ffff-4fff-8fff-ffffffffffff", BookingID: pgTestBookingID,
		UserID: otherUserID, CarID: pgTestCarID,
		Status: model.TripStatusActive, StartedAt: time.Now(),
		StartLocation: sharedmodel.Location{}, StartMileageKM: 0,
	})
	require.NoError(t, err)

	userID := pgTestUserID
	trips, err := repo.List(ctx, model.TripFilter{
		UserID:     &userID,
		Pagination: &sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, trips, 2)
	for _, trip := range trips {
		assert.Equal(t, pgTestUserID, trip.UserID)
	}
}

func TestTripRepo_List_FilterByTimeRange(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	base := time.Now().Truncate(time.Millisecond)

	_, err := repo.Create(ctx, model.TripCreate{
		ID: "11111111-aaaa-4aaa-8aaa-aaaaaaaaaaaa", BookingID: pgTestBookingID,
		UserID: pgTestUserID, CarID: pgTestCarID,
		Status: model.TripStatusActive, StartedAt: base.Add(-2 * time.Hour),
		StartLocation: sharedmodel.Location{}, StartMileageKM: 0,
	})
	require.NoError(t, err)
	_, err = repo.Create(ctx, model.TripCreate{
		ID: "22222222-aaaa-4aaa-8aaa-aaaaaaaaaaaa", BookingID: pgTestBookingID,
		UserID: pgTestUserID, CarID: pgTestCarID,
		Status: model.TripStatusActive, StartedAt: base,
		StartLocation: sharedmodel.Location{}, StartMileageKM: 0,
	})
	require.NoError(t, err)

	from := base.Add(-time.Hour)
	to := base.Add(time.Hour)
	trips, err := repo.List(ctx, model.TripFilter{
		TimeRange:  &sharedmodel.TimeRange{From: from, To: to},
		Pagination: &sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, trips, 1)
	assert.Equal(t, "22222222-aaaa-4aaa-8aaa-aaaaaaaaaaaa", trips[0].ID)
}

func TestTripRepo_Update(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	created := insertTrip(t, repo, "33333333-aaaa-4aaa-8aaa-aaaaaaaaaaaa")

	updated, err := repo.Update(ctx, created.ID, model.TripUpdate{
		Status:    ptrOf(model.TripStatusCompleted),
		UpdatedAt: time.Now(),
	})

	require.NoError(t, err)
	assert.Equal(t, model.TripStatusCompleted, updated.Status)
}

func TestTripRepo_Update_Conflict(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	created := insertTrip(t, repo, "44444444-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	stale := created.UpdatedAt.Add(-time.Second)

	_, err := repo.Update(ctx, created.ID, model.TripUpdate{
		Status:            ptrOf(model.TripStatusCompleted),
		UpdatedAt:         time.Now(),
		ExpectedUpdatedAt: &stale,
	})

	assert.ErrorIs(t, err, model.ErrConflict)
}

func TestTripRepo_Update_OptimisticLock_Success(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	ctx := context.Background()

	created := insertTrip(t, repo, "55555555-aaaa-4aaa-8aaa-aaaaaaaaaaaa")

	_, err := repo.Update(ctx, created.ID, model.TripUpdate{
		Status:            ptrOf(model.TripStatusCompleted),
		UpdatedAt:         time.Now(),
		ExpectedUpdatedAt: &created.UpdatedAt,
	})

	require.NoError(t, err)
}

func TestTripRepo_Transaction_Rollback(t *testing.T) {
	pool := requireDB(t)
	repo := postgres.NewTripRepo(discardLogger(), pool)
	statusRepo := postgres.NewTripStatusReadingRepo(discardLogger(), pool)
	transactor := postgres.NewTransactor(pool)
	ctx := context.Background()

	tripID := "66666666-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	actorID := pgTestUserID

	_ = transactor.InTx(ctx, func(ctx context.Context) error {
		_, err := repo.Create(ctx, model.TripCreate{
			ID: tripID, BookingID: pgTestBookingID, UserID: pgTestUserID, CarID: pgTestCarID,
			Status: model.TripStatusActive, StartedAt: time.Now(),
			StartLocation: sharedmodel.Location{}, StartMileageKM: 0,
		})
		if err != nil {
			return err
		}
		_, err = statusRepo.Create(ctx, model.TripStatusReadingCreate{
			TripID: tripID, FromStatus: model.TripStatus(""), ToStatus: model.TripStatusActive,
			ActorType: sharedmodel.ActorTypeUser, ActorID: &actorID, ChangedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		return context.Canceled // force rollback
	})

	_, err := repo.GetByID(ctx, tripID)
	assert.ErrorIs(t, err, model.ErrNotFound, "rolled-back trip must not exist")
}
