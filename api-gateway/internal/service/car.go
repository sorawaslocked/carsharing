package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type CarService struct {
	presenter CarPresenter
	log       *slog.Logger
}

func NewCarService(presenter CarPresenter, log *slog.Logger) *CarService {
	return &CarService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.CarService"),
	}
}

func (s *CarService) Create(ctx context.Context, data model.CarCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))
	log.Debug("creating car")

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating car", pkglog.Err(err))

		return "", err
	}

	log.Debug("car created", slog.String("id", id))

	return id, nil
}

func (s *CarService) Get(ctx context.Context, id string) (model.Car, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("getting car")

	car, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting car", pkglog.Err(err))

		return model.Car{}, err
	}

	return car, nil
}

func (s *CarService) List(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))
	log.Debug("listing cars")

	cars, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing cars", pkglog.Err(err))

		return nil, err
	}

	return cars, nil
}

func (s *CarService) Update(ctx context.Context, id string, data model.CarUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))
	log.Debug("updating car")

	if err := s.presenter.Update(ctx, id, data); err != nil {
		log.Warn("updating car", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))
	log.Debug("deleting car")

	if err := s.presenter.Delete(ctx, id); err != nil {
		log.Warn("deleting car", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarService) UpdateTelemetry(ctx context.Context, carID string, data model.CarTelemetryUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateTelemetry"), utils.MetadataFromCtx(ctx))
	log.Debug("updating car telemetry")

	if err := s.presenter.UpdateTelemetry(ctx, carID, data); err != nil {
		log.Warn("updating car telemetry", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarService) UpdateStatus(ctx context.Context, carID string, data model.CarStatusUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateStatus"), utils.MetadataFromCtx(ctx))
	log.Debug("updating car status")

	if err := s.presenter.UpdateStatus(ctx, carID, data); err != nil {
		log.Warn("updating car status", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarService) GetCarStatusHistory(ctx context.Context, carID string, filter model.CarStatusReadingFilter) ([]model.CarStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetCarStatusHistory"), utils.MetadataFromCtx(ctx))
	log.Debug("getting car status history")

	history, err := s.presenter.GetCarStatusHistory(ctx, carID, filter)
	if err != nil {
		log.Warn("getting car status history", pkglog.Err(err))

		return nil, err
	}

	return history, nil
}

func (s *CarService) GetCarTelemetryHistory(ctx context.Context, carID string, filter model.CarTelemetryReadingFilter) ([]model.CarTelemetryReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetCarTelemetryHistory"), utils.MetadataFromCtx(ctx))
	log.Debug("getting car telemetry history")

	history, err := s.presenter.GetCarTelemetryHistory(ctx, carID, filter)
	if err != nil {
		log.Warn("getting car telemetry history", pkglog.Err(err))

		return nil, err
	}

	return history, nil
}

func (s *CarService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))
	log.Debug("getting car image upload data")

	data, err := s.presenter.GetImageUploadData(ctx)
	if err != nil {
		log.Warn("getting car image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}
