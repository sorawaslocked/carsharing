package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
)

type BookingService interface {
	Create(ctx context.Context, data model.BookingCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.Booking, error)
	List(ctx context.Context, filter model.BookingListFilter) ([]model.Booking, error)
	Cancel(ctx context.Context, id string, reason *string) error
	UpdateStatus(ctx context.Context, id, status string, reason *string) error
	GetStatusHistory(ctx context.Context, filter model.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error)
}

type PricingRuleService interface {
	Create(ctx context.Context, data model.PricingRuleCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.PricingRule, error)
	List(ctx context.Context, filter model.PricingRuleListFilter) ([]model.PricingRule, error)
	Update(ctx context.Context, id string, data model.PricingRuleUpdate) error
	Delete(ctx context.Context, id string) error
}
