package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	sharedvalidation "carsharing/shared/validation"

	"github.com/go-playground/validator/v10"
)

type CarService struct {
	log      *slog.Logger
	validate *validator.Validate

	carModelRepo         CarModelRepository
	carRepo              CarRepository
	zoneRepo             ZoneRepository
	statusReadingRepo    CarStatusReadingRepository
	telemetryReadingRepo TelemetryReadingRepository
	objectStorage        ObjectStorage
	eventPublisher       EventPublisher
	carCreatedNotifier   CarCreatedNotifier
}

func NewCarService(
	log *slog.Logger,
	validate *validator.Validate,
	carModelRepo CarModelRepository,
	carRepo CarRepository,
	zoneRepo ZoneRepository,
	statusReadingRepo CarStatusReadingRepository,
	telemetryReadingRepo TelemetryReadingRepository,
	objectStorage ObjectStorage,
	eventPublisher EventPublisher,
) *CarService {
	return &CarService{
		log:                  pkglog.WithComponent(log, "service.CarService"),
		validate:             validate,
		carModelRepo:         carModelRepo,
		carRepo:              carRepo,
		zoneRepo:             zoneRepo,
		statusReadingRepo:    statusReadingRepo,
		telemetryReadingRepo: telemetryReadingRepo,
		objectStorage:        objectStorage,
		eventPublisher:       eventPublisher,
	}
}

func (s *CarService) SetCarCreatedNotifier(n CarCreatedNotifier) {
	s.carCreatedNotifier = n
}

func (s *CarService) Create(ctx context.Context, data validation.CarCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	if _, err := s.carModelRepo.FindByID(ctx, data.ModelID); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: finding car model by id", pkglog.Err(err))
		}
		return "", err
	}

	if data.ZoneID != nil {
		if _, err := s.zoneRepo.FindByID(ctx, *data.ZoneID); err != nil {
			if !errors.Is(err, model.ErrNotFound) {
				log.Error("repo: finding zone by id", pkglog.Err(err))
			}
			return "", err
		}
	}

	now := time.Now()

	car := model.Car{
		ModelID:          data.ModelID,
		VIN:              data.VIN,
		LicensePlate:     data.LicensePlate,
		Color:            data.Color,
		YearManufactured: data.YearManufactured,
		TelemetryID:      data.TelemetryID,
		ZoneID:           data.ZoneID,
		FuelLevel:        data.FuelLevel,
		BatteryLevel:     data.BatteryLevel,
		Notes:            data.Notes,
		Status:           model.CarStatusAvailable,
		LastSeenAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if data.MileageKM != nil {
		car.MileageKM = *data.MileageKM
	}
	if data.Location != nil {
		car.Location = sharedmodel.Location{
			Latitude:  data.Location.Latitude,
			Longitude: data.Location.Longitude,
		}
	}

	id, err := s.carRepo.Insert(ctx, car)
	if err != nil {
		log.Error("repo: inserting car", pkglog.Err(err))
		return "", err
	}

	car.ID = id
	if s.carCreatedNotifier != nil {
		s.carCreatedNotifier.OnCarCreated(car)
	}

	return id, nil
}

func (s *CarService) Get(ctx context.Context, id string) (model.Car, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.Car{}, err
	}

	car, err := s.carRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: finding car by id", pkglog.Err(err))
		}
		return model.Car{}, err
	}

	for i := range car.Images {
		url, err := s.objectStorage.GetPresignedURL(ctx, car.Images[i].Key)
		if err != nil {
			log.Error("object storage: getting presigned url", pkglog.Err(err))
			return model.Car{}, err
		}
		car.Images[i].URL = url
	}

	return car, nil
}

func (s *CarService) List(ctx context.Context, filter validation.CarFilter) ([]model.Car, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	cars, err := s.carRepo.Find(ctx, carFilter(filter))
	if err != nil {
		log.Error("repo: listing cars", pkglog.Err(err))
		return nil, err
	}

	for i := range cars {
		for j := range cars[i].Images {
			url, err := s.objectStorage.GetPresignedURL(ctx, cars[i].Images[j].Key)
			if err != nil {
				log.Error("object storage: getting presigned url", pkglog.Err(err))
				return nil, err
			}
			cars[i].Images[j].URL = url
		}
	}

	return cars, nil
}

func (s *CarService) Update(ctx context.Context, id string, data validation.CarUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	if data.ZoneID != nil {
		if _, err := s.zoneRepo.FindByID(ctx, *data.ZoneID); err != nil {
			if !errors.Is(err, model.ErrNotFound) {
				log.Error("repo: finding zone by id", pkglog.Err(err))
			}
			return err
		}
	}

	if err := s.carRepo.Update(ctx, id, model.CarUpdate{
		ModelID:      data.ModelID,
		LicensePlate: data.LicensePlate,
		Color:        data.Color,
		TelemetryID:  data.TelemetryID,
		ZoneID:       data.ZoneID,
		IsRetired:    data.IsRetired,
		Notes:        data.Notes,
		ImageKeys:    data.ImageKeys,
		UpdatedAt:    time.Now(),
	}); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: updating car", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarService) UpdateCarStatus(ctx context.Context, id string, data validation.CarStatusUpdate) error {
	md := utils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateCarStatus"), md)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	toStatus, _ := model.CarStatusFromString(data.Status)

	current, err := s.carRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: finding car by id", pkglog.Err(err))
		}
		return err
	}

	if err = transitionCarStatus(current.Status, toStatus); err != nil {
		return validation.Errors{"status": err}
	}

	now := time.Now()

	if err = s.carRepo.Update(ctx, id, model.CarUpdate{
		Status:    &toStatus,
		UpdatedAt: now,
	}); err != nil {
		log.Error("repo: updating car status", pkglog.Err(err))
		return err
	}

	actorType := sharedmodel.ActorTypeSystem
	var actorID *string
	if md.UserID != nil {
		actorType = sharedmodel.ActorTypeUser
		actorID = md.UserID
	}

	if s.statusReadingRepo != nil {
		if err = s.statusReadingRepo.Insert(ctx, model.CarStatusReading{
			CarID:      current.ID,
			FromStatus: current.Status,
			ToStatus:   toStatus,
			ActorType:  actorType,
			ActorID:    actorID,
			RecordedAt: now,
		}); err != nil {
			log.Error("repo: inserting status log entry",
				slog.String("carID", current.ID),
				slog.String("from", current.Status.String()),
				slog.String("to", toStatus.String()),
				pkglog.Err(err),
			)
		}
	}

	if s.eventPublisher != nil {
		if err := s.eventPublisher.PublishCarStatusUpdated(ctx, current.ID, current.Status.String(), toStatus.String()); err != nil {
			log.Error("event: publishing car status updated",
				slog.String("carID", current.ID),
				pkglog.Err(err),
			)
		}
	}

	log.Info("car status updated",
		slog.String("carID", current.ID),
		slog.String("from", current.Status.String()),
		slog.String("to", toStatus.String()),
		slog.String("actorType", string(actorType)),
	)

	return nil
}

func (s *CarService) UpdateCarTelemetry(ctx context.Context, id string, update validation.CarTelemetryUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateCarTelemetry"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, update); err != nil {
		return err
	}

	now := time.Now()

	carUpdate := model.CarUpdate{
		MileageKM:    &update.MileageKM,
		FuelLevel:    update.FuelLevel,
		BatteryLevel: update.BatteryLevel,
		LastSeenAt:   &now,
		UpdatedAt:    now,
	}
	if update.Location != nil {
		loc := sharedmodel.Location{Latitude: update.Location.Latitude, Longitude: update.Location.Longitude}
		carUpdate.Location = &loc
	}
	if err := s.carRepo.Update(ctx, id, carUpdate); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: updating car telemetry", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarService) ListCarStatusHistory(ctx context.Context, filter validation.CarStatusReadingFilter) ([]model.CarStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListCarStatusHistory"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	domainFilter := model.CarStatusReadingFilter{
		CarID:      filter.CarID,
		Pagination: &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset},
	}
	if filter.FromStatus != nil {
		s := model.CarStatus(*filter.FromStatus)
		domainFilter.FromStatus = &s
	}
	if filter.ToStatus != nil {
		s := model.CarStatus(*filter.ToStatus)
		domainFilter.ToStatus = &s
	}
	if filter.TimeRange != nil {
		tr := sharedmodel.TimeRange{}
		if filter.TimeRange.From != nil {
			tr.From = *filter.TimeRange.From
		}
		if filter.TimeRange.To != nil {
			tr.To = *filter.TimeRange.To
		}
		domainFilter.TimeRange = &tr
	}

	entries, err := s.statusReadingRepo.Find(ctx, domainFilter)
	if err != nil {
		log.Error("repo: listing car status history", pkglog.Err(err))
		return nil, err
	}

	return entries, nil
}

func (s *CarService) ListCarTelemetryHistory(ctx context.Context, filter validation.TelemetryReadingFilter) ([]model.TelemetryReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListCarTelemetryHistory"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	if s.telemetryReadingRepo == nil {
		return nil, nil
	}

	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	domainTelemetryFilter := model.TelemetryReadingFilter{
		CarID:      filter.CarID,
		Pagination: &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset},
	}
	if filter.TimeRange != nil {
		tr := sharedmodel.TimeRange{}
		if filter.TimeRange.From != nil {
			tr.From = *filter.TimeRange.From
		}
		if filter.TimeRange.To != nil {
			tr.To = *filter.TimeRange.To
		}
		domainTelemetryFilter.TimeRange = &tr
	}
	events, err := s.telemetryReadingRepo.Find(ctx, domainTelemetryFilter)
	if err != nil {
		log.Error("repo: listing car telemetry history", pkglog.Err(err))
		return nil, err
	}

	return events, nil
}

func (s *CarService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.carRepo.Delete(ctx, id); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: deleting car", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetCarImageUploadData(ctx)
	if err != nil {
		log.Error("object storage: getting upload data", pkglog.Err(err))
		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func (s *CarService) OnBookingCreated(ctx context.Context, event model.BookingCreatedEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnBookingCreated"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusReserved),
	}); err != nil {
		return err
	}

	log.Info("car reserved on booking created",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnTripStarted(ctx context.Context, event model.TripStartedEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnTripStarted"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusInUse),
	}); err != nil {
		return err
	}

	log.Info("car set to in_use on trip started",
		slog.String("carID", event.CarID),
		slog.String("tripID", event.TripID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnBookingCancelled(ctx context.Context, event model.BookingCancelledEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnBookingCancelled"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return err
	}

	log.Info("car released on booking cancelled",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnTripEnded(ctx context.Context, event model.TripEndedEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnTripEnded"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return err
	}

	log.Info("car released on trip ended",
		slog.String("carID", event.CarID),
		slog.String("tripID", event.TripID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnBookingExpired(ctx context.Context, event model.BookingExpiredEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnBookingExpired"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return err
	}

	log.Info("car released on booking expired",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnBookingCompleted(ctx context.Context, event model.BookingCompletedEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnBookingCompleted"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return err
	}

	log.Info("car released on booking completed",
		slog.String("carID", event.CarID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func (s *CarService) OnTripCancelled(ctx context.Context, event model.TripCancelledEvent) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "OnTripCancelled"), utils.MetadataFromCtx(ctx))

	if err := s.UpdateCarStatus(ctx, event.CarID, validation.CarStatusUpdate{
		Status: string(model.CarStatusAvailable),
	}); err != nil {
		return err
	}

	log.Info("car released on trip cancelled",
		slog.String("carID", event.CarID),
		slog.String("tripID", event.TripID),
		slog.String("bookingID", event.BookingID),
	)

	return nil
}

func carFilter(filter validation.CarFilter) model.CarFilter {
	repoFilter := model.CarFilter{}
	if filter.Status != nil {
		status, _ := model.CarStatusFromString(*filter.Status)
		repoFilter.Status = &status
	}
	if filter.ModelFilter != nil {
		mf := model.CarModelFilter{
			Brand:    filter.ModelFilter.Brand,
			Model:    filter.ModelFilter.Model,
			MinSeats: filter.ModelFilter.MinSeats,
		}
		if filter.ModelFilter.FuelType != nil {
			ft, _ := model.CarFuelTypeFromString(*filter.ModelFilter.FuelType)
			mf.FuelType = &ft
		}
		if filter.ModelFilter.Transmission != nil {
			tr, _ := model.CarTransmissionFromString(*filter.ModelFilter.Transmission)
			mf.Transmission = &tr
		}
		if filter.ModelFilter.BodyType != nil {
			bt, _ := model.CarBodyTypeFromString(*filter.ModelFilter.BodyType)
			mf.BodyType = &bt
		}
		if filter.ModelFilter.Class != nil {
			cl, _ := model.CarClassFromString(*filter.ModelFilter.Class)
			mf.Class = &cl
		}
		if filter.ModelFilter.Pagination == nil {
			filter.ModelFilter.Pagination = sharedvalidation.DefaultPagination()
		}
		mf.Pagination = &sharedmodel.Pagination{Limit: filter.ModelFilter.Pagination.Limit, Offset: filter.ModelFilter.Pagination.Offset}
		repoFilter.ModelFilter = &mf
	}
	if filter.LocationFilter != nil {
		repoFilter.LocationFilter = &model.LocationFilter{
			Location: sharedmodel.Location{
				Latitude:  filter.LocationFilter.Location.Latitude,
				Longitude: filter.LocationFilter.Location.Longitude,
			},
			RadiusKM: filter.LocationFilter.RadiusKM,
		}
	}
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	repoFilter.Pagination = &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset}
	return repoFilter
}
