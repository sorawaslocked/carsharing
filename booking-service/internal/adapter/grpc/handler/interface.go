package handler

import (
	"context"

	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/validation"
)

type BookingService interface {
	Create(ctx context.Context, data validation.BookingCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.Booking, error)
	List(ctx context.Context, filter validation.BookingListFilter) ([]model.Booking, error)
	Cancel(ctx context.Context, id string, reason *string) error
	UpdateStatus(ctx context.Context, id string, data validation.BookingStatusUpdate) error
	GetStatusHistory(ctx context.Context, filter validation.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error)
}

type PricingRuleService interface {
	Create(ctx context.Context, data validation.PricingRuleCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.PricingRule, error)
	List(ctx context.Context, filter validation.PricingRuleListFilter) ([]model.PricingRule, error)
	Update(ctx context.Context, id string, data validation.PricingRuleUpdate) error
	Delete(ctx context.Context, id string) error
}
