package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type InsuranceType string

const (
	InsuranceTypeOSAGO InsuranceType = "osago"
	InsuranceTypeKASKO InsuranceType = "kasko"
)

var validInsuranceTypes = map[InsuranceType]struct{}{
	InsuranceTypeOSAGO: {},
	InsuranceTypeKASKO: {},
}

func InsuranceTypeFromString(s string) (InsuranceType, bool) {
	it := InsuranceType(s)
	if _, ok := validInsuranceTypes[it]; !ok {
		return "", false
	}
	return it, true
}

func (t InsuranceType) String() string {
	return string(t)
}

type InsuranceStatus string

const (
	InsuranceStatusActive    InsuranceStatus = "active"
	InsuranceStatusExpired   InsuranceStatus = "expired"
	InsuranceStatusCancelled InsuranceStatus = "cancelled"
)

var validInsuranceStatuses = map[InsuranceStatus]struct{}{
	InsuranceStatusActive:    {},
	InsuranceStatusExpired:   {},
	InsuranceStatusCancelled: {},
}

func InsuranceStatusFromString(s string) (InsuranceStatus, bool) {
	is := InsuranceStatus(s)
	if _, ok := validInsuranceStatuses[is]; !ok {
		return "", false
	}
	return is, true
}

func (s InsuranceStatus) String() string {
	return string(s)
}

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
