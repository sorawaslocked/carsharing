package service

import (
	"context"
	"testing"
	"time"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/service/mocks"
	"carsharing/car-service/internal/validation"
	sharedvalidation "carsharing/shared/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestCarServiceHistory(t *testing.T, statusRepo CarStatusReadingRepository, telemetryRepo TelemetryReadingRepository) *CarService {
	t.Helper()
	return NewCarService(discardLogger(), newTestValidator(t), nil, nil, nil, statusRepo, telemetryRepo, nil, nil, noopCarCreatedNotifier{})
}

func TestListCarStatusHistory(t *testing.T) {
	carID := "c0000000-0000-4000-8000-000000000001"
	ctx := context.Background()

	t.Run("returns readings for valid filter", func(t *testing.T) {
		repo := mocks.NewMockCarStatusReadingRepository(t)
		svc := newTestCarServiceHistory(t, repo, nil)

		repo.EXPECT().
			Find(ctx, mock.MatchedBy(func(f model.CarStatusReadingFilter) bool {
				return f.CarID == carID && f.TimeRange == nil
			})).
			Return([]model.CarStatusReading{
				{CarID: carID, FromStatus: model.CarStatusAvailable, ToStatus: model.CarStatusReserved},
			}, nil)

		entries, err := svc.ListCarStatusHistory(ctx, validation.CarStatusReadingFilter{CarID: carID})
		assert.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Equal(t, model.CarStatusReserved, entries[0].ToStatus)
	})

	t.Run("time range is passed to repo", func(t *testing.T) {
		repo := mocks.NewMockCarStatusReadingRepository(t)
		svc := newTestCarServiceHistory(t, repo, nil)

		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		repo.EXPECT().
			Find(ctx, mock.MatchedBy(func(f model.CarStatusReadingFilter) bool {
				return f.TimeRange != nil &&
					f.TimeRange.From.Equal(from) &&
					f.TimeRange.To.Equal(to)
			})).
			Return([]model.CarStatusReading{}, nil)

		_, err := svc.ListCarStatusHistory(ctx, validation.CarStatusReadingFilter{
			CarID: carID,
			TimeRange: &sharedvalidation.TimeRange{
				From: &from,
				To:   &to,
			},
		})
		assert.NoError(t, err)
	})

	t.Run("from_status and to_status are passed to repo", func(t *testing.T) {
		repo := mocks.NewMockCarStatusReadingRepository(t)
		svc := newTestCarServiceHistory(t, repo, nil)

		fromStatus := string(model.CarStatusAvailable)
		toStatus := string(model.CarStatusReserved)

		repo.EXPECT().
			Find(ctx, mock.MatchedBy(func(f model.CarStatusReadingFilter) bool {
				return f.FromStatus != nil && *f.FromStatus == model.CarStatusAvailable &&
					f.ToStatus != nil && *f.ToStatus == model.CarStatusReserved
			})).
			Return([]model.CarStatusReading{}, nil)

		_, err := svc.ListCarStatusHistory(ctx, validation.CarStatusReadingFilter{
			CarID:      carID,
			FromStatus: &fromStatus,
			ToStatus:   &toStatus,
		})
		assert.NoError(t, err)
	})

	t.Run("invalid car_id returns validation error", func(t *testing.T) {
		svc := newTestCarServiceHistory(t, nil, nil)

		_, err := svc.ListCarStatusHistory(ctx, validation.CarStatusReadingFilter{CarID: "not-a-uuid"})
		assert.Error(t, err)
	})

	t.Run("reversed time range returns validation error", func(t *testing.T) {
		svc := newTestCarServiceHistory(t, nil, nil)

		from := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		_, err := svc.ListCarStatusHistory(ctx, validation.CarStatusReadingFilter{
			CarID: carID,
			TimeRange: &sharedvalidation.TimeRange{
				From: &from,
				To:   &to,
			},
		})
		assert.Error(t, err)
	})

	t.Run("repo error is returned", func(t *testing.T) {
		repo := mocks.NewMockCarStatusReadingRepository(t)
		svc := newTestCarServiceHistory(t, repo, nil)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrSql)

		_, err := svc.ListCarStatusHistory(ctx, validation.CarStatusReadingFilter{CarID: carID})
		assert.ErrorIs(t, err, model.ErrSql)
	})
}

func TestListCarTelemetryHistory(t *testing.T) {
	carID := "c0000000-0000-4000-8000-000000000001"
	ctx := context.Background()

	t.Run("returns readings for valid filter", func(t *testing.T) {
		repo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newTestCarServiceHistory(t, nil, repo)

		repo.EXPECT().
			Find(ctx, mock.MatchedBy(func(f model.TelemetryReadingFilter) bool {
				return f.CarID == carID && f.TimeRange == nil
			})).
			Return([]model.TelemetryReading{
				{ID: "tr-1", CarID: carID},
			}, nil)

		readings, err := svc.ListCarTelemetryHistory(ctx, validation.TelemetryReadingFilter{CarID: carID})
		assert.NoError(t, err)
		assert.Len(t, readings, 1)
		assert.Equal(t, "tr-1", readings[0].ID)
	})

	t.Run("time range is passed to repo", func(t *testing.T) {
		repo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newTestCarServiceHistory(t, nil, repo)

		from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

		repo.EXPECT().
			Find(ctx, mock.MatchedBy(func(f model.TelemetryReadingFilter) bool {
				return f.TimeRange != nil &&
					f.TimeRange.From.Equal(from) &&
					f.TimeRange.To.Equal(to)
			})).
			Return([]model.TelemetryReading{}, nil)

		_, err := svc.ListCarTelemetryHistory(ctx, validation.TelemetryReadingFilter{
			CarID: carID,
			TimeRange: &sharedvalidation.TimeRange{
				From: &from,
				To:   &to,
			},
		})
		assert.NoError(t, err)
	})

	t.Run("invalid car_id returns validation error", func(t *testing.T) {
		svc := newTestCarServiceHistory(t, nil, nil)

		_, err := svc.ListCarTelemetryHistory(ctx, validation.TelemetryReadingFilter{CarID: "not-a-uuid"})
		assert.Error(t, err)
	})

	t.Run("reversed time range returns validation error", func(t *testing.T) {
		svc := newTestCarServiceHistory(t, nil, nil)

		from := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		_, err := svc.ListCarTelemetryHistory(ctx, validation.TelemetryReadingFilter{
			CarID: carID,
			TimeRange: &sharedvalidation.TimeRange{
				From: &from,
				To:   &to,
			},
		})
		assert.Error(t, err)
	})

	t.Run("repo error is returned", func(t *testing.T) {
		repo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newTestCarServiceHistory(t, nil, repo)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrSql)

		_, err := svc.ListCarTelemetryHistory(ctx, validation.TelemetryReadingFilter{CarID: carID})
		assert.ErrorIs(t, err, model.ErrSql)
	})
}
