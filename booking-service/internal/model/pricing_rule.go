package model

import "time"

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
	Pagination Pagination
}
