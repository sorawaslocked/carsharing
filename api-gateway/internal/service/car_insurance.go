package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
)

type CarInsuranceService struct {
	presenter CarInsurancePresenter
}

func NewCarInsuranceService(presenter CarInsurancePresenter) *CarInsuranceService {
	return &CarInsuranceService{presenter: presenter}
}

func (s *CarInsuranceService) Create(ctx context.Context, data model.CarInsuranceCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *CarInsuranceService) Get(ctx context.Context, id string) (model.CarInsurance, error) {
	return s.presenter.Get(ctx, id)
}

func (s *CarInsuranceService) List(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	return s.presenter.List(ctx, filter)
}

func (s *CarInsuranceService) Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *CarInsuranceService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}

func (s *CarInsuranceService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	return s.presenter.GetImageUploadData(ctx)
}
