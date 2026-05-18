package model

// TripLiveFeed is the real-time snapshot streamed to clients during an active trip.
type TripLiveFeed struct {
	ElapsedSeconds     int64
	CurrentCostTenge   int32
	DistanceTraveledKM float64
}
