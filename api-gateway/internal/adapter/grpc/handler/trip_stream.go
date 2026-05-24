package handler

import (
	"context"
	"errors"
	"io"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	tripsvc "carsharing/protos/gen/service/trip"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (h *TripHandler) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamTripLiveFeed"), utils.MetadataFromCtx(ctx))

	stream, err := h.streamClient.StreamTripLiveFeed(ctx, &tripsvc.StreamTripLiveFeedRequest{TripId: tripID})
	if err != nil {
		log.Warn("streaming trip live feed", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Warn("receiving trip live feed stream", pkglog.Err(err))

			return dto.FromGrpcErr(err)
		}

		if err = send(model.TripLiveFeed{
			ElapsedSeconds:     msg.GetElapsedSeconds(),
			CurrentCostTenge:   msg.GetCurrentCostTenge(),
			DistanceTraveledKM: msg.GetDistanceTraveledKm(),
		}); err != nil {
			return err
		}
	}
}
