package handler

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"
	pkglog "carsharing/shared/pkg/log"

	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CarHandler struct {
	log        *slog.Logger
	carService CarService

	carsvc.UnimplementedCarServiceServer
}

func NewCarHandler(log *slog.Logger, carService CarService) *CarHandler {
	return &CarHandler{
		log:        pkglog.WithComponent(log, "adapter.grpc.handler.CarHandler"),
		carService: carService,
	}
}

func (h *CarHandler) CreateCar(ctx context.Context, req *carsvc.CreateCarRequest) (*carsvc.CreateCarResponse, error) {
	createInput := dto.FromCreateCarRequest(req)

	id, err := h.carService.Create(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateCarResponse{Id: id}, nil
}

func (h *CarHandler) GetCar(ctx context.Context, req *carsvc.GetCarRequest) (*carsvc.GetCarResponse, error) {
	car, err := h.carService.Get(ctx, req.Id)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarResponse{Car: dto.ToCarProto(car)}, nil
}

func (h *CarHandler) ListCars(ctx context.Context, req *carsvc.ListCarsRequest) (*carsvc.ListCarsResponse, error) {
	filterInput := dto.FromListCarsRequest(req)

	cars, err := h.carService.List(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.ListCarsResponse{Cars: dto.ToCarProtos(cars)}, nil
}

func (h *CarHandler) UpdateCar(ctx context.Context, req *carsvc.UpdateCarRequest) (*emptypb.Empty, error) {
	updateInput := dto.FromUpdateCarRequest(req)

	if err := h.carService.Update(ctx, req.Id, updateInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarHandler) UpdateCarStatus(ctx context.Context, req *carsvc.UpdateCarStatusRequest) (*emptypb.Empty, error) {
	statusInput := dto.FromUpdateCarStatusRequest(req)

	if err := h.carService.UpdateCarStatus(ctx, req.Id, statusInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarHandler) UpdateCarTelemetry(ctx context.Context, req *carsvc.UpdateCarTelemetryRequest) (*emptypb.Empty, error) {
	input := dto.FromUpdateCarTelemetryRequest(req)

	if err := h.carService.UpdateCarTelemetry(ctx, req.Id, input); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarHandler) DeleteCar(ctx context.Context, req *carsvc.DeleteCarRequest) (*emptypb.Empty, error) {
	if err := h.carService.Delete(ctx, req.Id); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarHandler) GetCarImageUploadData(ctx context.Context, _ *emptypb.Empty) (*carsvc.GetCarImageUploadDataResponse, error) {
	data, err := h.carService.GetImageUploadData(ctx)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarImageUploadDataResponse{
		UploadData: dto.ToImageUploadData(data),
	}, nil
}

func (h *CarHandler) GetCarStatusHistory(ctx context.Context, req *carsvc.GetCarStatusHistoryRequest) (*carsvc.GetCarStatusHistoryResponse, error) {
	filter := dto.FromGetCarStatusHistoryRequest(req)

	entries, err := h.carService.ListCarStatusHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarStatusHistoryResponse{Readings: dto.ToCarStatusReadingProtos(entries)}, nil
}

func (h *CarHandler) GetCarTelemetryHistory(ctx context.Context, req *carsvc.GetCarTelemetryHistoryRequest) (*carsvc.GetCarTelemetryHistoryResponse, error) {
	filter := dto.FromGetCarTelemetryHistoryRequest(req)

	readings, err := h.carService.ListCarTelemetryHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarTelemetryHistoryResponse{Readings: dto.ToCarTelemetryReadingProtos(readings)}, nil
}
