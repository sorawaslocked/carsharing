package validation

import (
	sharedvalidation "carsharing/shared/validation"
)

type BookingCreate struct {
	UserID           string `validate:"required,uuid4"`
	CarID            string `validate:"required,uuid4"`
	CommittedPeriods *int32 `validate:"omitempty,min=1"`
	PricingRuleID    string `validate:"required,uuid4"`
}

type BookingListFilter struct {
	UserID        *string `validate:"omitempty,uuid4"`
	CarID         *string `validate:"omitempty,uuid4"`
	Status        *string `validate:"omitempty,booking_status"`
	PricingRuleID *string `validate:"omitempty,uuid4"`
	Pagination    *sharedvalidation.Pagination
}

type BookingStatusUpdate struct {
	Status string  `validate:"required,booking_status"`
	Reason *string `validate:"omitempty,min=1,max=500"`
}

type BookingStatusHistoryFilter struct {
	BookingID  string `validate:"required,uuid4"`
	Pagination *sharedvalidation.Pagination
}

type idValidation struct {
	ID string `validate:"required,uuid4"`
}

type bookingStatusValidation struct {
	Status string `validate:"required,booking_status"`
}
