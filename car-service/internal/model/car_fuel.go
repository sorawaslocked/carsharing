package model

import "time"

type CarFuelState struct {
	CarID         string
	FuelPct       int
	FuelStatus    FuelStatus
	FuelUpdatedAt time.Time
	History       []int
}

type CarFuelAnomalyAlert struct {
	CarID          string
	DropPct        int
	MinutesElapsed float64
	RecordedAt     time.Time
}
