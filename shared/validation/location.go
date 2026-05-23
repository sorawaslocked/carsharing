package validation

type Location struct {
	Latitude  float64 `validate:"latitude_range"`
	Longitude float64 `validate:"longitude_range"`
}

type LocationFilter struct {
	Location Location `validate:""`
	RadiusKM float64  `validate:"radius_range"`
}
