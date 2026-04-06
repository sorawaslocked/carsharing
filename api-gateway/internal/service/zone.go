package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type ZoneService struct {
	presenter ZonePresenter
}

func NewZoneService(presenter ZonePresenter) *ZoneService {
	return &ZoneService{presenter: presenter}
}

func (s *ZoneService) Create(ctx context.Context, data model.ZoneCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *ZoneService) Get(ctx context.Context, id string) (model.Zone, error) {
	return s.presenter.Get(ctx, id)
}

func (s *ZoneService) GetAll(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error) {
	return s.presenter.GetAll(ctx, filter)
}

func (s *ZoneService) Update(ctx context.Context, id string, data model.ZoneUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *ZoneService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}
