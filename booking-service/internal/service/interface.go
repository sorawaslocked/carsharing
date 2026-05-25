package service

import (
	"context"
	"time"

	"carsharing/booking-service/internal/model"
	sharedmodel "carsharing/shared/model"
)

type BookingRepository interface {
	Create(ctx context.Context, data model.BookingCreate, expiresAt time.Time) (string, error)
	GetByID(ctx context.Context, id string) (model.Booking, error)
	List(ctx context.Context, filter model.BookingListFilter) ([]model.Booking, error)
	ListCreatedExpired(ctx context.Context, now time.Time) ([]model.Booking, error)
	UpdateStatus(ctx context.Context, id string, status model.BookingStatus, actorType sharedmodel.ActorType, actorID *string, reason *string) error
	GetStatusHistory(ctx context.Context, filter model.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error)
}

type PricingRuleRepository interface {
	Create(ctx context.Context, data model.PricingRuleCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.PricingRule, error)
	List(ctx context.Context, filter model.PricingRuleListFilter) ([]model.PricingRule, error)
	Update(ctx context.Context, id string, data model.PricingRuleUpdate) error
	Delete(ctx context.Context, id string) error
}

type EventPublisher interface {
	PublishBookingCreated(ctx context.Context, booking model.Booking) error
	PublishBookingCancelled(ctx context.Context, booking model.Booking, reason string) error
	PublishBookingExpired(ctx context.Context, booking model.Booking) error
	PublishBookingCompleted(ctx context.Context, booking model.Booking) error
}

type CarChecker interface {
	GetStatus(ctx context.Context, carID string) (model.CarStatus, error)
}

type CarModelChecker interface {
	Exists(ctx context.Context, modelID string) (bool, error)
}

type ZoneChecker interface {
	Exists(ctx context.Context, zoneID string) (bool, error)
}
