package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CarModelHandler struct {
	client carsvc.CarModelServiceClient
	log    *slog.Logger
}

func NewCarModelHandler(client carsvc.CarModelServiceClient, logger *slog.Logger) *CarModelHandler {
	return &CarModelHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.CarModelHandler"),
	}
}

func (h *CarModelHandler) Create(ctx context.Context, data model.CarModelCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")

	res, err := h.client.CreateCarModel(ctx, &carsvc.CreateCarModelRequest{
		Brand:        data.Brand,
		Model:        data.Model,
		Year:         int32(data.Year),
		FuelType:     data.FuelType,
		Transmission: data.Transmission,
		BodyType:     data.BodyType,
		Class:        data.Class,
		Seats:        int32(data.Seats),
		EngineVolume: data.EngineVolume,
		RangeKm:      data.RangeKM,
		Features:     data.Features,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *CarModelHandler) Get(ctx context.Context, id string) (model.CarModel, error) {
	logger := pkglog.WithMethod(h.log, "Get")

	res, err := h.client.GetCarModel(ctx, &carsvc.GetCarModelRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.CarModel{}, dto.FromGrpcErr(err)
	}

	return dto.CarModelFromProto(res.GetCarModel()), nil
}

func (h *CarModelHandler) List(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	logger := pkglog.WithMethod(h.log, "List")

	req := &carsvc.ListCarModelsRequest{
		Brand:        filter.Brand,
		Model:        filter.Model,
		FuelType:     filter.FuelType,
		Transmission: filter.Transmission,
		BodyType:     filter.BodyType,
		Class:        filter.Class,
	}
	if filter.MinSeats != nil {
		v := int32(*filter.MinSeats)
		req.MinSeats = &v
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListCarModels(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	models := make([]model.CarModel, len(res.GetCarModels()))
	for i, m := range res.GetCarModels() {
		models[i] = dto.CarModelFromProto(m)
	}

	return models, nil
}

func (h *CarModelHandler) Update(ctx context.Context, id string, data model.CarModelUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")

	req := &carsvc.UpdateCarModelRequest{
		Id:           id,
		Brand:        data.Brand,
		Model:        data.Model,
		FuelType:     data.FuelType,
		Transmission: data.Transmission,
		BodyType:     data.BodyType,
		Class:        data.Class,
		EngineVolume: data.EngineVolume,
		RangeKm:      data.RangeKM,
		Features:     data.Features,
		ImageKeys:    data.ImageKeys,
	}
	if data.Year != nil {
		v := int32(*data.Year)
		req.Year = &v
	}
	if data.Seats != nil {
		v := int32(*data.Seats)
		req.Seats = &v
	}

	_, err := h.client.UpdateCarModel(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarModelHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")

	_, err := h.client.DeleteCarModel(ctx, &carsvc.DeleteCarModelRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarModelHandler) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	logger := pkglog.WithMethod(h.log, "GetImageUploadData")

	res, err := h.client.GetCarModelImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}
