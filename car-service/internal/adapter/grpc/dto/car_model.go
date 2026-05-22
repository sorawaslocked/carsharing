package dto

import (
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"

	basecar "github.com/sorawaslocked/car-rental-protos/gen/base/car"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateCarModelRequest(req *carsvc.CreateCarModelRequest) validation.CarModelCreate {
	return validation.CarModelCreate{
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         int16(req.Year),
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		Seats:        int8(req.Seats),
		EngineVolume: req.EngineVolume,
		RangeKM:      req.RangeKm,
		Features:     req.Features,
	}
}

func FromListCarModelsRequest(req *carsvc.ListCarModelsRequest) validation.CarModelFilter {
	filter := validation.CarModelFilter{
		Brand:        req.Brand,
		Model:        req.Model,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
	}
	if req.MinSeats != nil {
		v := int8(*req.MinSeats)
		filter.MinSeats = &v
	}
	if req.Pagination != nil {
		filter.Pagination = &sharedmodel.Pagination{
			Limit:  int64(req.Pagination.Limit),
			Offset: int64(req.Pagination.Offset),
		}
	}

	return filter
}

func FromUpdateCarModelRequest(req *carsvc.UpdateCarModelRequest) validation.CarModelUpdate {
	update := validation.CarModelUpdate{
		Brand:        req.Brand,
		Model:        req.Model,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		EngineVolume: req.EngineVolume,
		Features:     req.Features,
		ImageKeys:    req.ImageKeys,
	}
	if req.Year != nil {
		v := int16(*req.Year)
		update.Year = &v
	}
	if req.Seats != nil {
		v := int8(*req.Seats)
		update.Seats = &v
	}
	if req.RangeKm != nil {
		v := *req.RangeKm
		update.RangeKM = &v
	}

	return update
}

func ToCarModelProto(cm model.CarModel) *basecar.CarModel {
	return &basecar.CarModel{
		Id:           cm.ID,
		Brand:        cm.Brand,
		Model:        cm.Model,
		Year:         int32(cm.Year),
		FuelType:     string(cm.FuelType),
		Transmission: string(cm.Transmission),
		BodyType:     string(cm.BodyType),
		Class:        string(cm.Class),
		Seats:        int32(cm.Seats),
		EngineVolume: cm.EngineVolume,
		RangeKm:      int32(cm.RangeKM),
		Features:     cm.Features,
		ImageUrls:    imageURLsFromImages(cm.Images),
		CreatedAt:    timestamppb.New(cm.CreatedAt),
		UpdatedAt:    timestamppb.New(cm.UpdatedAt),
	}
}

func ToCarModelProtos(cms []model.CarModel) []*basecar.CarModel {
	protos := make([]*basecar.CarModel, len(cms))
	for i, cm := range cms {
		protos[i] = ToCarModelProto(cm)
	}
	return protos
}
