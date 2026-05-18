package dto

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	basebookingpb "github.com/sorawaslocked/car-rental-protos/gen/base/booking"
)

func BookingFromProto(b *basebookingpb.Booking) model.Booking {
	booking := model.Booking{
		ID:            b.GetId(),
		UserID:        b.GetUserId(),
		CarID:         b.GetCarId(),
		Status:        b.GetStatus(),
		PricingRuleID: b.GetPricingRuleId(),
	}

	if b.CommittedPeriods != nil {
		v := b.GetCommittedPeriods()
		booking.CommittedPeriods = &v
	}

	if s := b.GetPricingSnapshot(); s != nil {
		booking.PricingSnapshot = PricingSnapshotFromProto(s)
	}

	if b.GetExpiresAt() != nil {
		v := b.GetExpiresAt().AsTime()
		booking.ExpiresAt = &v
	}
	if b.GetCreatedAt() != nil {
		booking.CreatedAt = b.GetCreatedAt().AsTime()
	}
	if b.GetUpdatedAt() != nil {
		booking.UpdatedAt = b.GetUpdatedAt().AsTime()
	}

	return booking
}

func BookingStatusReadingFromProto(r *basebookingpb.BookingStatusReading) model.BookingStatusReading {
	reading := model.BookingStatusReading{
		ID:         r.GetId(),
		BookingID:  r.GetBookingId(),
		FromStatus: r.GetFromStatus(),
		ToStatus:   r.GetToStatus(),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
	}

	if r.GetChangedAt() != nil {
		reading.ChangedAt = r.GetChangedAt().AsTime()
	}

	return reading
}
