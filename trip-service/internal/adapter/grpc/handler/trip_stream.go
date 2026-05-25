package handler

import (
	"errors"
	"io"
	"log/slog"

	"google.golang.org/grpc"

	tripsvc "carsharing/protos/gen/service/trip"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
)

type TripStreamHandler struct {
	tripsvc.UnimplementedTripStreamServiceServer
	log     *slog.Logger
	service TripService
}

func NewTripStreamHandler(log *slog.Logger, service TripService) *TripStreamHandler {
	return &TripStreamHandler{
		log:     pkglog.WithComponent(log, "adapter.grpc.handler.TripStreamHandler"),
		service: service,
	}
}

func (h *TripStreamHandler) StreamTripLiveFeed(req *tripsvc.StreamTripLiveFeedRequest, stream grpc.ServerStreamingServer[tripsvc.StreamTripLiveFeedResponse]) error {
	ctx := stream.Context()
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamTripLiveFeed"), md)

	trip, err := h.service.GetTrip(ctx, req.TripId)
	if err != nil {
		log.Warn("getting trip for stream auth check", pkglog.Err(err))
		return dto.ToStatusError(err)
	}

	if !hasTripManagerRole(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return dto.ToStatusError(model.ErrInsufficientPermissions)
	}

	err = h.service.StreamTripLiveFeed(ctx, req.TripId, func(feed model.TripLiveFeed) error {
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

func hasTripManagerRole(roles []sharedmodel.Role) bool {
	for _, r := range roles {
		if r == sharedmodel.RoleAdmin || r == sharedmodel.RoleBookingManager {
			return true
		}
	}
	return false
}
