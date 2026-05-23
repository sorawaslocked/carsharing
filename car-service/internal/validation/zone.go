package validation

import sharedvalidation "carsharing/shared/validation"

type ZoneFilter struct {
	Type       *string `validate:"omitempty,zonetype"`
	IsActive   *bool   `validate:"omitempty"`
	Pagination *sharedvalidation.Pagination
}

type ZoneCreate struct {
	Name            string `validate:"required,min=1,max=100"`
	Type            string `validate:"required,zonetype"`
	BoundaryGeoJSON string `validate:"required"`
	FeeAdjustment   int32  `validate:"min=-100000,max=100000"`
}

type ZoneUpdate struct {
	Name            *string `validate:"omitempty,min=1,max=100"`
	Type            *string `validate:"omitempty,zonetype"`
	BoundaryGeoJSON *string `validate:"omitempty"`
	FeeAdjustment   *int32  `validate:"omitempty,min=-100000,max=100000"`
	IsActive        *bool   `validate:"omitempty"`
}
