package model

import "time"

type CarInsurance struct {
	ID        string
	CarID     string
	Type      InsuranceType
	Provider  string
	PolicyNum string

	StartsAt  time.Time
	ExpiresAt time.Time

	CostTenge int32
	Status    InsuranceStatus

	Notes  *string
	Images []Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarInsuranceFilter struct {
	CarID  *string
	Type   *InsuranceType
	Status *InsuranceStatus

	ExpiringWithinDays *int32

	Pagination
}

type CarInsuranceUpdate struct {
	Provider  *string
	PolicyNum *string
	StartsAt  *time.Time
	ExpiresAt *time.Time
	CostTenge *int32
	Status    *InsuranceStatus
	Notes     *string
	ImageKeys []string
	UpdatedAt time.Time
}

type CarInsuranceFilterInput struct {
	CarID  *string `validate:"omitempty,uuid"`
	Type   *string `validate:"omitempty,insurancetype"`
	Status *string `validate:"omitempty,insurancestatus"`

	ExpiringWithinDays *int32 `validate:"omitempty,min=1,max=365"`

	PaginationInput
}

type CarInsuranceCreateInput struct {
	CarID     string    `validate:"required,uuid"`
	Type      string    `validate:"required,insurancetype"`
	Provider  string    `validate:"required,min=1,max=100"`
	PolicyNum string    `validate:"required,min=1,max=100"`
	StartsAt  time.Time `validate:"required"`
	ExpiresAt time.Time `validate:"required,gtfield=StartsAt"`
	CostTenge int32     `validate:"min=0"`
	Notes     *string   `validate:"omitempty,min=1,max=500"`
}

type CarInsuranceUpdateInput struct {
	Provider  *string    `validate:"omitempty,min=1,max=100"`
	PolicyNum *string    `validate:"omitempty,min=1,max=100"`
	StartsAt  *time.Time `validate:"omitempty"`
	ExpiresAt *time.Time `validate:"omitempty"`
	CostTenge *int32     `validate:"omitempty,min=0"`
	Status    *string    `validate:"omitempty,insurancestatus"`
	Notes     *string    `validate:"omitempty,min=1,max=500"`
	ImageKeys []string   `validate:"omitempty,max=10,dive,min=1"`
}
