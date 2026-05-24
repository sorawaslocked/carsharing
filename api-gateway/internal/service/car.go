package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
)

type CarService struct {
	presenter CarPresenter
}

func NewCarService(presenter CarPresenter) *CarService {
	return &CarService{presenter: presenter}
}

func (s *CarService) Create(ctx context.Context, data model.CarCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *CarService) Get(ctx context.Context, id string) (model.Car, error) {
	return s.presenter.Get(ctx, id)
}

func (s *CarService) List(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	return s.presenter.List(ctx, filter)
}

func (s *CarService) Update(ctx context.Context, id string, data model.CarUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *CarService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}

func (s *CarService) UpdateTelemetry(ctx context.Context, carID string, data model.CarTelemetryUpdate) error {
	return s.presenter.UpdateTelemetry(ctx, carID, data)
}

func (s *CarService) UpdateStatus(ctx context.Context, carID string, data model.CarStatusUpdate) error {
	return s.presenter.UpdateStatus(ctx, carID, data)
}

func (s *CarService) GetCarStatusHistory(ctx context.Context, carID string, filter model.CarStatusReadingFilter) ([]model.CarStatusReading, error) {
	return s.presenter.GetCarStatusHistory(ctx, carID, filter)
}

func (s *CarService) GetCarTelemetryHistory(ctx context.Context, carID string, filter model.CarTelemetryReadingFilter) ([]model.CarTelemetryReading, error) {
	return s.presenter.GetCarTelemetryHistory(ctx, carID, filter)
}

func (s *CarService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	return s.presenter.GetImageUploadData(ctx)
}
