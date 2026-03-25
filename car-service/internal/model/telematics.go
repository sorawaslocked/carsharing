package model

import "time"

type CarTelematicsEvent struct {
	ID           string
	CarID        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	OdometerKM   int64
	RecordedAt   time.Time
	ReceivedAt   time.Time
}

type TelematicsUpdate struct {
	CarID        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	OdometerKM   int64
	RecordedAt   time.Time
}
