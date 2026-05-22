package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

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
	Images []sharedmodel.Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarInsuranceFilter struct {
	CarID  *string
	Type   *InsuranceType
	Status *InsuranceStatus

	ExpiringWithinDays *int32

	Pagination *sharedmodel.Pagination
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
