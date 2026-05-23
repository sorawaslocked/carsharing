package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type PricingRuleType string

const (
	PricingRuleTypeTime     PricingRuleType = "time"
	PricingRuleTypeDistance PricingRuleType = "distance"
	PricingRuleTypeCombined PricingRuleType = "combined"
)

var validPricingRuleTypes = map[PricingRuleType]struct{}{
	PricingRuleTypeTime:     {},
	PricingRuleTypeDistance: {},
	PricingRuleTypeCombined: {},
}

func PricingRuleTypeFromString(s string) (PricingRuleType, bool) {
	t := PricingRuleType(s)
	_, ok := validPricingRuleTypes[t]
	return t, ok
}

type PricingSnapshot struct {
	RateTenge         int32
	RatePerKMTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
}

type PricingRule struct {
	ID                string
	ModelID           *string
	ZoneID            *string
	Class             *string
	Type              string
	RateTenge         int32
	RatePerKMTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type PricingRuleCreate struct {
	ModelID           *string
	ZoneID            *string
	Class             *string
	Type              string
	RateTenge         int32
	RatePerKMTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
}

type PricingRuleUpdate struct {
	ModelID           *string
	ZoneID            *string
	Class             *string
	Type              *string
	RateTenge         *int32
	RatePerKMTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
	IsActive          *bool
}

type PricingRuleListFilter struct {
	ModelID    *string
	ZoneID     *string
	Class      *string
	Type       *string
	IsActive   *bool
	Pagination sharedmodel.Pagination
}
