package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
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

func (s *CarService) GetAll(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	return s.presenter.GetAll(ctx, filter)
}

func (s *CarService) Update(ctx context.Context, id string, data model.CarUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *CarService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}

func (s *CarService) GetCarStatusLog(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error) {
	return s.presenter.GetCarStatusLog(ctx, filter)
}

func (s *CarService) GetCarFuelHistory(ctx context.Context, filter model.CarFuelReadingFilter) ([]model.CarFuelReading, error) {
	return s.presenter.GetCarFuelHistory(ctx, filter)
}

func (s *CarService) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.presenter.GetImageUploadData(ctx)
}
