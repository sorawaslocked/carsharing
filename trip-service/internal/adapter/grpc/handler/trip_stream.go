package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"google.golang.org/grpc"

	tripsvc "carsharing/protos/gen/service/trip"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
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
		log:     pkglog.WithComponent(log, "adapter.grpc.handler.TripStreamHandler"),
		service: service,
	}
}

func (h *TripStreamHandler) StreamTripLiveFeed(req *tripsvc.StreamTripLiveFeedRequest, stream grpc.ServerStreamingServer[tripsvc.StreamTripLiveFeedResponse]) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamTripLiveFeed"), pkgutils.MetadataFromCtx(stream.Context()))

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
		log.Warn("streaming trip live feed", pkglog.Err(err))
		return dto.ToStatusError(err)
	}
	return nil
}
