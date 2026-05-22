package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarTelematicsEvent struct {
	ID           string
	CarID        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	OdometerKM   int64
	ActorID      *string
	ActorType    string
	Metadata     map[string]any
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
	ActorID      *string
	ActorType    string
	RecordedAt   time.Time
}

type CarTelematicsUpdateInput struct {
	MileageKM    int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *Location
}

type TelematicsEventFilter struct {
	CarID *string
	From  *time.Time
	To    *time.Time

	Pagination *sharedmodel.Pagination
}
