package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarStatusLogEntry struct {
	ID         string
	CarID      string
	FromStatus CarStatus
	ToStatus   CarStatus
	ActorType  string
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	RecordedAt time.Time
}

type CarStatusLogFilter struct {
	CarID *string

	Pagination *sharedmodel.Pagination
}

type CarStatusTransitionInput struct {
	CarID     string
	ToStatus  CarStatus
	ActorType string
	ActorID   *string
	Reason    *string
	Metadata  map[string]any
}
