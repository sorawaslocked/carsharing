package model_test

import (
	"testing"

	"carsharing/booking-service/internal/model"
)

func TestValidateTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    model.BookingStatus
		to      model.BookingStatus
		wantErr error
	}{
		{"created → expired", model.BookingStatusCreated, model.BookingStatusExpired, nil},
		{"created → completed", model.BookingStatusCreated, model.BookingStatusCompleted, nil},
		{"created → cancelled", model.BookingStatusCreated, model.BookingStatusCancelled, nil},
		{"created → created", model.BookingStatusCreated, model.BookingStatusCreated, model.ErrInvalidBookingStatusTransition},
		{"expired → created", model.BookingStatusExpired, model.BookingStatusCreated, model.ErrInvalidBookingStatusTransition},
		{"expired → completed", model.BookingStatusExpired, model.BookingStatusCompleted, model.ErrInvalidBookingStatusTransition},
		{"expired → cancelled", model.BookingStatusExpired, model.BookingStatusCancelled, model.ErrInvalidBookingStatusTransition},
		{"completed → created", model.BookingStatusCompleted, model.BookingStatusCreated, model.ErrInvalidBookingStatusTransition},
		{"completed → expired", model.BookingStatusCompleted, model.BookingStatusExpired, model.ErrInvalidBookingStatusTransition},
		{"completed → cancelled", model.BookingStatusCompleted, model.BookingStatusCancelled, model.ErrInvalidBookingStatusTransition},
		{"cancelled → created", model.BookingStatusCancelled, model.BookingStatusCreated, model.ErrInvalidBookingStatusTransition},
		{"cancelled → expired", model.BookingStatusCancelled, model.BookingStatusExpired, model.ErrInvalidBookingStatusTransition},
		{"cancelled → completed", model.BookingStatusCancelled, model.BookingStatusCompleted, model.ErrInvalidBookingStatusTransition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.ValidateTransition(tt.from, tt.to)
			if err != tt.wantErr {
				t.Errorf("ValidateTransition(%q, %q) = %v, want %v", tt.from, tt.to, err, tt.wantErr)
			}
		})
	}
}

func TestParseBookingStatus(t *testing.T) {
	tests := []struct {
		input  string
		want   model.BookingStatus
		wantOk bool
	}{
		{"created", model.BookingStatusCreated, true},
		{"expired", model.BookingStatusExpired, true},
		{"completed", model.BookingStatusCompleted, true},
		{"cancelled", model.BookingStatusCancelled, true},
		{"", "", false},
		{"CREATED", "", false},
		{"Expired", "", false},
		{"unknown", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := model.ParseBookingStatus(tt.input)
			if ok != tt.wantOk {
				t.Errorf("ParseBookingStatus(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("ParseBookingStatus(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
