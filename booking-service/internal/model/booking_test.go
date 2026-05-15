package model_test

import (
	"testing"

	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
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
		{"created → created", model.BookingStatusCreated, model.BookingStatusCreated, model.ErrInvalidTransition},
		{"expired → created", model.BookingStatusExpired, model.BookingStatusCreated, model.ErrInvalidTransition},
		{"expired → completed", model.BookingStatusExpired, model.BookingStatusCompleted, model.ErrInvalidTransition},
		{"expired → cancelled", model.BookingStatusExpired, model.BookingStatusCancelled, model.ErrInvalidTransition},
		{"completed → created", model.BookingStatusCompleted, model.BookingStatusCreated, model.ErrInvalidTransition},
		{"completed → expired", model.BookingStatusCompleted, model.BookingStatusExpired, model.ErrInvalidTransition},
		{"completed → cancelled", model.BookingStatusCompleted, model.BookingStatusCancelled, model.ErrInvalidTransition},
		{"cancelled → created", model.BookingStatusCancelled, model.BookingStatusCreated, model.ErrInvalidTransition},
		{"cancelled → expired", model.BookingStatusCancelled, model.BookingStatusExpired, model.ErrInvalidTransition},
		{"cancelled → completed", model.BookingStatusCancelled, model.BookingStatusCompleted, model.ErrInvalidTransition},
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
		input   string
		want    model.BookingStatus
		wantErr error
	}{
		{"created", model.BookingStatusCreated, nil},
		{"expired", model.BookingStatusExpired, nil},
		{"completed", model.BookingStatusCompleted, nil},
		{"cancelled", model.BookingStatusCancelled, nil},
		{"", "", model.ErrInvalidStatus},
		{"CREATED", "", model.ErrInvalidStatus},
		{"Expired", "", model.ErrInvalidStatus},
		{"unknown", "", model.ErrInvalidStatus},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := model.ParseBookingStatus(tt.input)
			if err != tt.wantErr {
				t.Errorf("ParseBookingStatus(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseBookingStatus(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
