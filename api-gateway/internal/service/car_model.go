package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarModelService struct {
	presenter CarModelPresenter
}

func NewCarModelService(presenter CarModelPresenter) *CarModelService {
	return &CarModelService{presenter: presenter}
}

func (s *CarModelService) Create(ctx context.Context, data model.CarModelCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *CarModelService) Get(ctx context.Context, id string) (model.CarModel, error) {
	return s.presenter.Get(ctx, id)
}

func (s *CarModelService) GetAll(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	return s.presenter.GetAll(ctx, filter)
}

func (s *CarModelService) Update(ctx context.Context, id string, data model.CarModelUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *CarModelService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}

func (s *CarModelService) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.presenter.GetImageUploadData(ctx)
}
