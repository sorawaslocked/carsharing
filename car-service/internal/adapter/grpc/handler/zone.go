package handler

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"
	pkglog "carsharing/shared/pkg/log"

	carsvc "carsharing/protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ZoneHandler struct {
	log         *slog.Logger
	zoneService ZoneService

	carsvc.UnimplementedZoneServiceServer
}

func NewZoneHandler(log *slog.Logger, zoneService ZoneService) *ZoneHandler {
	return &ZoneHandler{
		log:         pkglog.WithComponent(log, "grpc.handler.ZoneHandler"),
		zoneService: zoneService,
	}
}

func (h *ZoneHandler) CreateZone(ctx context.Context, req *carsvc.CreateZoneRequest) (*carsvc.CreateZoneResponse, error) {
	createInput := dto.FromCreateZoneRequest(req)

	id, err := h.zoneService.Create(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateZoneResponse{Id: id}, nil
}

func (h *ZoneHandler) GetZone(ctx context.Context, req *carsvc.GetZoneRequest) (*carsvc.GetZoneResponse, error) {
	zone, err := h.zoneService.Get(ctx, req.Id)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetZoneResponse{Zone: dto.ToZoneProto(zone)}, nil
}

func (h *ZoneHandler) ListZones(ctx context.Context, req *carsvc.ListZonesRequest) (*carsvc.ListZonesResponse, error) {
	filterInput := dto.FromListZonesRequest(req)

	zones, err := h.zoneService.List(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.ListZonesResponse{Zones: dto.ToZoneProtos(zones)}, nil
}

func (h *ZoneHandler) UpdateZone(ctx context.Context, req *carsvc.UpdateZoneRequest) (*emptypb.Empty, error) {
	updateInput := dto.FromUpdateZoneRequest(req)

	if err := h.zoneService.Update(ctx, req.Id, updateInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ZoneHandler) DeleteZone(ctx context.Context, req *carsvc.DeleteZoneRequest) (*emptypb.Empty, error) {
	if err := h.zoneService.Delete(ctx, req.Id); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}
