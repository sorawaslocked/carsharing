package model

import sharedmodel "carsharing/shared/model"

import "time"

type PricingRule struct {
	ID string

	ModelID *string
	Class   *string

	Type      string
	RateTenge int32

	RatePerKMTenge *int32
	FreeMinutes    *int32
	MinChargeTenge *int32

	OvertimePolicy    *string
	OvertimeRateTenge *int32

	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PricingRuleFilter struct {
	ModelID  *string
	Class    *string
	Type     *string
	IsActive *bool

	Pagination *sharedmodel.Pagination
}

type PricingRuleCreate struct {
	ModelID *string
	Class   *string

	Type      string
	RateTenge int32

	RatePerKMTenge *int32
	FreeMinutes    *int32
	MinChargeTenge *int32

	OvertimePolicy    *string
	OvertimeRateTenge *int32
}

type PricingRuleUpdate struct {
	ModelID *string
	Class   *string

	Type      *string
	RateTenge *int32

	RatePerKMTenge *int32
	FreeMinutes    *int32
	MinChargeTenge *int32

	OvertimePolicy    *string
	OvertimeRateTenge *int32

	IsActive *bool
}
