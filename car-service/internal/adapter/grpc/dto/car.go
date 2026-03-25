package dto

import (
	"car-rental-car-service/internal/model"

	"github.com/sorawaslocked/car-rental-protos/gen/base"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateCarRequest(req *carsvc.CreateCarRequest) model.CarCreateInput {
	return model.CarCreateInput{
		ModelID:          req.ModelID,
		VIN:              req.VIN,
		LicensePlate:     req.LicensePlate,
		Color:            req.Color,
		YearManufactured: int16(req.YearManufactured),
		MileageKM:        req.MileageKM,
		FuelLevel:        req.FuelLevel,
		BatteryLevel:     req.BatteryLevel,
		Notes:            req.Notes,
	}
}

func FromGetCarRequest(req *carsvc.GetCarRequest) model.CarFilterInput {
	return model.CarFilterInput{
		ID: &req.ID,
	}
}

func FromGetCarsRequest(req *carsvc.GetCarsRequest) model.CarFilterInput {
	carModelFilterIsNil := true
	carModelFilter := model.CarModelFilterInput{}

	if req.Brand != nil {
		carModelFilter.Brand = req.Brand
		carModelFilterIsNil = false
	}
	if req.Model != nil {
		carModelFilter.Model = req.Model
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.FuelType != nil {
		carModelFilter.FuelType = req.FuelType
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.Transmission != nil {
		carModelFilter.Transmission = req.Transmission
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.BodyType != nil {
		carModelFilter.BodyType = req.BodyType
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.Class != nil {
		carModelFilter.Class = req.Class
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.MinSeats != nil {
		carModelFilter.MinSeats = new(int8(*req.MinSeats))
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}

	locationFilterIsNil := true
	locationFilter := model.LocationFilter{
		Location: model.Location{},
	}

	if req.Latitude != nil {
		locationFilter.Location.Latitude = *req.Latitude
		locationFilterIsNil = false
	}
	if req.Longitude != nil {
		locationFilter.Location.Longitude = *req.Longitude
		if locationFilterIsNil {
			locationFilterIsNil = false
		}
	}
	if req.RadiusKM != nil {
		locationFilter.RadiusKM = *req.RadiusKM
		if locationFilterIsNil {
			locationFilterIsNil = false
		}
	}

	carFilter := model.CarFilterInput{
		Status: req.Status,
		PaginationInput: model.PaginationInput{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if !carModelFilterIsNil {
		carFilter.ModelFilter = &carModelFilter
	}
	if !locationFilterIsNil {
		carFilter.LocationFilter = &locationFilter
	}

	return carFilter
}

func FromGetAvailableCars(req *carsvc.GetAvailableCarsRequest) model.CarFilterInput {
	carModelFilterIsNil := true
	carModelFilter := model.CarModelFilterInput{}

	if req.Brand != nil {
		carModelFilter.Brand = req.Brand
		carModelFilterIsNil = false
	}
	if req.Model != nil {
		carModelFilter.Model = req.Model
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.FuelType != nil {
		carModelFilter.FuelType = req.FuelType
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.Transmission != nil {
		carModelFilter.Transmission = req.Transmission
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.BodyType != nil {
		carModelFilter.BodyType = req.BodyType
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.Class != nil {
		carModelFilter.Class = req.Class
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}
	if req.MinSeats != nil {
		carModelFilter.MinSeats = new(int8(*req.MinSeats))
		if carModelFilterIsNil {
			carModelFilterIsNil = false
		}
	}

	locationFilterIsNil := true
	locationFilter := model.LocationFilter{
		Location: model.Location{},
	}

	if req.Latitude != nil {
		locationFilter.Location.Latitude = *req.Latitude
		locationFilterIsNil = false
	}
	if req.Longitude != nil {
		locationFilter.Location.Longitude = *req.Longitude
		if locationFilterIsNil {
			locationFilterIsNil = false
		}
	}
	if req.RadiusKM != nil {
		locationFilter.RadiusKM = *req.RadiusKM
		if locationFilterIsNil {
			locationFilterIsNil = false
		}
	}

	carFilter := model.CarFilterInput{
		PaginationInput: model.PaginationInput{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if !carModelFilterIsNil {
		carFilter.ModelFilter = &carModelFilter
	}
	if !locationFilterIsNil {
		carFilter.LocationFilter = &locationFilter
	}

	return carFilter
}

func FromUpdateCarRequest(req *carsvc.UpdateCarRequest) (model.CarFilterInput, model.CarUpdateInput) {
	carFilter := model.CarFilterInput{
		ID: &req.ID,
	}

	update := model.CarUpdateInput{
		ModelID:      req.ModelID,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		Notes:        req.Notes,
	}

	return carFilter, update
}

func FromUpdateCarStatusRequest(req *carsvc.UpdateCarStatusRequest) (model.CarFilterInput, model.CarStatusUpdateInput) {
	carFilter := model.CarFilterInput{
		ID: &req.ID,
	}

	statusUpdate := model.CarStatusUpdateInput{
		Status: req.Status,
	}

	return carFilter, statusUpdate
}

func FromDeleteCarRequest(req *carsvc.DeleteCarRequest) model.CarFilterInput {
	carFilter := model.CarFilterInput{
		ID: &req.ID,
	}

	return carFilter
}

func ToCarProto(c model.Car) *base.Car {
	return &base.Car{
		ID:               c.ID,
		ModelID:          c.ModelID,
		VIN:              c.VIN,
		LicensePlate:     c.LicensePlate,
		Color:            c.Color,
		YearManufactured: int32(c.YearManufactured),
		MileageKM:        c.MileageKM,
		FuelLevel:        c.FuelLevel,
		BatteryLevel:     c.BatteryLevel,
		Latitude:         c.Location.Latitude,
		Longitude:        c.Location.Longitude,
		Status:           string(c.Status),
		Notes:            c.Notes,
		LastSeenAt:       timestamppb.New(c.LastSeenAt),
		CreatedAt:        timestamppb.New(c.CreatedAt),
		UpdatedAt:        timestamppb.New(c.UpdatedAt),
	}
}

func ToCarProtos(cs []model.Car) []*base.Car {
	cProtos := make([]*base.Car, len(cs))

	for i, c := range cs {
		cProtos[i] = ToCarProto(c)
	}

	return cProtos
}
