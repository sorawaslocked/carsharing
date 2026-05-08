package model

type Pagination struct {
	Limit  int64
	Offset int64
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type ImageUploadData struct {
	PresignedPutURL string
	ObjectKey       string
}

type PricingSnapshot struct {
	RateTenge         int32
	RatePerKMTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
}
