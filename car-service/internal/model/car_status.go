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
	ActorType  CarStatusActor
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	ChangedAt  time.Time
}

type CarStatusLogFilter struct {
	CarID *string

	Pagination *sharedmodel.Pagination
}

type CarStatusTransitionInput struct {
	CarID     string
	ToStatus  CarStatus
	ActorType CarStatusActor
	ActorID   *string
	Reason    *string
	Metadata  map[string]any
}
