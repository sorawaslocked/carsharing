package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

func (s *TripService) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	return s.presenter.StreamTripLiveFeed(ctx, tripID, send)
}
