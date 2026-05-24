package dto

import (
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	sharedvalidation "carsharing/shared/validation"

	"carsharing/protos/gen/base"
	basecar "carsharing/protos/gen/base/car"
	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateCarRequest(req *carsvc.CreateCarRequest) validation.CarCreate {
	input := validation.CarCreate{
		ModelID:          req.ModelId,
		VIN:              req.Vin,
		LicensePlate:     req.LicensePlate,
		Color:            req.Color,
		YearManufactured: int16(req.YearManufactured),
		TelemetryID:      req.TelemetryId,
		MileageKM:        req.MileageKm,
		FuelLevel:        req.FuelLevel,
		BatteryLevel:     req.BatteryLevel,
		Notes:            req.Notes,
	}
	if req.Location != nil {
		input.Location = &sharedvalidation.Location{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}
	return input
}

func FromListCarsRequest(req *carsvc.ListCarsRequest) validation.CarFilter {
	carModelFilterIsNil := true
	carModelFilter := validation.CarModelFilter{}

	if req.Brand != nil {
		carModelFilter.Brand = req.Brand
		carModelFilterIsNil = false
	}
	if req.Model != nil {
		carModelFilter.Model = req.Model
		carModelFilterIsNil = false
	}
	if req.FuelType != nil {
		carModelFilter.FuelType = req.FuelType
		carModelFilterIsNil = false
	}
	if req.Transmission != nil {
		carModelFilter.Transmission = req.Transmission
		carModelFilterIsNil = false
	}
	if req.BodyType != nil {
		carModelFilter.BodyType = req.BodyType
		carModelFilterIsNil = false
	}
	if req.Class != nil {
		carModelFilter.Class = req.Class
		carModelFilterIsNil = false
	}
	if req.MinSeats != nil {
		v := int8(*req.MinSeats)
		carModelFilter.MinSeats = &v
		carModelFilterIsNil = false
	}

	locationFilterIsNil := true
	locationFilter := sharedvalidation.LocationFilter{}

	if req.Location != nil {
		locationFilter.Location.Latitude = req.Location.Latitude
		locationFilter.Location.Longitude = req.Location.Longitude
		locationFilterIsNil = false
	}
	if req.RadiusM != nil {
		locationFilter.RadiusKM = float64(*req.RadiusM) / 1000
		locationFilterIsNil = false
	}

	carFilter := validation.CarFilter{
		Status: req.Status,
	}
	if req.Pagination != nil {
		carFilter.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}

	if !carModelFilterIsNil {
		carFilter.ModelFilter = &carModelFilter
	}
	if !locationFilterIsNil {
		carFilter.LocationFilter = &locationFilter
	}

	return carFilter
}

func FromUpdateCarRequest(req *carsvc.UpdateCarRequest) validation.CarUpdate {
	return validation.CarUpdate{
		ModelID:      req.ModelId,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		TelemetryID:  req.TelemetryId,
		ZoneID:       req.ZoneId,
		IsRetired:    req.IsRetired,
		Notes:        req.Notes,
		ImageKeys:    req.ImageKeys,
	}
}

func FromUpdateCarStatusRequest(req *carsvc.UpdateCarStatusRequest) validation.CarStatusUpdate {
	return validation.CarStatusUpdate{
		Status: req.Status,
	}
}

func FromUpdateCarTelemetryRequest(req *carsvc.UpdateCarTelemetryRequest) validation.CarTelemetryUpdate {
	input := validation.CarTelemetryUpdate{
		FuelLevel:    req.FuelLevel,
		BatteryLevel: req.BatteryLevel,
	}
	if req.MileageKm != nil {
		input.MileageKM = *req.MileageKm
	}
	if req.Location != nil {
		input.Location = &sharedvalidation.Location{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}
	return input
}

func FromGetCarStatusHistoryRequest(req *carsvc.GetCarStatusHistoryRequest) validation.CarStatusReadingFilter {
	filter := validation.CarStatusReadingFilter{
		CarID:      req.CarId,
		FromStatus: req.FromStatus,
		ToStatus:   req.ToStatus,
	}
	if req.Pagination != nil {
		filter.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func FromGetCarTelemetryHistoryRequest(req *carsvc.GetCarTelemetryHistoryRequest) validation.TelemetryReadingFilter {
	filter := validation.TelemetryReadingFilter{
		CarID: req.CarId,
	}
	if req.From != nil {
		t := req.From.AsTime()
		filter.From = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		filter.To = &t
	}
	if req.Pagination != nil {
		filter.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func ToCarProto(c model.Car) *basecar.Car {
	proto := &basecar.Car{
		Id:               c.ID,
		ModelId:          c.ModelID,
		Vin:              c.VIN,
		LicensePlate:     c.LicensePlate,
		Color:            c.Color,
		YearManufactured: int32(c.YearManufactured),
		MileageKm:        c.MileageKM,
		FuelLevel:        c.FuelLevel,
		BatteryLevel:     c.BatteryLevel,
		TelemetryId:      c.TelemetryID,
		ZoneId:           c.ZoneID,
		FuelStatus:       c.FuelStatus,
		IsRetired:        c.IsRetired,
		Status:           string(c.Status),
		Notes:            c.Notes,
		ImageUrls:        imageURLsFromImages(c.Images),
		LastSeenAt:       timestamppb.New(c.LastSeenAt),
		CreatedAt:        timestamppb.New(c.CreatedAt),
		UpdatedAt:        timestamppb.New(c.UpdatedAt),
	}
	if c.Location.Latitude != 0 || c.Location.Longitude != 0 {
		proto.Location = &base.Location{
			Latitude:  c.Location.Latitude,
			Longitude: c.Location.Longitude,
		}
	}
	return proto
}

func ToCarProtos(cs []model.Car) []*basecar.Car {
	protos := make([]*basecar.Car, len(cs))
	for i, c := range cs {
		protos[i] = ToCarProto(c)
	}
	return protos
}

func ToSlimCarProto(c model.Car) *basecar.SlimCar {
	proto := &basecar.SlimCar{
		Id:           c.ID,
		ModelId:      c.ModelID,
		LicensePlate: c.LicensePlate,
		Color:        c.Color,
		Status:       string(c.Status),
		FuelLevel:    c.FuelLevel,
	}
	if c.Location.Latitude != 0 || c.Location.Longitude != 0 {
		proto.Location = &base.Location{
			Latitude:  c.Location.Latitude,
			Longitude: c.Location.Longitude,
		}
	}
	return proto
}

func ToCarStatusReadingProtos(entries []model.CarStatusReading) []*basecar.CarStatusReading {
	protos := make([]*basecar.CarStatusReading, len(entries))
	for i, e := range entries {
		proto := &basecar.CarStatusReading{
			Id:         e.ID,
			CarId:      e.CarID,
			FromStatus: string(e.FromStatus),
			ToStatus:   string(e.ToStatus),
			ActorType:  string(e.ActorType),
			ActorId:    e.ActorID,
			Reason:     e.Reason,
			RecordedAt: timestamppb.New(e.RecordedAt),
		}
		if len(e.Metadata) > 0 {
			meta, _ := structpb.NewStruct(e.Metadata)
			proto.Metadata = meta
		}
		protos[i] = proto
	}
	return protos
}

func ToCarTelemetryReadingProtos(readings []model.TelemetryReading) []*basecar.CarTelemetryReading {
	protos := make([]*basecar.CarTelemetryReading, len(readings))
	for i, r := range readings {
		proto := &basecar.CarTelemetryReading{
			Id:           r.ID,
			CarId:        r.CarID,
			FuelPct:      r.FuelPct,
			FuelRawPct:   r.FuelRawPct,
			BatteryLevel: r.BatteryLevel,
			MileageKm:    r.MileageKM,
			ActorType:    string(r.ActorType),
			ActorId:      r.ActorID,
			Reason:       r.Reason,
			RecordedAt:   timestamppb.New(r.RecordedAt),
		}
		if r.Location != nil {
			proto.Location = &base.Location{
				Latitude:  r.Location.Latitude,
				Longitude: r.Location.Longitude,
			}
		}
		if len(r.Metadata) > 0 {
			meta, _ := structpb.NewStruct(r.Metadata)
			proto.Metadata = meta
		}
		protos[i] = proto
	}
	return protos
}

func ToImageUploadData(data sharedmodel.ImageUploadData) *base.ImageUploadData {
	return &base.ImageUploadData{
		PresignedPutUrl: data.PresignedPutURL,
		ObjectKey:       data.ObjectKey,
	}
}

func imageURLsFromImages(images []sharedmodel.Image) []string {
	if len(images) == 0 {
		return nil
	}
	urls := make([]string, 0, len(images))
	for _, img := range images {
		if img.URL != "" {
			urls = append(urls, img.URL)
		}
	}
	return urls
}
