package model

import "time"

type CarTelemetry struct {
	CarID      string
	Location   Location
	FuelLevel  *float32
	MileageKM  int64
	RecordedAt time.Time
}
