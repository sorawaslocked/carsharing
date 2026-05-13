package dto

type SlimCarLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SlimCar struct {
	ID           string          `json:"id"`
	ModelID      string          `json:"modelId"`
	LicensePlate string          `json:"licensePlate"`
	Color        string          `json:"color"`
	Location     SlimCarLocation `json:"location"`
	FuelLevel    float32         `json:"fuelLevel"`
	Status       string          `json:"status"`
}

type CarFleetMessage struct {
	Cars []SlimCar `json:"cars"`
}

type CarTelemetryMessage struct {
	FuelLevel    float32         `json:"fuelLevel"`
	BatteryLevel float32         `json:"batteryLevel"`
	MileageKM    int64           `json:"mileageKm"`
	Location     SlimCarLocation `json:"location"`
	RecordedAt   string          `json:"recordedAt"`
}

type CarStatusMessage struct {
	CarID      string `json:"carId"`
	FromStatus string `json:"fromStatus"`
	ToStatus   string `json:"toStatus"`
}
