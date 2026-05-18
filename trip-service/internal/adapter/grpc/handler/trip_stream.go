package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"google.golang.org/grpc"

	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"

	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
	pkglog "carsharing/trip-service/internal/pkg/log"
)

type TripStreamService interface {
	StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error
}

type TripStreamHandler struct {
	tripsvc.UnimplementedTripStreamServiceServer
	log     *slog.Logger
	service TripStreamService
}

func NewTripStreamHandler(log *slog.Logger, service TripStreamService) *TripStreamHandler {
	return &TripStreamHandler{
		log:     pkglog.WithComponent(log, "handler.TripStreamHandler"),
		service: service,
	}
}

func (h *TripStreamHandler) StreamTripLiveFeed(req *tripsvc.StreamTripLiveFeedRequest, stream grpc.ServerStreamingServer[tripsvc.StreamTripLiveFeedResponse]) error {
	err := h.service.StreamTripLiveFeed(stream.Context(), req.TripId, func(feed model.TripLiveFeed) error {
		return stream.Send(&tripsvc.StreamTripLiveFeedResponse{
			ElapsedSeconds:     feed.ElapsedSeconds,
			CurrentCostTenge:   feed.CurrentCostTenge,
			DistanceTraveledKm: feed.DistanceTraveledKM,
		})
	})
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return dto.ToStatusError(err)
	}
	return nil
}
