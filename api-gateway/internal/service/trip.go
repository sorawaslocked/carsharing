package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

type TripService struct {
	presenter TripPresenter
}

func NewTripService(presenter TripPresenter) *TripService {
	return &TripService{presenter: presenter}
}

func (s *TripService) Start(ctx context.Context, bookingID string) (string, error) {
	return s.presenter.Start(ctx, bookingID)
}

func (s *TripService) Get(ctx context.Context, id string) (model.Trip, error) {
	return s.presenter.Get(ctx, id)
}

func (s *TripService) List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error) {
	return s.presenter.List(ctx, filter)
}

func (s *TripService) End(ctx context.Context, id string) error {
	return s.presenter.End(ctx, id)
}

func (s *TripService) Cancel(ctx context.Context, id string, reason *string) error {
	return s.presenter.Cancel(ctx, id, reason)
}

func (s *TripService) GetSummary(ctx context.Context, id string) (model.TripSummary, error) {
	return s.presenter.GetSummary(ctx, id)
}

func (s *TripService) GetStatusHistory(ctx context.Context, id string, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	return s.presenter.GetStatusHistory(ctx, id, filter)
}
