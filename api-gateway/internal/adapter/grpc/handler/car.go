package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	carsvc "carsharing/protos/gen/service/car"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarHandler struct {
	client       carsvc.CarServiceClient
	streamClient carsvc.CarStreamServiceClient
	log          *slog.Logger
}

func NewCarHandler(client carsvc.CarServiceClient, streamClient carsvc.CarStreamServiceClient, logger *slog.Logger) *CarHandler {
	return &CarHandler{
		client:       client,
		streamClient: streamClient,
		log:          pkglog.WithComponent(logger, "grpc.CarHandler"),
	}
}

func (h *CarHandler) Create(ctx context.Context, data model.CarCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.CreateCarRequest{
		ModelId:          data.ModelID,
		Vin:              data.VIN,
		LicensePlate:     data.LicensePlate,
		Color:            data.Color,
		YearManufactured: int32(data.YearManufactured),
		TelemetryId:      data.TelemetryID,
		ZoneId:           data.ZoneID,
		MileageKm:        data.MileageKM,
		FuelLevel:        data.FuelLevel,
		BatteryLevel:     data.BatteryLevel,
		Notes:            data.Notes,
	}
	if data.Location != nil {
		req.Location = dto.LocationToProto(*data.Location)
	}
	res, err := h.client.CreateCar(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *CarHandler) Get(ctx context.Context, id string) (model.Car, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetCar(ctx, &carsvc.GetCarRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.Car{}, dto.FromGrpcErr(err)
	}

	return dto.CarFromProto(res.GetCar()), nil
}

func (h *CarHandler) List(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.ListCarsRequest{
		Brand:        filter.Brand,
		Model:        filter.Model,
		FuelType:     filter.FuelType,
		Transmission: filter.Transmission,
		BodyType:     filter.BodyType,
		Class:        filter.Class,
		ZoneId:       filter.ZoneID,
		Status:       filter.Status,
		IsRetired:    filter.IsRetired,
		MinFuelLevel: filter.MinFuelLevel,
		RadiusM:      filter.RadiusM,
	}
	if filter.MinSeats != nil {
		v := int32(*filter.MinSeats)
		req.MinSeats = &v
	}
	if filter.Location != nil {
		req.Location = dto.LocationToProto(*filter.Location)
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListCars(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	cars := make([]model.Car, len(res.GetCars()))
	for i, c := range res.GetCars() {
		cars[i] = dto.CarFromProto(c)
	}

	return cars, nil
}

func (h *CarHandler) Update(ctx context.Context, id string, data model.CarUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.UpdateCarRequest{
		Id:           id,
		ModelId:      data.ModelID,
		LicensePlate: data.LicensePlate,
		Color:        data.Color,
		TelemetryId:  data.TelemetryID,
		ZoneId:       data.ZoneID,
		IsRetired:    data.IsRetired,
		Notes:        data.Notes,
		ImageKeys:    data.ImageKeys,
	}

	_, err := h.client.UpdateCar(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.DeleteCar(ctx, &carsvc.DeleteCarRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarHandler) UpdateTelemetry(ctx context.Context, carID string, data model.CarTelemetryUpdate) error {
	logger := pkglog.WithMethod(h.log, "UpdateTelemetry")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.UpdateCarTelemetryRequest{
		Id:           carID,
		MileageKm:    data.MileageKM,
		FuelLevel:    data.FuelLevel,
		BatteryLevel: data.BatteryLevel,
		Reason:       data.Reason,
	}
	if data.Location != nil {
		req.Location = dto.LocationToProto(*data.Location)
	}
	if data.Metadata != nil {
		s, err := structpb.NewStruct(data.Metadata)
		if err == nil {
			req.Metadata = s
		}
	}

	_, err := h.client.UpdateCarTelemetry(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarHandler) UpdateStatus(ctx context.Context, carID string, data model.CarStatusUpdate) error {
	logger := pkglog.WithMethod(h.log, "UpdateStatus")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.UpdateCarStatusRequest{
		Id:     carID,
		Status: data.Status,
		Reason: data.Reason,
	}
	if data.Metadata != nil {
		s, err := structpb.NewStruct(data.Metadata)
		if err == nil {
			req.Metadata = s
		}
	}

	_, err := h.client.UpdateCarStatus(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarHandler) GetCarStatusHistory(ctx context.Context, carID string, filter model.CarStatusReadingFilter) ([]model.CarStatusReading, error) {
	logger := pkglog.WithMethod(h.log, "GetCarStatusHistory")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.GetCarStatusHistoryRequest{CarId: carID}
	if filter.TimeRange != nil {
		tr := &basepb.TimeRange{}
		if !filter.TimeRange.From.IsZero() {
			tr.From = timestamppb.New(filter.TimeRange.From)
		}
		if !filter.TimeRange.To.IsZero() {
			tr.To = timestamppb.New(filter.TimeRange.To)
		}
		req.TimeRange = tr
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.GetCarStatusHistory(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	entries := make([]model.CarStatusReading, len(res.GetReadings()))
	for i, r := range res.GetReadings() {
		entries[i] = dto.CarStatusReadingFromProto(r)
	}

	return entries, nil
}

func (h *CarHandler) GetCarTelemetryHistory(ctx context.Context, carID string, filter model.CarTelemetryReadingFilter) ([]model.CarTelemetryReading, error) {
	logger := pkglog.WithMethod(h.log, "GetCarTelemetryHistory")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.GetCarTelemetryHistoryRequest{CarId: carID}
	if filter.TimeRange != nil {
		tr := &basepb.TimeRange{}
		if !filter.TimeRange.From.IsZero() {
			tr.From = timestamppb.New(filter.TimeRange.From)
		}
		if !filter.TimeRange.To.IsZero() {
			tr.To = timestamppb.New(filter.TimeRange.To)
		}
		req.TimeRange = tr
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.GetCarTelemetryHistory(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}
		return nil, dto.FromGrpcErr(err)
	}

	readings := make([]model.CarTelemetryReading, len(res.GetReadings()))
	for i, r := range res.GetReadings() {
		readings[i] = dto.CarTelemetryReadingFromProto(r)
	}
	return readings, nil
}

func (h *CarHandler) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	logger := pkglog.WithMethod(h.log, "GetImageUploadData")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetCarImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}
		return sharedmodel.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}
