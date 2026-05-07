package model

import "time"

type CarPricingRule struct {
	ID string

	// Scope selectors — nil means applies to all values of that dimension
	ModelID *string
	ZoneID  *string
	Class   *string

	RatePerMinuteTenge int32
	RatePerKMTenge     int32
	FreeMinutes        int32
	MinChargeTenge     int32

	IsActive  bool
	StartsAt  *time.Time
	ExpiresAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarPricingRuleFilter struct {
	ModelID  *string
	ZoneID   *string
	Class    *string
	IsActive *bool

	Pagination *Pagination
}

type CarPricingRuleCreate struct {
	ModelID *string
	ZoneID  *string
	Class   *string

	RatePerMinuteTenge int32
	RatePerKMTenge     int32
	FreeMinutes        int32
	MinChargeTenge     int32

	StartsAt  *time.Time
	ExpiresAt *time.Time
}

type CarPricingRuleUpdate struct {
	ModelID *string
	ZoneID  *string
	Class   *string

	RatePerMinuteTenge *int32
	RatePerKMTenge     *int32
	FreeMinutes        *int32
	MinChargeTenge     *int32

	IsActive  *bool
	StartsAt  *time.Time
	ExpiresAt *time.Time
}
