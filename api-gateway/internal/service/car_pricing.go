package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarPricingRuleService struct {
	presenter CarPricingRulePresenter
}

func NewCarPricingRuleService(presenter CarPricingRulePresenter) *CarPricingRuleService {
	return &CarPricingRuleService{presenter: presenter}
}

func (s *CarPricingRuleService) Create(ctx context.Context, data model.CarPricingRuleCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *CarPricingRuleService) Get(ctx context.Context, id string) (model.CarPricingRule, error) {
	return s.presenter.Get(ctx, id)
}

func (s *CarPricingRuleService) List(ctx context.Context, filter model.CarPricingRuleFilter) ([]model.CarPricingRule, error) {
	return s.presenter.List(ctx, filter)
}

func (s *CarPricingRuleService) Update(ctx context.Context, id string, data model.CarPricingRuleUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *CarPricingRuleService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}
