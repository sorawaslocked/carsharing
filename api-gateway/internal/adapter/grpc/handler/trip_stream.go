package handler

import (
	"context"
	"errors"
	"io"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/utils"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"
)

func (h *TripHandler) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	logger := pkglog.WithMethod(h.log, "StreamTripLiveFeed")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	stream, err := h.streamClient.StreamTripLiveFeed(ctx, &tripsvc.StreamTripLiveFeedRequest{TripId: tripID})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}
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
			if dto.IsSystemErr(err) {
				logger.Error("stream recv failed", pkglog.Err(err))
			}
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
