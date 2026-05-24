package dto

import (
	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/validation"
	basebookingpb "carsharing/protos/gen/base/booking"
	servicebookingpb "carsharing/protos/gen/service/booking"
	sharedvalidation "carsharing/shared/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PricingRuleToProto(r model.PricingRule) *basebookingpb.PricingRule {
	return &basebookingpb.PricingRule{
		Id:                r.ID,
		CarModelId:        r.ModelID,
		ZoneId:            r.ZoneID,
		CarClass:          r.Class,
		Type:              r.Type,
		RateTenge:         r.RateTenge,
		RatePerKmTenge:    r.RatePerKMTenge,
		FreeMinutes:       r.FreeMinutes,
		MinChargeTenge:    r.MinChargeTenge,
		OvertimePolicy:    r.OvertimePolicy,
		OvertimeRateTenge: r.OvertimeRateTenge,
		IsActive:          r.IsActive,
		CreatedAt:         timestamppb.New(r.CreatedAt),
		UpdatedAt:         timestamppb.New(r.UpdatedAt),
	}
}

func PricingRuleCreateFromProto(req *servicebookingpb.CreatePricingRuleRequest) validation.PricingRuleCreate {
	return validation.PricingRuleCreate{
		ModelID:           req.ModelId,
		ZoneID:            req.ZoneId,
		Class:             req.Class,
		Type:              req.Type,
		RateTenge:         req.RateTenge,
		RatePerKMTenge:    req.RatePerKmTenge,
		FreeMinutes:       req.FreeMinutes,
		MinChargeTenge:    req.MinChargeTenge,
		OvertimePolicy:    req.OvertimePolicy,
		OvertimeRateTenge: req.OvertimeRateTenge,
	}
}

func PricingRuleUpdateFromProto(req *servicebookingpb.UpdatePricingRuleRequest) validation.PricingRuleUpdate {
	return validation.PricingRuleUpdate{
		ModelID:           req.ModelId,
		ZoneID:            req.ZoneId,
		Class:             req.Class,
		Type:              req.Type,
		RateTenge:         req.RateTenge,
		RatePerKMTenge:    req.RatePerKmTenge,
		FreeMinutes:       req.FreeMinutes,
		MinChargeTenge:    req.MinChargeTenge,
		OvertimePolicy:    req.OvertimePolicy,
		OvertimeRateTenge: req.OvertimeRateTenge,
		IsActive:          req.IsActive,
	}
}

func PricingRuleListFilterFromProto(req *servicebookingpb.ListPricingRulesRequest) validation.PricingRuleListFilter {
	filter := validation.PricingRuleListFilter{
		ModelID:  req.ModelId,
		ZoneID:   req.ZoneId,
		Class:    req.Class,
		Type:     req.Type,
		IsActive: req.IsActive,
	}
	if req.Pagination != nil {
		p := sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
		filter.Pagination = &p
	}
	return filter
}
