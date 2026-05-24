package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarStatusReading struct {
	ID         string
	CarID      string
	FromStatus CarStatus
	ToStatus   CarStatus
	ActorType  sharedmodel.ActorType
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	RecordedAt time.Time
}

type CarStatusReadingFilter struct {
	CarID      string
	FromStatus *CarStatus
	ToStatus   *CarStatus

	Pagination *sharedmodel.Pagination
}

type CarStatusTransition struct {
	CarID     string
	ToStatus  CarStatus
	ActorType sharedmodel.ActorType
	ActorID   *string
	Reason    *string
	Metadata  map[string]any
}
