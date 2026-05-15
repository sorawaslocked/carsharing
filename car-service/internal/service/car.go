package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"

	"github.com/go-playground/validator/v10"
)

type CarService struct {
	carRepo            CarRepository
	statusLogRepo      CarStatusLogRepository
	telematicsRepo     TelematicsRepository
	objectStorage      ObjectStorage
	eventPublisher     EventPublisher
	carCreatedNotifier CarCreatedNotifier

	validate *validator.Validate
	log      *slog.Logger
}

func NewCarService(
	carRepo CarRepository,
	statusLogRepo CarStatusLogRepository,
	telematicsRepo TelematicsRepository,
	objectStorage ObjectStorage,
	eventPublisher EventPublisher,
	validate *validator.Validate,
	log *slog.Logger,
) *CarService {
	s := &CarService{
		carRepo:        carRepo,
		statusLogRepo:  statusLogRepo,
		telematicsRepo: telematicsRepo,
		objectStorage:  objectStorage,
		eventPublisher: eventPublisher,
		validate:       validate,
	}

	s.log = pkglog.WithComponent(log, "service.CarService")

	return s
}

func (s *CarService) SetCarCreatedNotifier(n CarCreatedNotifier) {
	s.carCreatedNotifier = n
}

func (s *CarService) Create(ctx context.Context, createInput model.CarCreateInput) (string, error) {
	const method = "Create"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	err := validation.ValidateInput(s.validate, createInput)
	if err != nil {
		return "", handleError(logger, err)
	}

	now := time.Now()

	car := model.Car{
		ModelID:          createInput.ModelID,
		VIN:              createInput.VIN,
		LicensePlate:     createInput.LicensePlate,
		Color:            createInput.Color,
		YearManufactured: createInput.YearManufactured,
		MileageKM:        createInput.MileageKM,
		FuelLevel:        createInput.FuelLevel,
		BatteryLevel:     createInput.BatteryLevel,
		Notes:            createInput.Notes,
		Status:           model.CarStatusAvailable,
		LastSeenAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	id, err := s.carRepo.Insert(ctx, car)
	if err != nil {
		return "", handleError(logger, err)
	}

	car.ID = id
	if s.carCreatedNotifier != nil {
		s.carCreatedNotifier.OnCarCreated(car)
	}

	return id, nil
}

func (s *CarService) Get(ctx context.Context, id string) (model.Car, error) {
	const method = "Get"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	car, err := s.carRepo.FindByID(ctx, id)
	if err != nil {
		return model.Car{}, handleError(logger, err)
	}

	for i := range car.Images {
		url, err := s.objectStorage.GetPresignedURL(ctx, *car.Images[i].Key)
		if err != nil {
			return model.Car{}, handleError(logger, err)
		}
		car.Images[i].URL = &url
	}

	return car, nil
}

func (s *CarService) GetAll(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error) {
	const method = "GetAll"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	err := validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return nil, handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput)

	cars, err := s.carRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	for i := range cars {
		for j := range cars[i].Images {
			url, err := s.objectStorage.GetPresignedURL(ctx, *cars[i].Images[j].Key)
			if err != nil {
				return nil, handleError(logger, err)
			}
			cars[i].Images[j].URL = &url
		}
	}

	return cars, nil
}

func (s *CarService) Update(ctx context.Context, id string, updateInput model.CarUpdateInput) error {
	const method = "Update"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	err := validation.ValidateInput(s.validate, updateInput)
	if err != nil {
		return handleError(logger, err)
	}

	update := model.CarUpdate{
		ModelID:      updateInput.ModelID,
		LicensePlate: updateInput.LicensePlate,
		Color:        updateInput.Color,
		Notes:        updateInput.Notes,
		ImageKeys:    updateInput.ImageKeys,
		UpdatedAt:    time.Now(),
	}

	err = s.carRepo.Update(ctx, id, update)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarService) UpdateCarStatus(ctx context.Context, id string, statusInput model.CarStatusUpdateInput) error {
	const method = "UpdateCarStatus"
	md := utils.MetadataFromCtx(ctx)
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, statusInput); err != nil {
		return handleError(logger, err)
	}
	toStatus, _ := model.ParseCarStatus(statusInput.Status)

	current, err := s.carRepo.FindByID(ctx, id)
	if err != nil {
		return handleError(logger, err)
	}

	if err = transitionCarStatus(current.Status, toStatus); err != nil {
		return handleError(logger, err)
	}

	now := time.Now()

	if err = s.carRepo.Update(ctx, id, model.CarUpdate{
		Status:    &toStatus,
		UpdatedAt: now,
	}); err != nil {
		return handleError(logger, err)
	}

	actorType := model.CarStatusActorSystem
	var actorID *string
	if md.UserID != nil {
		actorType = model.CarStatusActorUser
		actorID = md.UserID
	}

	if s.statusLogRepo != nil {
		if err = s.statusLogRepo.Insert(ctx, model.CarStatusLogEntry{
			CarID:      current.ID,
			FromStatus: current.Status,
			ToStatus:   toStatus,
			ActorType:  actorType,
			ActorID:    actorID,
			ChangedAt:  now,
		}); err != nil {
			logger.Error("failed to write status log entry",
				slog.String("carID", current.ID),
				slog.String("from", string(current.Status)),
				slog.String("to", string(toStatus)),
			)
		}
	}

	if s.eventPublisher != nil {
		if pubErr := s.eventPublisher.PublishCarStatusUpdated(ctx, current.ID, string(current.Status), string(toStatus)); pubErr != nil {
			logger.Error("failed to publish car status updated event",
				slog.String("carID", current.ID),
				pkglog.Err(pubErr),
			)
		}
	}

	logger.Info("car status updated",
		slog.String("carID", current.ID),
		slog.String("from", string(current.Status)),
		slog.String("to", string(toStatus)),
		slog.String("actorType", string(actorType)),
	)

	return nil
}

func (s *CarService) UpdateCarTelemetry(ctx context.Context, id string, input model.CarTelematicsUpdateInput) error {
	const method = "UpdateCarTelemetry"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	now := time.Now()

	update := model.CarUpdate{
		MileageKM:    &input.MileageKM,
		FuelLevel:    input.FuelLevel,
		BatteryLevel: input.BatteryLevel,
		Location:     input.Location,
		LastSeenAt:   &now,
		UpdatedAt:    now,
	}

	if err := s.carRepo.Update(ctx, id, update); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarService) GetCarStatusHistory(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error) {
	const method = "GetCarStatusHistory"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if s.statusLogRepo == nil {
		return nil, nil
	}

	entries, err := s.statusLogRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return entries, nil
}

func (s *CarService) GetCarFuelHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error) {
	const method = "GetCarFuelHistory"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	return s.findTelematicsEvents(ctx, logger, filter)
}

func (s *CarService) GetCarLocationHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error) {
	const method = "GetCarLocationHistory"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	return s.findTelematicsEvents(ctx, logger, filter)
}

func (s *CarService) GetCarBatteryHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error) {
	const method = "GetCarBatteryHistory"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	return s.findTelematicsEvents(ctx, logger, filter)
}

func (s *CarService) GetCarMileageHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error) {
	const method = "GetCarMileageHistory"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	return s.findTelematicsEvents(ctx, logger, filter)
}

func (s *CarService) findTelematicsEvents(ctx context.Context, logger *slog.Logger, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error) {
	if s.telematicsRepo == nil {
		return nil, nil
	}

	events, err := s.telematicsRepo.FindEvents(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return events, nil
}

func (s *CarService) Delete(ctx context.Context, id string) error {
	const method = "Delete"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.carRepo.Delete(ctx, id); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarService) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	const method = "GetImageUploadData"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetCarImageUploadData(ctx)
	if err != nil {
		return model.ImageUploadData{}, handleError(logger, err)
	}

	return data, nil
}

func (s *CarService) OnBookingCreated(ctx context.Context, event model.BookingCreatedEvent) error {
	const method = "OnBookingCreated"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusReserved),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car reserved on booking created",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnTripStarted(ctx context.Context, event model.TripStartedEvent) error {
	const method = "OnTripStarted"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusInUse),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car set to in_use on trip started",
		slog.String("carID", event.CarID),
		slog.String("tripID", event.TripID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnBookingCancelled(ctx context.Context, event model.BookingCancelledEvent) error {
	const method = "OnBookingCancelled"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car released on booking cancelled",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnTripEnded(ctx context.Context, event model.TripEndedEvent) error {
	const method = "OnTripEnded"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car released on trip ended",
		slog.String("carID", event.CarID),
		slog.String("tripID", event.TripID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnBookingExpired(ctx context.Context, event model.BookingExpiredEvent) error {
	const method = "OnBookingExpired"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car released on booking expired",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnBookingCompleted(ctx context.Context, event model.BookingCompletedEvent) error {
	const method = "OnBookingCompleted"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car released on booking completed",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnTripCancelled(ctx context.Context, event model.TripCancelledEvent) error {
	const method = "OnTripCancelled"
	logger := pkglog.WithMethod(s.log, method)
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, model.CarStatusUpdateInput{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("car released on trip cancelled",
		slog.String("carID", event.CarID),
		slog.String("tripID", event.TripID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}
