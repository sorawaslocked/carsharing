package dto

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	basecarpb "github.com/sorawaslocked/car-rental-protos/gen/base/car"
)

func CarPricingRuleFromProto(r *basecarpb.CarPricingRule) model.CarPricingRule {
	rule := model.CarPricingRule{
		ID:                 r.GetId(),
		ModelID:            r.ModelId,
		ZoneID:             r.ZoneId,
		Class:              r.Class,
		RatePerMinuteTenge: r.GetRatePerMinuteTenge(),
		RatePerKMTenge:     r.GetRatePerKmTenge(),
		FreeMinutes:        r.GetFreeMinutes(),
		MinChargeTenge:     r.GetMinChargeTenge(),
		IsActive:           r.GetIsActive(),
	}
	if r.GetStartsAt() != nil {
		t := r.GetStartsAt().AsTime()
		rule.StartsAt = &t
	}
	if r.GetExpiresAt() != nil {
		t := r.GetExpiresAt().AsTime()
		rule.ExpiresAt = &t
	}
	if r.GetCreatedAt() != nil {
		rule.CreatedAt = r.GetCreatedAt().AsTime()
	}
	if r.GetUpdatedAt() != nil {
		rule.UpdatedAt = r.GetUpdatedAt().AsTime()
	}
	return rule
}
