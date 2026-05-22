package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type Zone struct {
	ID              string
	Name            string
	Type            ZoneType
	BoundaryGeoJSON string
	FeeAdjustment   int32
	IsActive        bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ZoneFilter struct {
	Type     *ZoneType
	IsActive *bool

	Pagination *sharedmodel.Pagination
}

type ZoneUpdate struct {
	Name            *string
	Type            *ZoneType
	BoundaryGeoJSON *string
	FeeAdjustment   *int32
	IsActive        *bool
	UpdatedAt       time.Time
}
