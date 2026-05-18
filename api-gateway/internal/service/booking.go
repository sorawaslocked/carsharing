package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type BookingService struct {
	presenter BookingPresenter
}

func NewBookingService(presenter BookingPresenter) *BookingService {
	return &BookingService{presenter: presenter}
}

func (s *BookingService) Create(ctx context.Context, data model.BookingCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *BookingService) Get(ctx context.Context, id string) (model.Booking, error) {
	return s.presenter.Get(ctx, id)
}

func (s *BookingService) List(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error) {
	return s.presenter.List(ctx, filter)
}

func (s *BookingService) Cancel(ctx context.Context, id string) error {
	return s.presenter.Cancel(ctx, id)
}

func (s *BookingService) UpdateStatus(ctx context.Context, id string, data model.BookingStatusUpdate) error {
	return s.presenter.UpdateStatus(ctx, id, data)
}

func (s *BookingService) GetStatusHistory(ctx context.Context, id string, filter model.BookingStatusReadingFilter) ([]model.BookingStatusReading, error) {
	return s.presenter.GetStatusHistory(ctx, id, filter)
}
