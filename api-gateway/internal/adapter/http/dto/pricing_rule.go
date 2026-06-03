package dto

import (
	"time"

	"carsharing/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
)

type PricingRuleResponse struct {
	Rule PricingRule `json:"rule"`
}

type PricingRulesResponse struct {
	Rules []PricingRule `json:"rules"`
}

type PricingRule struct {
	ID      string  `json:"id"`
	ModelID *string `json:"modelID,omitempty"`
	Class   *string `json:"class,omitempty" validate:"omitempty,oneof=economy compact comfort business luxury"`

	Type      string `json:"type" validate:"oneof=by_minute by_hour by_day"`
	RateTenge int32  `json:"rateTenge" validate:"min=1"`

	RatePerKMTenge *int32 `json:"ratePerKMTenge,omitempty" validate:"omitempty,min=0"`
	FreeMinutes    *int32 `json:"freeMinutes,omitempty" validate:"omitempty,min=0"`
	MinChargeTenge *int32 `json:"minChargeTenge,omitempty" validate:"omitempty,min=0"`

	OvertimePolicy    *string `json:"overtimePolicy,omitempty"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge,omitempty" validate:"omitempty,min=0"`

	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PricingRuleCreateRequest struct {
	ModelID *string `json:"modelID"`
	Class   *string `json:"class" validate:"omitempty,oneof=economy compact comfort business luxury"`

	Type      string `json:"type" binding:"required,oneof=by_minute by_hour by_day"`
	RateTenge int32  `json:"rateTenge" binding:"required,min=1"`

	RatePerKMTenge *int32 `json:"ratePerKMTenge" validate:"omitempty,min=0"`
	FreeMinutes    *int32 `json:"freeMinutes" validate:"omitempty,min=0"`
	MinChargeTenge *int32 `json:"minChargeTenge" validate:"omitempty,min=0"`

	OvertimePolicy    *string `json:"overtimePolicy"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge" validate:"omitempty,min=0"`
}

type PricingRuleUpdateRequest struct {
	ModelID *string `json:"modelID"`
	Class   *string `json:"class" validate:"omitempty,oneof=economy compact comfort business luxury"`

	Type      *string `json:"type" validate:"omitempty,oneof=by_minute by_hour by_day"`
	RateTenge *int32  `json:"rateTenge" validate:"omitempty,min=1"`

	RatePerKMTenge *int32 `json:"ratePerKMTenge" validate:"omitempty,min=0"`
	FreeMinutes    *int32 `json:"freeMinutes" validate:"omitempty,min=0"`
	MinChargeTenge *int32 `json:"minChargeTenge" validate:"omitempty,min=0"`

	OvertimePolicy    *string `json:"overtimePolicy"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge" validate:"omitempty,min=0"`

	IsActive *bool `json:"isActive"`
}

func FromPricingRuleCreateRequest(ctx *gin.Context) (model.PricingRuleCreate, error) {
	var req PricingRuleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.PricingRuleCreate{}, err
	}
	return model.PricingRuleCreate{
		ModelID:           req.ModelID,
		Class:             req.Class,
		Type:              req.Type,
		RateTenge:         req.RateTenge,
		RatePerKMTenge:    req.RatePerKMTenge,
		FreeMinutes:       req.FreeMinutes,
		MinChargeTenge:    req.MinChargeTenge,
		OvertimePolicy:    req.OvertimePolicy,
		OvertimeRateTenge: req.OvertimeRateTenge,
	}, nil
}

func FromPricingRuleUpdateRequest(ctx *gin.Context) (model.PricingRuleUpdate, error) {
	var req PricingRuleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.PricingRuleUpdate{}, err
	}
	return model.PricingRuleUpdate{
		ModelID:           req.ModelID,
		Class:             req.Class,
		Type:              req.Type,
		RateTenge:         req.RateTenge,
		RatePerKMTenge:    req.RatePerKMTenge,
		FreeMinutes:       req.FreeMinutes,
		MinChargeTenge:    req.MinChargeTenge,
		OvertimePolicy:    req.OvertimePolicy,
		OvertimeRateTenge: req.OvertimeRateTenge,
		IsActive:          req.IsActive,
	}, nil
}

func PricingRuleFilterFromCtx(ctx *gin.Context) (model.PricingRuleFilter, error) {
	f := model.PricingRuleFilter{}

	if v := ctx.Query("modelID"); v != "" {
		f.ModelID = &v
	}
	if v := ctx.Query("class"); v != "" {
		f.Class = &v
	}
	if v := ctx.Query("type"); v != "" {
		f.Type = &v
	}
	if v := ctx.Query("isActive"); v != "" {
		b := v == "true"
		f.IsActive = &b
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.PricingRuleFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func ToPricingRuleResponse(m model.PricingRule) PricingRule {
	return PricingRule{
		ID:                m.ID,
		ModelID:           m.ModelID,
		Class:             m.Class,
		Type:              m.Type,
		RateTenge:         m.RateTenge,
		RatePerKMTenge:    m.RatePerKMTenge,
		FreeMinutes:       m.FreeMinutes,
		MinChargeTenge:    m.MinChargeTenge,
		OvertimePolicy:    m.OvertimePolicy,
		OvertimeRateTenge: m.OvertimeRateTenge,
		IsActive:          m.IsActive,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
