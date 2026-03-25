package model

type Location struct {
	Latitude  float64
	Longitude float64
}

type LocationFilter struct {
	Location Location
	RadiusKM float64
}

type Pagination struct {
	Limit  *int64
	Offset *int64
}

type PaginationInput struct {
	Limit  *int64 `validate:"omitempty,min=1,max=100"`
	Offset *int64 `validate:"omitempty,min=0"`
}
