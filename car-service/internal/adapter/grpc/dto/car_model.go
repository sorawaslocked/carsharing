package dto

import (
	"car-rental-car-service/internal/model"

	"github.com/sorawaslocked/car-rental-protos/gen/base"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateCarModelRequest(req *carsvc.CreateCarModelRequest) model.CarModelCreateInput {
	return model.CarModelCreateInput{
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         int16(req.Year),
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		Seats:        int8(req.Seats),
		EngineVolume: req.EngineVolume,
		RangeKM:      req.RangeKM,
		Features:     req.Features,
	}
}

func FromGetCarModelRequest(req *carsvc.GetCarModelRequest) model.CarModelFilterInput {
	return model.CarModelFilterInput{
		ID: &req.ID,
	}
}

func FromGetCarModelsRequest(req *carsvc.GetCarModelsRequest) model.CarModelFilterInput {
	filter := model.CarModelFilterInput{
		Brand:        req.Brand,
		Model:        req.Model,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		PaginationInput: model.PaginationInput{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
	if req.MinSeats != nil {
		filter.MinSeats = new(int8(*req.MinSeats))
	}

	return filter
}

func FromUpdateCarModelRequest(req *carsvc.UpdateCarModelRequest) (model.CarModelFilterInput, model.CarModelUpdateInput) {
	filter := model.CarModelFilterInput{
		ID: &req.ID,
	}

	update := model.CarModelUpdateInput{
		Brand:        req.Brand,
		Model:        req.Model,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		EngineVolume: req.EngineVolume,
		RangeKM:      req.RangeKM,
		Features:     req.Features,
	}
	if req.Year != nil {
		update.Year = new(int16(*req.Year))
	}
	if req.Seats != nil {
		update.Seats = new(int8(*req.Seats))
	}

	return filter, update
}

func FromDeleteCarModelRequest(req *carsvc.DeleteCarModelRequest) model.CarModelFilterInput {
	return model.CarModelFilterInput{
		ID: &req.ID,
	}
}

func ToCarModelProto(cm model.CarModel) *base.CarModel {
	return &base.CarModel{
		ID:           cm.ID,
		Brand:        cm.Brand,
		Model:        cm.Model,
		Year:         int32(cm.Year),
		FuelType:     string(cm.FuelType),
		Transmission: string(cm.Transmission),
		BodyType:     string(cm.BodyType),
		Class:        string(cm.Class),
		Seats:        int32(cm.Seats),
		EngineVolume: cm.EngineVolume,
		RangeKM:      cm.RangeKM,
		Features:     cm.Features,
		CreatedAt:    timestamppb.New(cm.CreatedAt),
		UpdatedAt:    timestamppb.New(cm.UpdatedAt),
	}
}

func ToCarModelProtos(cms []model.CarModel) []*base.CarModel {
	cmProtos := make([]*base.CarModel, len(cms))

	for i, cm := range cms {
		cmProtos[i] = ToCarModelProto(cm)
	}

	return cmProtos
}
