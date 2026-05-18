package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

type PricingRuleService struct {
	presenter PricingRulePresenter
}

func NewPricingRuleService(presenter PricingRulePresenter) *PricingRuleService {
	return &PricingRuleService{presenter: presenter}
}

func (s *PricingRuleService) Create(ctx context.Context, data model.PricingRuleCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *PricingRuleService) Get(ctx context.Context, id string) (model.PricingRule, error) {
	return s.presenter.Get(ctx, id)
}

func (s *PricingRuleService) List(ctx context.Context, filter model.PricingRuleFilter) ([]model.PricingRule, error) {
	return s.presenter.List(ctx, filter)
}

func (s *PricingRuleService) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *PricingRuleService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}
