package dto

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/model"

	"github.com/sorawaslocked/car-rental-protos/gen/base"
	basecar "github.com/sorawaslocked/car-rental-protos/gen/base/car"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromUpdateCarTelemetryRequest(req *carsvc.UpdateCarTelemetryRequest) model.CarTelematicsUpdateInput {
	input := model.CarTelematicsUpdateInput{}
	if req.MileageKm != nil {
		input.MileageKM = *req.MileageKm
	}
	input.FuelLevel = req.FuelLevel
	input.BatteryLevel = req.BatteryLevel
	if req.Location != nil {
		input.Location = &model.Location{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		}
	}
	return input
}

func FromGetCarStatusHistoryRequest(req *carsvc.GetCarStatusHistoryRequest) model.CarStatusLogFilter {
	filter := model.CarStatusLogFilter{
		CarID: &req.CarId,
	}
	if req.Pagination != nil {
		limit := int64(req.Pagination.Limit)
		offset := int64(req.Pagination.Offset)
		filter.Pagination = model.Pagination{
			Limit:  &limit,
			Offset: &offset,
		}
	}
	return filter
}

func FromGetCarFuelHistoryRequest(req *carsvc.GetCarFuelHistoryRequest) model.TelematicsEventFilter {
	filter := model.TelematicsEventFilter{
		CarID: &req.CarId,
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
		limit := int64(req.Pagination.Limit)
		offset := int64(req.Pagination.Offset)
		filter.Pagination = model.Pagination{
			Limit:  &limit,
			Offset: &offset,
		}
	}
	return filter
}

func FromGetTelematicsHistoryRequest(carID string, from, to *timestamppb.Timestamp, pagination *base.Pagination) model.TelematicsEventFilter {
	filter := model.TelematicsEventFilter{
		CarID: &carID,
	}
	if from != nil {
		t := from.AsTime()
		filter.From = &t
	}
	if to != nil {
		t := to.AsTime()
		filter.To = &t
	}
	if pagination != nil {
		limit := pagination.Limit
		offset := pagination.Offset
		filter.Pagination = model.Pagination{
			Limit:  &limit,
			Offset: &offset,
		}
	}
	return filter
}

func ToCarStatusReadingProtos(entries []model.CarStatusLogEntry) []*basecar.CarStatusReading {
	protos := make([]*basecar.CarStatusReading, len(entries))
	for i, e := range entries {
		protos[i] = &basecar.CarStatusReading{
			Id:         e.ID,
			CarId:      e.CarID,
			FromStatus: string(e.FromStatus),
			ToStatus:   string(e.ToStatus),
			ActorType:  string(e.ActorType),
			ActorId:    e.ActorID,
			Reason:     e.Reason,
			RecordedAt: timestamppb.New(e.ChangedAt),
		}
	}
	return protos
}

func ToCarFuelReadingProtos(events []model.CarTelematicsEvent) []*basecar.CarFuelReading {
	protos := make([]*basecar.CarFuelReading, len(events))
	for i, e := range events {
		proto := &basecar.CarFuelReading{
			CarId:      e.CarID,
			RecordedAt: timestamppb.New(e.RecordedAt),
		}
		if e.FuelLevel != nil {
			proto.FuelPct = *e.FuelLevel
		}
		protos[i] = proto
	}
	return protos
}

func ToCarLocationReadingProtos(events []model.CarTelematicsEvent) []*basecar.CarLocationReading {
	protos := make([]*basecar.CarLocationReading, len(events))
	for i, e := range events {
		proto := &basecar.CarLocationReading{
			Id:         e.ID,
			CarId:      e.CarID,
			RecordedAt: timestamppb.New(e.RecordedAt),
		}
		if e.Latitude != 0 || e.Longitude != 0 {
			proto.Location = &base.Location{
				Latitude:  e.Latitude,
				Longitude: e.Longitude,
			}
		}
		protos[i] = proto
	}
	return protos
}

func ToCarBatteryReadingProtos(events []model.CarTelematicsEvent) []*basecar.CarBatteryReading {
	protos := make([]*basecar.CarBatteryReading, len(events))
	for i, e := range events {
		proto := &basecar.CarBatteryReading{
			Id:         e.ID,
			CarId:      e.CarID,
			RecordedAt: timestamppb.New(e.RecordedAt),
		}
		if e.BatteryLevel != nil {
			proto.BatteryLevel = *e.BatteryLevel
		}
		protos[i] = proto
	}
	return protos
}

func ToCarMileageReadingProtos(events []model.CarTelematicsEvent) []*basecar.CarMileageReading {
	protos := make([]*basecar.CarMileageReading, len(events))
	for i, e := range events {
		protos[i] = &basecar.CarMileageReading{
			Id:         e.ID,
			CarId:      e.CarID,
			MileageKm:  e.OdometerKM,
			RecordedAt: timestamppb.New(e.RecordedAt),
		}
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

func FromCreateCarRequest(req *carsvc.CreateCarRequest) model.CarCreateInput {
	return model.CarCreateInput{
		ModelID:          req.ModelId,
		VIN:              req.Vin,
		LicensePlate:     req.LicensePlate,
		Color:            req.Color,
		YearManufactured: int16(req.YearManufactured),
		Notes:            ptrToSlice(req.Notes),
	}
}

func FromListCarsRequest(req *carsvc.ListCarsRequest) model.CarFilterInput {
	carModelFilterIsNil := true
	carModelFilter := model.CarModelFilterInput{}

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
	locationFilter := model.LocationFilter{}

	if req.Location != nil {
		locationFilter.Location.Latitude = req.Location.Latitude
		locationFilter.Location.Longitude = req.Location.Longitude
		locationFilterIsNil = false
	}
	if req.RadiusM != nil {
		locationFilter.RadiusKM = float64(*req.RadiusM) / 1000
		locationFilterIsNil = false
	}

	carFilter := model.CarFilterInput{
		Status: req.Status,
	}
	if req.Pagination != nil {
		limit := int64(req.Pagination.Limit)
		offset := int64(req.Pagination.Offset)
		carFilter.PaginationInput = model.PaginationInput{
			Limit:  &limit,
			Offset: &offset,
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

func FromUpdateCarRequest(req *carsvc.UpdateCarRequest) model.CarUpdateInput {
	return model.CarUpdateInput{
		ModelID:      req.ModelId,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		Notes:        ptrToSlice(req.Notes),
		ImageKeys:    req.ImageKeys,
	}
}

func FromUpdateCarStatusRequest(req *carsvc.UpdateCarStatusRequest) model.CarStatusUpdateInput {
	return model.CarStatusUpdateInput{
		Status: req.Status,
	}
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
		Status:           string(c.Status),
		Notes:            sliceToPtr(c.Notes),
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

func ToImageUploadData(uploadURL, objectKey string) *base.ImageUploadData {
	return &base.ImageUploadData{
		PresignedPutUrl: uploadURL,
		ObjectKey:       objectKey,
	}
}

func imageURLsFromImages(images []model.Image) []string {
	if len(images) == 0 {
		return nil
	}
	urls := make([]string, 0, len(images))
	for _, img := range images {
		if img.URL != nil {
			urls = append(urls, *img.URL)
		}
	}
	return urls
}

func ptrToSlice(s *string) []string {
	if s == nil {
		return nil
	}
	return []string{*s}
}

func sliceToPtr(ss []string) *string {
	if len(ss) == 0 {
		return nil
	}
	return &ss[0]
}
