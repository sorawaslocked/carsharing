package model

type Location struct {
	Latitude  float64
	Longitude float64
}

type LocationFilter struct {
	Location Location
	RadiusKM float64
}
