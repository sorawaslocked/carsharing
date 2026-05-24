package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	carsvc "carsharing/protos/gen/service/car"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type ZoneHandler struct {
	client carsvc.ZoneServiceClient
	log    *slog.Logger
}

func NewZoneHandler(client carsvc.ZoneServiceClient, logger *slog.Logger) *ZoneHandler {
	return &ZoneHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.ZoneHandler"),
	}
}

func (h *ZoneHandler) Create(ctx context.Context, data model.ZoneCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.CreateZone(ctx, &carsvc.CreateZoneRequest{
		Name:            data.Name,
		Type:            data.Type,
		BoundaryGeoJson: data.BoundaryGeoJSON,
		FeeAdjustment:   data.FeeAdjustment,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *ZoneHandler) Get(ctx context.Context, id string) (model.Zone, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetZone(ctx, &carsvc.GetZoneRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.Zone{}, dto.FromGrpcErr(err)
	}

	return dto.ZoneFromProto(res.GetZone()), nil
}

func (h *ZoneHandler) List(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.ListZones(ctx, &carsvc.ListZonesRequest{
		Type:     filter.Type,
		IsActive: filter.IsActive,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	zones := make([]model.Zone, len(res.GetZones()))
	for i, z := range res.GetZones() {
		zones[i] = dto.ZoneFromProto(z)
	}

	return zones, nil
}

func (h *ZoneHandler) Update(ctx context.Context, id string, data model.ZoneUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.UpdateZone(ctx, &carsvc.UpdateZoneRequest{
		Id:              id,
		Name:            data.Name,
		Type:            data.Type,
		BoundaryGeoJson: data.BoundaryGeoJSON,
		FeeAdjustment:   data.FeeAdjustment,
		IsActive:        data.IsActive,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *ZoneHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.DeleteZone(ctx, &carsvc.DeleteZoneRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}
