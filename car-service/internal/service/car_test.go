package service

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/service/mocks"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type noopCarCreatedNotifier struct{}

func (noopCarCreatedNotifier) OnCarCreated(model.Car) {}

func newTestCarService(t *testing.T, carRepo CarRepository, statusLogRepo CarStatusReadingRepository, eventPub EventPublisher) *CarService {
	t.Helper()
	v := validator.New()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	_ = validation.RegisterCustomValidators(v, log)
	return NewCarService(log, v, nil, carRepo, statusLogRepo, nil, nil, eventPub, noopCarCreatedNotifier{})
}

func TestCarUpdate(t *testing.T) {
	ctx := context.Background()
	carID := "c0000000-0000-4000-8000-000000000001"

	newSvc := func(t *testing.T, carRepo CarRepository) *CarService {
		t.Helper()
		return NewCarService(discardLogger(), newTestValidator(t), nil, carRepo, nil, nil, nil, nil, noopCarCreatedNotifier{})
	}

	t.Run("car updated successfully", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo)

		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)

		err := svc.Update(ctx, carID, validation.CarUpdate{})
		assert.NoError(t, err)
	})

	t.Run("car not found returns ErrCarNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo)

		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(model.ErrCarNotFound)

		err := svc.Update(ctx, carID, validation.CarUpdate{})
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})
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
	carID := "c0000000-0000-4000-8000-000000000001"
	ctx := context.Background()

	t.Run("valid transition updates car and emits log and event", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		statusLogRepo := mocks.NewMockCarStatusReadingRepository(t)
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
			Insert(ctx, mock.MatchedBy(func(e model.CarStatusReading) bool {
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

		err := svc.UpdateCarStatus(ctx, carID, validation.CarStatusUpdate{
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

		err := svc.UpdateCarStatus(ctx, carID, validation.CarStatusUpdate{
			Status: string(model.CarStatusInUse), // available → in_use not in transition map
		})
		assert.Error(t, err)
	})

	t.Run("car not found returns ErrCarNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newTestCarService(t, carRepo, nil, nil)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{}, model.ErrCarNotFound)

		err := svc.UpdateCarStatus(ctx, carID, validation.CarStatusUpdate{
			Status: string(model.CarStatusReserved),
		})
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})

	t.Run("status log insert failure is non-fatal", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		statusLogRepo := mocks.NewMockCarStatusReadingRepository(t)
		eventPub := mocks.NewMockEventPublisher(t)
		svc := newTestCarService(t, carRepo, statusLogRepo, eventPub)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusReserved}, nil)
		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)
		statusLogRepo.EXPECT().Insert(ctx, mock.Anything).Return(model.ErrSql)
		eventPub.EXPECT().
			PublishCarStatusUpdated(ctx, carID,
				string(model.CarStatusReserved),
				string(model.CarStatusAvailable),
			).
			Return(nil)

		err := svc.UpdateCarStatus(ctx, carID, validation.CarStatusUpdate{
			Status: string(model.CarStatusAvailable),
		})
		assert.NoError(t, err)
	})

	t.Run("event publish failure is non-fatal", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		statusLogRepo := mocks.NewMockCarStatusReadingRepository(t)
		eventPub := mocks.NewMockEventPublisher(t)
		svc := newTestCarService(t, carRepo, statusLogRepo, eventPub)

		carRepo.EXPECT().
			FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusReserved}, nil)
		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)
		statusLogRepo.EXPECT().Insert(ctx, mock.Anything).Return(nil)
		eventPub.EXPECT().
			PublishCarStatusUpdated(ctx, carID, mock.Anything, mock.Anything).
			Return(model.ErrSql)

		err := svc.UpdateCarStatus(ctx, carID, validation.CarStatusUpdate{
			Status: string(model.CarStatusAvailable),
		})
		assert.NoError(t, err)
	})
}

func TestCarServiceCreate(t *testing.T) {
	ctx := context.Background()
	modelID := "m0000000-0000-4000-8000-000000000001"

	newSvc := func(t *testing.T, modelRepo CarModelRepository, carRepo CarRepository, notifier CarCreatedNotifier) *CarService {
		t.Helper()
		return NewCarService(discardLogger(), newTestValidator(t), modelRepo, carRepo, nil, nil, nil, nil, notifier)
	}

	validInput := validation.CarCreate{
		ModelID:          modelID,
		VIN:              "12345678901234567",
		LicensePlate:     "ABC123",
		Color:            "red",
		YearManufactured: 2022,
		TelemetryID:      "tel-001",
	}

	t.Run("happy path returns inserted id", func(t *testing.T) {
		modelRepo := mocks.NewMockCarModelRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, modelRepo, carRepo, noopCarCreatedNotifier{})

		modelRepo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{ID: modelID}, nil)
		carRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(c model.Car) bool {
			return c.VIN == "12345678901234567" && c.Status == model.CarStatusAvailable
		})).Return("car-123", nil)

		id, err := svc.Create(ctx, validInput)
		assert.NoError(t, err)
		assert.Equal(t, "car-123", id)
	})

	t.Run("car model not found returns ErrCarModelNotFound", func(t *testing.T) {
		modelRepo := mocks.NewMockCarModelRepository(t)
		svc := newSvc(t, modelRepo, nil, noopCarCreatedNotifier{})

		modelRepo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{}, model.ErrCarModelNotFound)

		_, err := svc.Create(ctx, validInput)
		assert.ErrorIs(t, err, model.ErrCarModelNotFound)
	})

	t.Run("duplicate VIN returns ErrDuplicateVIN", func(t *testing.T) {
		modelRepo := mocks.NewMockCarModelRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, modelRepo, carRepo, noopCarCreatedNotifier{})

		modelRepo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{ID: modelID}, nil)
		carRepo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrDuplicateVIN)

		_, err := svc.Create(ctx, validInput)
		assert.ErrorIs(t, err, model.ErrDuplicateVIN)
	})

	t.Run("duplicate license plate returns ErrDuplicateLicensePlate", func(t *testing.T) {
		modelRepo := mocks.NewMockCarModelRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, modelRepo, carRepo, noopCarCreatedNotifier{})

		modelRepo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{ID: modelID}, nil)
		carRepo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrDuplicateLicensePlate)

		_, err := svc.Create(ctx, validInput)
		assert.ErrorIs(t, err, model.ErrDuplicateLicensePlate)
	})

	t.Run("validation rejects missing required fields", func(t *testing.T) {
		svc := newSvc(t, nil, nil, noopCarCreatedNotifier{})

		_, err := svc.Create(ctx, validation.CarCreate{})
		assert.Error(t, err)
	})
}

func TestCarServiceGet(t *testing.T) {
	ctx := context.Background()
	carID := "c0000000-0000-4000-8000-000000000001"

	newSvc := func(t *testing.T, carRepo CarRepository, storage ObjectStorage) *CarService {
		t.Helper()
		return NewCarService(discardLogger(), newTestValidator(t), nil, carRepo, nil, nil, storage, nil, noopCarCreatedNotifier{})
	}

	t.Run("returns car with no images", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo, nil)

		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{ID: carID, Status: model.CarStatusAvailable}, nil)

		got, err := svc.Get(ctx, carID)
		assert.NoError(t, err)
		assert.Equal(t, carID, got.ID)
	})

	t.Run("populates presigned URLs for car images", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newSvc(t, carRepo, storage)

		key := "cars/photo.jpg"
		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{
			ID:     carID,
			Images: []sharedmodel.Image{{Key: key}},
		}, nil)
		storage.EXPECT().GetPresignedURL(ctx, key).Return("https://cdn/photo.jpg", nil)

		got, err := svc.Get(ctx, carID)
		assert.NoError(t, err)
		assert.Equal(t, "https://cdn/photo.jpg", got.Images[0].URL)
	})

	t.Run("car not found returns ErrCarNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo, nil)

		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{}, model.ErrCarNotFound)

		_, err := svc.Get(ctx, carID)
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})

	t.Run("object storage error is propagated", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newSvc(t, carRepo, storage)

		key := "cars/photo.jpg"
		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{
			ID:     carID,
			Images: []sharedmodel.Image{{Key: key}},
		}, nil)
		storage.EXPECT().GetPresignedURL(ctx, key).Return("", model.ErrObjectStorage)

		_, err := svc.Get(ctx, carID)
		assert.Error(t, err)
	})
}

func TestCarServiceDelete(t *testing.T) {
	ctx := context.Background()
	carID := "c0000000-0000-4000-8000-000000000001"

	newSvc := func(t *testing.T, carRepo CarRepository) *CarService {
		t.Helper()
		return NewCarService(discardLogger(), newTestValidator(t), nil, carRepo, nil, nil, nil, nil, noopCarCreatedNotifier{})
	}

	t.Run("happy path returns nil", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo)

		carRepo.EXPECT().Delete(ctx, carID).Return(nil)

		assert.NoError(t, svc.Delete(ctx, carID))
	})

	t.Run("car not found returns ErrCarNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo)

		carRepo.EXPECT().Delete(ctx, carID).Return(model.ErrCarNotFound)

		assert.ErrorIs(t, svc.Delete(ctx, carID), model.ErrCarNotFound)
	})
}

func TestCarServiceUpdateCarTelemetry(t *testing.T) {
	ctx := context.Background()
	carID := "c0000000-0000-4000-8000-000000000001"

	newSvc := func(t *testing.T, carRepo CarRepository, telemetryRepo TelemetryReadingRepository) *CarService {
		t.Helper()
		return NewCarService(discardLogger(), newTestValidator(t), nil, carRepo, nil, telemetryRepo, nil, nil, noopCarCreatedNotifier{})
	}

	t.Run("happy path updates car telemetry", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		telemetryRepo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newSvc(t, carRepo, telemetryRepo)

		carRepo.EXPECT().Update(ctx, carID, mock.MatchedBy(func(u model.CarUpdate) bool {
			return u.MileageKM != nil && *u.MileageKM == 50_000
		})).Return(nil)
		telemetryRepo.EXPECT().Insert(ctx, mock.Anything).Return(nil)

		err := svc.UpdateCarTelemetry(ctx, carID, validation.CarTelemetryUpdate{MileageKM: 50_000})
		assert.NoError(t, err)
	})

	t.Run("car not found returns ErrCarNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo, nil)

		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(model.ErrCarNotFound)

		err := svc.UpdateCarTelemetry(ctx, carID, validation.CarTelemetryUpdate{MileageKM: 50_000})
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newSvc(t, carRepo, nil)

		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(model.ErrSql)

		err := svc.UpdateCarTelemetry(ctx, carID, validation.CarTelemetryUpdate{MileageKM: 50_000})
		assert.ErrorIs(t, err, model.ErrSql)
	})
}

func TestCarServiceGetImageUploadData(t *testing.T) {
	ctx := context.Background()

	newSvc := func(t *testing.T, storage ObjectStorage) *CarService {
		t.Helper()
		return NewCarService(discardLogger(), newTestValidator(t), nil, nil, nil, nil, storage, nil, noopCarCreatedNotifier{})
	}

	t.Run("returns upload data from object storage", func(t *testing.T) {
		storage := mocks.NewMockObjectStorage(t)
		svc := newSvc(t, storage)

		want := sharedmodel.ImageUploadData{ObjectKey: "cars/abc.jpg", PresignedPutURL: "https://upload.example.com"}
		storage.EXPECT().GetCarImageUploadData(ctx).Return(want, nil)

		got, err := svc.GetImageUploadData(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("storage error is propagated", func(t *testing.T) {
		storage := mocks.NewMockObjectStorage(t)
		svc := newSvc(t, storage)

		storage.EXPECT().GetCarImageUploadData(ctx).Return(sharedmodel.ImageUploadData{}, model.ErrObjectStorage)

		_, err := svc.GetImageUploadData(ctx)
		assert.Error(t, err)
	})
}

func TestEventHandlers(t *testing.T) {
	carID := "c0000000-0000-4000-8000-000000000001"
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
			statusLogRepo := mocks.NewMockCarStatusReadingRepository(t)
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
