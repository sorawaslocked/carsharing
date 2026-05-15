package service

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/service/mocks"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestCarService(t *testing.T, carRepo CarRepository, statusLogRepo CarStatusLogRepository, eventPub EventPublisher) *CarService {
	t.Helper()
	v := validator.New()
	_ = validation.RegisterCustomValidators(v)
	return NewCarService(carRepo, statusLogRepo, nil, nil, eventPub, v,
		slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestTransitionCarStatus(t *testing.T) {
	tests := []struct {
		from    model.CarStatus
		to      model.CarStatus
		allowed bool
	}{
		// from available
		{model.CarStatusAvailable, model.CarStatusReserved, true},
		{model.CarStatusAvailable, model.CarStatusMaintenance, true},
		{model.CarStatusAvailable, model.CarStatusOutOfService, true},
		{model.CarStatusAvailable, model.CarStatusInUse, false},
		// from reserved
		{model.CarStatusReserved, model.CarStatusInUse, true},
		{model.CarStatusReserved, model.CarStatusAvailable, true},
		{model.CarStatusReserved, model.CarStatusMaintenance, true},
		{model.CarStatusReserved, model.CarStatusOutOfService, true},
		// from in_use
		{model.CarStatusInUse, model.CarStatusAvailable, true},
		{model.CarStatusInUse, model.CarStatusMaintenance, true},
		{model.CarStatusInUse, model.CarStatusOutOfService, true},
		{model.CarStatusInUse, model.CarStatusReserved, false},
		// from maintenance
		{model.CarStatusMaintenance, model.CarStatusAvailable, true},
		{model.CarStatusMaintenance, model.CarStatusOutOfService, true},
		{model.CarStatusMaintenance, model.CarStatusReserved, false},
		{model.CarStatusMaintenance, model.CarStatusInUse, false},
		// from out_of_service (terminal)
		{model.CarStatusOutOfService, model.CarStatusAvailable, false},
		{model.CarStatusOutOfService, model.CarStatusReserved, false},
		{model.CarStatusOutOfService, model.CarStatusInUse, false},
		{model.CarStatusOutOfService, model.CarStatusMaintenance, false},
	}

	for _, tt := range tests {
		err := transitionCarStatus(tt.from, tt.to)
		if tt.allowed {
			assert.NoError(t, err, "%s → %s should be allowed", tt.from, tt.to)
		} else {
			assert.Error(t, err, "%s → %s should be rejected", tt.from, tt.to)
		}
	}
}

func TestUpdateCarStatus(t *testing.T) {
	carID := "car-123"
	ctx := context.Background()

	t.Run("valid transition updates car and emits log and event", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		statusLogRepo := mocks.NewMockCarStatusLogRepository(t)
		eventPub := mocks.NewMockEventPublisher(t)
		svc := newTestCarService(t, carRepo, statusLogRepo, eventPub)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusAvailable}, nil)
		carRepo.EXPECT().
			Update(ctx, carID, mock.MatchedBy(func(u model.CarUpdate) bool {
				return u.Status != nil && *u.Status == model.CarStatusReserved
			})).
			Return(nil)
		statusLogRepo.EXPECT().
			Insert(ctx, mock.MatchedBy(func(e model.CarStatusLogEntry) bool {
				return e.CarID == carID &&
					e.FromStatus == model.CarStatusAvailable &&
					e.ToStatus == model.CarStatusReserved
			})).
			Return(nil)
		eventPub.EXPECT().
			PublishCarStatusUpdated(ctx, carID,
				string(model.CarStatusAvailable),
				string(model.CarStatusReserved),
			).
			Return(nil)

		err := svc.UpdateCarStatus(ctx, carID, model.CarStatusUpdateInput{
			Status: string(model.CarStatusReserved),
		})
		assert.NoError(t, err)
	})

	t.Run("invalid transition returns error without touching the repo", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newTestCarService(t, carRepo, nil, nil)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusAvailable}, nil)

		err := svc.UpdateCarStatus(ctx, carID, model.CarStatusUpdateInput{
			Status: string(model.CarStatusInUse), // available → in_use not in transition map
		})
		assert.Error(t, err)
	})

	t.Run("car not found returns ErrNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newTestCarService(t, carRepo, nil, nil)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{}, model.ErrNotFound)

		err := svc.UpdateCarStatus(ctx, carID, model.CarStatusUpdateInput{
			Status: string(model.CarStatusReserved),
		})
		assert.ErrorIs(t, err, model.ErrNotFound)
	})

	t.Run("status log insert failure is non-fatal", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		statusLogRepo := mocks.NewMockCarStatusLogRepository(t)
		eventPub := mocks.NewMockEventPublisher(t)
		svc := newTestCarService(t, carRepo, statusLogRepo, eventPub)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusReserved}, nil)
		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)
		statusLogRepo.EXPECT().Insert(ctx, mock.Anything).Return(model.ErrInternalServerError)
		eventPub.EXPECT().
			PublishCarStatusUpdated(ctx, carID,
				string(model.CarStatusReserved),
				string(model.CarStatusAvailable),
			).
			Return(nil)

		err := svc.UpdateCarStatus(ctx, carID, model.CarStatusUpdateInput{
			Status: string(model.CarStatusAvailable),
		})
		assert.NoError(t, err)
	})

	t.Run("event publish failure is non-fatal", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		statusLogRepo := mocks.NewMockCarStatusLogRepository(t)
		eventPub := mocks.NewMockEventPublisher(t)
		svc := newTestCarService(t, carRepo, statusLogRepo, eventPub)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusReserved}, nil)
		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)
		statusLogRepo.EXPECT().Insert(ctx, mock.Anything).Return(nil)
		eventPub.EXPECT().
			PublishCarStatusUpdated(ctx, carID, mock.Anything, mock.Anything).
			Return(model.ErrInternalServerError)

		err := svc.UpdateCarStatus(ctx, carID, model.CarStatusUpdateInput{
			Status: string(model.CarStatusAvailable),
		})
		assert.NoError(t, err)
	})
}

func TestEventHandlers(t *testing.T) {
	carID := "car-123"
	ctx := context.Background()

	tests := []struct {
		name       string
		fromStatus model.CarStatus
		toStatus   model.CarStatus
		trigger    func(svc *CarService) error
	}{
		{
			name:       "OnBookingCreated reserves available car",
			fromStatus: model.CarStatusAvailable,
			toStatus:   model.CarStatusReserved,
			trigger: func(svc *CarService) error {
				return svc.OnBookingCreated(ctx, model.BookingCreatedEvent{
					BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
		{
			name:       "OnBookingCancelled releases reserved car",
			fromStatus: model.CarStatusReserved,
			toStatus:   model.CarStatusAvailable,
			trigger: func(svc *CarService) error {
				return svc.OnBookingCancelled(ctx, model.BookingCancelledEvent{
					BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
		{
			name:       "OnBookingExpired releases reserved car",
			fromStatus: model.CarStatusReserved,
			toStatus:   model.CarStatusAvailable,
			trigger: func(svc *CarService) error {
				return svc.OnBookingExpired(ctx, model.BookingExpiredEvent{
					BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
		{
			name:       "OnBookingCompleted releases in-use car",
			fromStatus: model.CarStatusInUse,
			toStatus:   model.CarStatusAvailable,
			trigger: func(svc *CarService) error {
				return svc.OnBookingCompleted(ctx, model.BookingCompletedEvent{
					BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
		{
			name:       "OnTripStarted sets reserved car to in_use",
			fromStatus: model.CarStatusReserved,
			toStatus:   model.CarStatusInUse,
			trigger: func(svc *CarService) error {
				return svc.OnTripStarted(ctx, model.TripStartedEvent{
					TripID: "trip-1", BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
		{
			name:       "OnTripEnded releases in-use car",
			fromStatus: model.CarStatusInUse,
			toStatus:   model.CarStatusAvailable,
			trigger: func(svc *CarService) error {
				return svc.OnTripEnded(ctx, model.TripEndedEvent{
					TripID: "trip-1", BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
		{
			name:       "OnTripCancelled releases in-use car",
			fromStatus: model.CarStatusInUse,
			toStatus:   model.CarStatusAvailable,
			trigger: func(svc *CarService) error {
				return svc.OnTripCancelled(ctx, model.TripCancelledEvent{
					TripID: "trip-1", BookingID: "booking-1", CarID: carID, UserID: "user-1",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carRepo := mocks.NewMockCarRepository(t)
			statusLogRepo := mocks.NewMockCarStatusLogRepository(t)
			eventPub := mocks.NewMockEventPublisher(t)
			svc := newTestCarService(t, carRepo, statusLogRepo, eventPub)

			carRepo.EXPECT().
				FindByID(ctx, carID).
				Return(model.Car{ID: carID, Status: tt.fromStatus}, nil)
			carRepo.EXPECT().
				Update(ctx, carID, mock.MatchedBy(func(u model.CarUpdate) bool {
					return u.Status != nil && *u.Status == tt.toStatus
				})).
				Return(nil)
			statusLogRepo.EXPECT().Insert(ctx, mock.Anything).Return(nil)
			eventPub.EXPECT().
				PublishCarStatusUpdated(ctx, carID,
					string(tt.fromStatus),
					string(tt.toStatus),
				).
				Return(nil)

			assert.NoError(t, tt.trigger(svc))
		})
	}
}
