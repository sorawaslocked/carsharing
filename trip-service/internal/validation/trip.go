package validation

import (
	sharedvalidation "carsharing/shared/validation"
)

type TripStart struct {
	BookingID string `validate:"required,uuid4"`
}

type TripEnd struct {
	ID string `validate:"required,uuid4"`
}

type TripCancel struct {
	ID     string  `validate:"required,uuid4"`
	Reason *string `validate:"omitempty,min=1,max=500"`
}

type TripFilter struct {
	UserID     *string `validate:"omitempty,uuid4"`
	CarID      *string `validate:"omitempty,uuid4"`
	Status     *string `validate:"omitempty,trip_status"`
	TimeRange  *sharedvalidation.TimeRange
	Pagination *sharedvalidation.Pagination
}

type TripStatusHistoryFilter struct {
	TripID     string `validate:"required,uuid4"`
	TimeRange  *sharedvalidation.TimeRange
	Pagination *sharedvalidation.Pagination
}

type idValidation struct {
	ID string `validate:"required,uuid4"`
}
