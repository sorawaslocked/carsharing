package validation

import (
	sharedvalidation "carsharing/shared/validation"
)

type PricingRuleCreate struct {
	ModelID           *string `validate:"omitempty,uuid4"`
	ZoneID            *string `validate:"omitempty,uuid4"`
	Class             *string `validate:"omitempty,min=1,max=50"`
	Type              string  `validate:"required,pricing_rule_type"`
	RateTenge         int32   `validate:"required,min=1"`
	RatePerKMTenge    *int32  `validate:"omitempty,min=0"`
	FreeMinutes       *int32  `validate:"omitempty,min=0"`
	MinChargeTenge    *int32  `validate:"omitempty,min=0"`
	OvertimePolicy    *string `validate:"omitempty,min=1,max=100"`
	OvertimeRateTenge *int32  `validate:"omitempty,min=0"`
}

type PricingRuleUpdate struct {
	ModelID           *string `validate:"omitempty,uuid4"`
	ZoneID            *string `validate:"omitempty,uuid4"`
	Class             *string `validate:"omitempty,min=1,max=50"`
	Type              *string `validate:"omitempty,pricing_rule_type"`
	RateTenge         *int32  `validate:"omitempty,min=1"`
	RatePerKMTenge    *int32  `validate:"omitempty,min=0"`
	FreeMinutes       *int32  `validate:"omitempty,min=0"`
	MinChargeTenge    *int32  `validate:"omitempty,min=0"`
	OvertimePolicy    *string `validate:"omitempty,min=1,max=100"`
	OvertimeRateTenge *int32  `validate:"omitempty,min=0"`
	IsActive          *bool
}

type PricingRuleListFilter struct {
	ModelID    *string `validate:"omitempty,uuid4"`
	ZoneID     *string `validate:"omitempty,uuid4"`
	Class      *string `validate:"omitempty,min=1,max=50"`
	Type       *string `validate:"omitempty,pricing_rule_type"`
	IsActive   *bool
	Pagination *sharedvalidation.Pagination
}
