package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
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

func (s *CarService) GetCarFuelHistory(ctx context.Context, carID string, filter model.CarFuelReadingFilter) ([]model.CarFuelReading, error) {
	return s.presenter.GetCarFuelHistory(ctx, carID, filter)
}

func (s *CarService) GetCarLocationHistory(ctx context.Context, carID string, filter model.CarLocationReadingFilter) ([]model.CarLocationReading, error) {
	return s.presenter.GetCarLocationHistory(ctx, carID, filter)
}

func (s *CarService) GetCarBatteryHistory(ctx context.Context, carID string, filter model.CarBatteryReadingFilter) ([]model.CarBatteryReading, error) {
	return s.presenter.GetCarBatteryHistory(ctx, carID, filter)
}

func (s *CarService) GetCarMileageHistory(ctx context.Context, carID string, filter model.CarMileageReadingFilter) ([]model.CarMileageReading, error) {
	return s.presenter.GetCarMileageHistory(ctx, carID, filter)
}

func (s *CarService) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.presenter.GetImageUploadData(ctx)
}
