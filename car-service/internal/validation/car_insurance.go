package validation

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarInsuranceFilter struct {
	CarID              *string `validate:"omitempty,uuid"`
	Type               *string `validate:"omitempty,insurancetype"`
	Status             *string `validate:"omitempty,insurancestatus"`
	ExpiringWithinDays *int32  `validate:"omitempty,min=1,max=365"`
	Pagination         *sharedmodel.Pagination
}

type CarInsuranceCreate struct {
	CarID     string    `validate:"required,uuid"`
	Type      string    `validate:"required,insurancetype"`
	Provider  string    `validate:"required,min=1,max=100"`
	PolicyNum string    `validate:"required,min=1,max=100"`
	StartsAt  time.Time `validate:"required"`
	ExpiresAt time.Time `validate:"required,gtfield=StartsAt"`
	CostTenge int32     `validate:"min=0"`
	Notes     *string   `validate:"omitempty,min=1,max=500"`
}

type CarInsuranceUpdate struct {
	Provider  *string    `validate:"omitempty,min=1,max=100"`
	PolicyNum *string    `validate:"omitempty,min=1,max=100"`
	StartsAt  *time.Time `validate:"omitempty"`
	ExpiresAt *time.Time `validate:"omitempty"`
	CostTenge *int32     `validate:"omitempty,min=0"`
	Status    *string    `validate:"omitempty,insurancestatus"`
	Notes     *string    `validate:"omitempty,min=1,max=500"`
	ImageKeys []string   `validate:"omitempty,max=10,dive,min=1"`
}
