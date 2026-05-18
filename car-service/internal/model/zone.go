package model

import "time"

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

	Pagination
}

type ZoneUpdate struct {
	Name            *string
	Type            *ZoneType
	BoundaryGeoJSON *string
	FeeAdjustment   *int32
	IsActive        *bool
	UpdatedAt       time.Time
}

type ZoneFilterInput struct {
	Type     *string `validate:"omitempty,zonetype"`
	IsActive *bool   `validate:"omitempty"`

	PaginationInput
}

type ZoneCreateInput struct {
	Name            string `validate:"required,min=1,max=100"`
	Type            string `validate:"required,zonetype"`
	BoundaryGeoJSON string `validate:"required"`
	FeeAdjustment   int32  `validate:"min=-100000,max=100000"`
}

type ZoneUpdateInput struct {
	Name            *string `validate:"omitempty,min=1,max=100"`
	Type            *string `validate:"omitempty,zonetype"`
	BoundaryGeoJSON *string `validate:"omitempty"`
	FeeAdjustment   *int32  `validate:"omitempty,min=-100000,max=100000"`
	IsActive        *bool   `validate:"omitempty"`
}
