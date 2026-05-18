package handler

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CarHandler struct {
	carService CarService

	log *slog.Logger

	carsvc.UnimplementedCarServiceServer
}

func NewCarHandler(carService CarService, log *slog.Logger) *CarHandler {
	s := &CarHandler{
		carService: carService,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "CarHandler"),
		),
	)

	return s
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

	cars, err := h.carService.GetAll(ctx, filterInput)
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
		UploadData: dto.ToImageUploadData(data.URL, data.ObjectKey),
	}, nil
}

func (h *CarHandler) UpdateCarTelemetry(ctx context.Context, req *carsvc.UpdateCarTelemetryRequest) (*emptypb.Empty, error) {
	input := dto.FromUpdateCarTelemetryRequest(req)

	if err := h.carService.UpdateCarTelemetry(ctx, req.Id, input); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarHandler) GetCarStatusHistory(ctx context.Context, req *carsvc.GetCarStatusHistoryRequest) (*carsvc.GetCarStatusHistoryResponse, error) {
	filter := dto.FromGetCarStatusHistoryRequest(req)

	entries, err := h.carService.GetCarStatusHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarStatusHistoryResponse{Readings: dto.ToCarStatusReadingProtos(entries)}, nil
}

func (h *CarHandler) GetCarFuelHistory(ctx context.Context, req *carsvc.GetCarFuelHistoryRequest) (*carsvc.GetCarFuelHistoryResponse, error) {
	filter := dto.FromGetCarFuelHistoryRequest(req)

	readings, err := h.carService.GetCarFuelHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarFuelHistoryResponse{Readings: dto.ToCarFuelReadingProtos(readings)}, nil
}

func (h *CarHandler) GetCarLocationHistory(ctx context.Context, req *carsvc.GetCarLocationHistoryRequest) (*carsvc.GetCarLocationHistoryResponse, error) {
	filter := dto.FromGetTelematicsHistoryRequest(req.CarId, req.From, req.To, req.Pagination)

	events, err := h.carService.GetCarLocationHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarLocationHistoryResponse{Readings: dto.ToCarLocationReadingProtos(events)}, nil
}

func (h *CarHandler) GetCarBatteryHistory(ctx context.Context, req *carsvc.GetCarBatteryHistoryRequest) (*carsvc.GetCarBatteryHistoryResponse, error) {
	filter := dto.FromGetTelematicsHistoryRequest(req.CarId, req.From, req.To, req.Pagination)

	events, err := h.carService.GetCarBatteryHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarBatteryHistoryResponse{Readings: dto.ToCarBatteryReadingProtos(events)}, nil
}

func (h *CarHandler) GetCarMileageHistory(ctx context.Context, req *carsvc.GetCarMileageHistoryRequest) (*carsvc.GetCarMileageHistoryResponse, error) {
	filter := dto.FromGetTelematicsHistoryRequest(req.CarId, req.From, req.To, req.Pagination)

	events, err := h.carService.GetCarMileageHistory(ctx, filter)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarMileageHistoryResponse{Readings: dto.ToCarMileageReadingProtos(events)}, nil
}
