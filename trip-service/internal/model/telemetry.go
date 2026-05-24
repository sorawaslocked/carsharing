package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarTelemetry struct {
	CarID      string
	Location   sharedmodel.Location
	FuelLevel  *float32
	MileageKM  int64
	RecordedAt time.Time
}
