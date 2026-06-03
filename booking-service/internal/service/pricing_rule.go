package service

import (
	"context"
	"errors"
	"log/slog"

	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	sharedvalidation "carsharing/shared/validation"

	"github.com/go-playground/validator/v10"
)

type PricingRuleService struct {
	log             *slog.Logger
	validate        *validator.Validate
	ruleRepo        PricingRuleRepository
	carModelChecker CarModelChecker
}

func NewPricingRuleService(
	log *slog.Logger,
	validate *validator.Validate,
	ruleRepo PricingRuleRepository,
	carModelChecker CarModelChecker,
) *PricingRuleService {
	return &PricingRuleService{
		log:             pkglog.WithComponent(log, "service.PricingRuleService"),
		validate:        validate,
		ruleRepo:        ruleRepo,
		carModelChecker: carModelChecker,
	}
}

func (s *PricingRuleService) Create(ctx context.Context, data validation.PricingRuleCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	if data.ModelID != nil {
		exists, err := s.carModelChecker.Exists(ctx, *data.ModelID)
		if err != nil {
			log.Error("checking car model existence", pkglog.Err(err))
			return "", err
		}
		if !exists {
			return "", model.ErrCarModelNotFound
		}
	}

	ruleData := model.PricingRuleCreate{
		ModelID:           data.ModelID,
		Class:             data.Class,
		Type:              data.Type,
		RateTenge:         data.RateTenge,
		RatePerKMTenge:    data.RatePerKMTenge,
		FreeMinutes:       data.FreeMinutes,
		MinChargeTenge:    data.MinChargeTenge,
		OvertimePolicy:    data.OvertimePolicy,
		OvertimeRateTenge: data.OvertimeRateTenge,
	}

	id, err := s.ruleRepo.Create(ctx, ruleData)
	if err != nil {
		log.Error("repo: creating pricing rule", pkglog.Err(err))
		return "", err
	}

	return id, nil
}

func (s *PricingRuleService) GetByID(ctx context.Context, id string) (model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetByID"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.PricingRule{}, err
	}

	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrPricingRuleNotFound) {
			log.Error("repo: getting pricing rule", pkglog.Err(err))
		}
		return model.PricingRule{}, err
	}

	return rule, nil
}

func (s *PricingRuleService) List(ctx context.Context, filter validation.PricingRuleListFilter) ([]model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	rules, err := s.ruleRepo.List(ctx, pricingRuleListFilter(filter))
	if err != nil {
		log.Error("repo: listing pricing rules", pkglog.Err(err))
		return nil, err
	}

	return rules, nil
}

func (s *PricingRuleService) Update(ctx context.Context, id string, data validation.PricingRuleUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	if data.ModelID != nil {
		exists, err := s.carModelChecker.Exists(ctx, *data.ModelID)
		if err != nil {
			log.Error("checking car model existence", pkglog.Err(err))
			return err
		}
		if !exists {
			return model.ErrCarModelNotFound
		}
	}

	repoUpdate := model.PricingRuleUpdate{
		ModelID:           data.ModelID,
		Class:             data.Class,
		Type:              data.Type,
		RateTenge:         data.RateTenge,
		RatePerKMTenge:    data.RatePerKMTenge,
		FreeMinutes:       data.FreeMinutes,
		MinChargeTenge:    data.MinChargeTenge,
		OvertimePolicy:    data.OvertimePolicy,
		OvertimeRateTenge: data.OvertimeRateTenge,
		IsActive:          data.IsActive,
	}

	if err := s.ruleRepo.Update(ctx, id, repoUpdate); err != nil {
		if !errors.Is(err, model.ErrPricingRuleNotFound) {
			log.Error("repo: updating pricing rule", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *PricingRuleService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.ruleRepo.Delete(ctx, id); err != nil {
		if !errors.Is(err, model.ErrPricingRuleNotFound) {
			log.Error("repo: deleting pricing rule", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func pricingRuleListFilter(f validation.PricingRuleListFilter) model.PricingRuleListFilter {
	if f.Pagination == nil {
		f.Pagination = sharedvalidation.DefaultPagination()
	}
	return model.PricingRuleListFilter{
		ModelID:    f.ModelID,
		Class:      f.Class,
		Type:       f.Type,
		IsActive:   f.IsActive,
		Pagination: sharedmodel.Pagination{Limit: f.Pagination.Limit, Offset: f.Pagination.Offset},
	}
}
