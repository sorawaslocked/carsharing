package model

import "time"

type Zone struct {
	ID              string
	Name            string
	Type            string
	BoundaryGeoJSON string
	FeeAdjustment   int32
	IsActive        bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ZoneFilter struct {
	Type     *string
	IsActive *bool
}

type ZoneCreate struct {
	Name            string
	Type            string
	BoundaryGeoJSON string
	FeeAdjustment   int32
}

type ZoneUpdate struct {
	Name            *string
	Type            *string
	BoundaryGeoJSON *string
	FeeAdjustment   *int32
	IsActive        *bool
}
