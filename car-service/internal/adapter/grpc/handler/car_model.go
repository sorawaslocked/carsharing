package handler

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"
	pkglog "carsharing/shared/pkg/log"

	carsvc "carsharing/protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CarModelHandler struct {
	log             *slog.Logger
	carModelService CarModelService

	carsvc.UnimplementedCarModelServiceServer
}

func NewCarModelHandler(log *slog.Logger, carModelService CarModelService) *CarModelHandler {
	return &CarModelHandler{
		log:             pkglog.WithComponent(log, "grpc.handler.CarModelHandler"),
		carModelService: carModelService,
	}
}

func (h *CarModelHandler) CreateCarModel(ctx context.Context, req *carsvc.CreateCarModelRequest) (*carsvc.CreateCarModelResponse, error) {
	createInput := dto.FromCreateCarModelRequest(req)

	id, err := h.carModelService.Create(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateCarModelResponse{Id: id}, nil
}

func (h *CarModelHandler) GetCarModel(ctx context.Context, req *carsvc.GetCarModelRequest) (*carsvc.GetCarModelResponse, error) {
	carModel, err := h.carModelService.Get(ctx, req.Id)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarModelResponse{CarModel: dto.ToCarModelProto(carModel)}, nil
}

func (h *CarModelHandler) ListCarModels(ctx context.Context, req *carsvc.ListCarModelsRequest) (*carsvc.ListCarModelsResponse, error) {
	filterInput := dto.FromListCarModelsRequest(req)

	carModels, err := h.carModelService.List(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.ListCarModelsResponse{
		CarModels: dto.ToCarModelProtos(carModels),
	}, nil
}

func (h *CarModelHandler) UpdateCarModel(ctx context.Context, req *carsvc.UpdateCarModelRequest) (*emptypb.Empty, error) {
	updateInput := dto.FromUpdateCarModelRequest(req)

	if err := h.carModelService.Update(ctx, req.Id, updateInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarModelHandler) DeleteCarModel(ctx context.Context, req *carsvc.DeleteCarModelRequest) (*emptypb.Empty, error) {
	if err := h.carModelService.Delete(ctx, req.Id); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarModelHandler) GetCarModelImageUploadData(ctx context.Context, _ *emptypb.Empty) (*carsvc.GetCarModelImageUploadDataResponse, error) {
	data, err := h.carModelService.GetImageUploadData(ctx)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarModelImageUploadDataResponse{
		UploadData: dto.ToImageUploadData(data),
	}, nil
}
