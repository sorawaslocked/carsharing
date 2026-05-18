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
	ZoneID  *string `json:"zoneID,omitempty"`
	Class   *string `json:"class,omitempty"`

	Type      string `json:"type"`
	RateTenge int32  `json:"rateTenge"`

	RatePerKMTenge *int32 `json:"ratePerKMTenge,omitempty"`
	FreeMinutes    *int32 `json:"freeMinutes,omitempty"`
	MinChargeTenge *int32 `json:"minChargeTenge,omitempty"`

	OvertimePolicy    *string `json:"overtimePolicy,omitempty"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge,omitempty"`

	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PricingRuleCreateRequest struct {
	ModelID *string `json:"modelID"`
	ZoneID  *string `json:"zoneID"`
	Class   *string `json:"class"`

	Type      string `json:"type"`
	RateTenge int32  `json:"rateTenge"`

	RatePerKMTenge *int32 `json:"ratePerKMTenge"`
	FreeMinutes    *int32 `json:"freeMinutes"`
	MinChargeTenge *int32 `json:"minChargeTenge"`

	OvertimePolicy    *string `json:"overtimePolicy"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge"`
}

type PricingRuleUpdateRequest struct {
	ModelID *string `json:"modelID"`
	ZoneID  *string `json:"zoneID"`
	Class   *string `json:"class"`

	Type      *string `json:"type"`
	RateTenge *int32  `json:"rateTenge"`

	RatePerKMTenge *int32 `json:"ratePerKMTenge"`
	FreeMinutes    *int32 `json:"freeMinutes"`
	MinChargeTenge *int32 `json:"minChargeTenge"`

	OvertimePolicy    *string `json:"overtimePolicy"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge"`

	IsActive *bool `json:"isActive"`
}

func FromPricingRuleCreateRequest(ctx *gin.Context) (model.PricingRuleCreate, error) {
	var req PricingRuleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.PricingRuleCreate{}, err
	}
	return model.PricingRuleCreate{
		ModelID:           req.ModelID,
		ZoneID:            req.ZoneID,
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
		ZoneID:            req.ZoneID,
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
	if v := ctx.Query("zoneID"); v != "" {
		f.ZoneID = &v
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
		ZoneID:            m.ZoneID,
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
