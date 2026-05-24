package wsdto

type TripLiveFeedMessage struct {
	ElapsedSeconds     int64   `json:"elapsedSeconds"`
	CurrentCostTenge   int32   `json:"currentCostTenge"`
	DistanceTraveledKM float64 `json:"distanceTraveledKm"`
}
