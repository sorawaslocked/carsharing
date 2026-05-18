package dto

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	basebookingpb "github.com/sorawaslocked/car-rental-protos/gen/base/booking"
)

func PricingRuleFromProto(r *basebookingpb.PricingRule) model.PricingRule {
	rule := model.PricingRule{
		ID:                r.GetId(),
		ZoneID:            r.ZoneId,
		Type:              r.GetType(),
		RateTenge:         r.GetRateTenge(),
		RatePerKMTenge:    r.RatePerKmTenge,
		FreeMinutes:       r.FreeMinutes,
		MinChargeTenge:    r.MinChargeTenge,
		OvertimePolicy:    r.OvertimePolicy,
		OvertimeRateTenge: r.OvertimeRateTenge,
		IsActive:          r.GetIsActive(),
	}

	if r.CarModelId != nil {
		rule.ModelID = r.CarModelId
	}
	if r.CarClass != nil {
		rule.Class = r.CarClass
	}

	if r.GetCreatedAt() != nil {
		rule.CreatedAt = r.GetCreatedAt().AsTime()
	}
	if r.GetUpdatedAt() != nil {
		rule.UpdatedAt = r.GetUpdatedAt().AsTime()
	}

	return rule
}
