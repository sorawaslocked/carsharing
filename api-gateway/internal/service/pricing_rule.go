package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type PricingRuleService struct {
	presenter PricingRulePresenter
	log       *slog.Logger
}

func NewPricingRuleService(presenter PricingRulePresenter, log *slog.Logger) *PricingRuleService {
	return &PricingRuleService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.PricingRuleService"),
	}
}

func (s *PricingRuleService) Create(ctx context.Context, data model.PricingRuleCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating pricing rule", pkglog.Err(err))

		return "", err
	}

	return id, nil
}

func (s *PricingRuleService) Get(ctx context.Context, id string) (model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	rule, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting pricing rule", pkglog.Err(err))

		return model.PricingRule{}, err
	}

	return rule, nil
}

func (s *PricingRuleService) List(ctx context.Context, filter model.PricingRuleFilter) ([]model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	rules, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing pricing rules", pkglog.Err(err))

		return nil, err
	}

	return rules, nil
}

func (s *PricingRuleService) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Update(ctx, id, data); err != nil {
		log.Warn("updating pricing rule", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *PricingRuleService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Delete(ctx, id); err != nil {
		log.Warn("deleting pricing rule", pkglog.Err(err))

		return err
	}

	return nil
}
