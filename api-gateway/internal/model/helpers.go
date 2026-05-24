package model

type PricingSnapshot struct {
	RateTenge         int32
	RatePerKMTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
}
