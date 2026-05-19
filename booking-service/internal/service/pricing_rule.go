package service

import (
	"context"
	"log/slog"

	"carsharing/booking-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type PricingRuleService struct {
	log      *slog.Logger
	ruleRepo PricingRuleRepository
}

func NewPricingRuleService(log *slog.Logger, ruleRepo PricingRuleRepository) *PricingRuleService {
	return &PricingRuleService{
		log:      pkglog.WithComponent(log, "service.PricingRuleService"),
		ruleRepo: ruleRepo,
	}
}

func (s *PricingRuleService) Create(ctx context.Context, data model.PricingRuleCreate) (string, error) {
	log := pkglog.WithMethod(s.log, "Create")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	id, err := s.ruleRepo.Create(ctx, data)
	if err != nil {
		log.Error("failed to create pricing rule", pkglog.Err(err))
		return "", err
	}

	return id, nil
}

func (s *PricingRuleService) GetByID(ctx context.Context, id string) (model.PricingRule, error) {
	log := pkglog.WithMethod(s.log, "GetByID")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to get pricing rule", pkglog.Err(err))
		}
		return model.PricingRule{}, err
	}

	return rule, nil
}

func (s *PricingRuleService) List(ctx context.Context, filter model.PricingRuleListFilter) ([]model.PricingRule, error) {
	log := pkglog.WithMethod(s.log, "List")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	rules, err := s.ruleRepo.List(ctx, filter)
	if err != nil {
		log.Error("failed to list pricing rules", pkglog.Err(err))
		return nil, err
	}

	return rules, nil
}

func (s *PricingRuleService) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	log := pkglog.WithMethod(s.log, "Update")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	if err := s.ruleRepo.Update(ctx, id, data); err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to update pricing rule", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *PricingRuleService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMethod(s.log, "Delete")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	if err := s.ruleRepo.Delete(ctx, id); err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to delete pricing rule", pkglog.Err(err))
		}
		return err
	}

	return nil
}
