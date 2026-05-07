package model

import "time"

type CarInsurance struct {
	ID        string
	CarID     string
	Type      string
	Provider  string
	PolicyNum string

	StartsAt  time.Time
	ExpiresAt time.Time

	CostTenge int32
	Status    string

	ImageURLs []string
	Notes     *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarInsuranceFilter struct {
	CarID  *string
	Type   *string
	Status *string

	ExpiringWithinDays *int32

	Pagination *Pagination
}

type CarInsuranceCreate struct {
	CarID     string
	Type      string
	Provider  string
	PolicyNum string
	StartsAt  time.Time
	ExpiresAt time.Time
	CostTenge int32
	Notes     *string
}

type CarInsuranceUpdate struct {
	Provider  *string
	PolicyNum *string
	StartsAt  *time.Time
	ExpiresAt *time.Time
	CostTenge *int32
	Status    *string
	Notes     *string
	ImageKeys []string
}
