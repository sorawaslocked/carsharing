package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/dto"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CarModelHandler struct {
	carModelService CarModelService

	log *slog.Logger

	carsvc.UnimplementedCarModelServiceServer
}

func NewCarModelHandler(carModelService CarModelService, log *slog.Logger) *CarModelHandler {
	s := &CarModelHandler{
		carModelService: carModelService,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "CarModelHandler"),
		),
	)

	return s
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

	imageURLs, _ := h.carModelService.GetImageURLs(ctx, req.Id)

	proto := dto.ToCarModelProto(carModel)
	proto.ImageUrls = imageURLs

	return &carsvc.GetCarModelResponse{CarModel: proto}, nil
}

func (h *CarModelHandler) ListCarModels(ctx context.Context, req *carsvc.ListCarModelsRequest) (*carsvc.ListCarModelsResponse, error) {
	filterInput := dto.FromListCarModelsRequest(req)

	carModels, err := h.carModelService.GetAll(ctx, filterInput)
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
		UploadData: dto.ToImageUploadData(data.URL, data.ObjectKey),
	}, nil
}
