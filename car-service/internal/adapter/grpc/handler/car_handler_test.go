package handler

import (
	"context"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	carsvc "carsharing/protos/gen/service/car"
	sharedmodel "carsharing/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCarHandlerCreateCar(t *testing.T) {
	ctx := context.Background()

	t.Run("returns id from service", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.Anything).Return("car-123", nil)

		resp, err := h.CreateCar(ctx, &carsvc.CreateCarRequest{
			ModelId: "model-1", Vin: "12345678901234567",
			LicensePlate: "ABC123", Color: "red", YearManufactured: 2022,
		})
		assert.NoError(t, err)
		assert.Equal(t, "car-123", resp.Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.Anything).Return("", errInternal)

		_, err := h.CreateCar(ctx, &carsvc.CreateCarRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerGetCar(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns populated car proto", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Get(ctx, carID).Return(model.Car{
			ID: carID, LicensePlate: "ABC-001", Status: model.CarStatusAvailable,
		}, nil)

		resp, err := h.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
		assert.NoError(t, err)
		assert.Equal(t, carID, resp.Car.Id)
		assert.Equal(t, "ABC-001", resp.Car.LicensePlate)
		assert.Equal(t, string(model.CarStatusAvailable), resp.Car.Status)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Get(ctx, carID).Return(model.Car{}, model.ErrCarNotFound)

		_, err := h.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarHandlerListCars(t *testing.T) {
	ctx := context.Background()

	t.Run("returns car list", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().List(ctx, mock.Anything).Return([]model.Car{
			{ID: "c-1"}, {ID: "c-2"},
		}, nil)

		resp, err := h.ListCars(ctx, &carsvc.ListCarsRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Cars, 2)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().List(ctx, mock.Anything).Return(nil, errInternal)

		_, err := h.ListCars(ctx, &carsvc.ListCarsRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerUpdateCar(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)

		resp, err := h.UpdateCar(ctx, &carsvc.UpdateCarRequest{Id: carID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Update(ctx, carID, mock.Anything).Return(model.ErrCarNotFound)

		_, err := h.UpdateCar(ctx, &carsvc.UpdateCarRequest{Id: carID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarHandlerUpdateCarStatus(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().UpdateCarStatus(ctx, carID, mock.MatchedBy(func(in validation.CarStatusUpdate) bool {
			return in.Status == "reserved"
		})).Return(nil)

		resp, err := h.UpdateCarStatus(ctx, &carsvc.UpdateCarStatusRequest{Id: carID, Status: "reserved"})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().UpdateCarStatus(ctx, carID, mock.Anything).Return(errInternal)

		_, err := h.UpdateCarStatus(ctx, &carsvc.UpdateCarStatusRequest{Id: carID, Status: "reserved"})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerDeleteCar(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Delete(ctx, carID).Return(nil)

		resp, err := h.DeleteCar(ctx, &carsvc.DeleteCarRequest{Id: carID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().Delete(ctx, carID).Return(model.ErrCarNotFound)

		_, err := h.DeleteCar(ctx, &carsvc.DeleteCarRequest{Id: carID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarHandlerGetCarImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().GetImageUploadData(ctx).Return(sharedmodel.ImageUploadData{
			PresignedPutURL: "https://upload.example.com/car",
			ObjectKey:       "cars/photo.jpg",
		}, nil)

		resp, err := h.GetCarImageUploadData(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "https://upload.example.com/car", resp.UploadData.PresignedPutUrl)
		assert.Equal(t, "cars/photo.jpg", resp.UploadData.ObjectKey)
	})
}

func TestCarHandlerUpdateCarTelemetry(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		mileage := int64(50_000)
		svc.EXPECT().UpdateCarTelemetry(ctx, carID, mock.MatchedBy(func(in validation.CarTelemetryUpdate) bool {
			return in.MileageKM == 50_000
		})).Return(nil)

		resp, err := h.UpdateCarTelemetry(ctx, &carsvc.UpdateCarTelemetryRequest{
			Id: carID, MileageKm: &mileage,
		})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})
}

func TestCarHandlerGetCarStatusHistory(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns status readings", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().ListCarStatusHistory(ctx, mock.Anything).Return([]model.CarStatusReading{
			{CarID: carID, FromStatus: model.CarStatusAvailable, ToStatus: model.CarStatusReserved},
		}, nil)

		resp, err := h.GetCarStatusHistory(ctx, &carsvc.GetCarStatusHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.Equal(t, string(model.CarStatusAvailable), resp.Readings[0].FromStatus)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().ListCarStatusHistory(ctx, mock.Anything).Return(nil, errInternal)

		_, err := h.GetCarStatusHistory(ctx, &carsvc.GetCarStatusHistoryRequest{CarId: carID})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerGetCarTelemetryHistory(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns telemetry readings", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().ListCarTelemetryHistory(ctx, mock.Anything).Return([]model.TelemetryReading{
			{ID: "tr-1", CarID: carID},
		}, nil)

		resp, err := h.GetCarTelemetryHistory(ctx, &carsvc.GetCarTelemetryHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.Equal(t, "tr-1", resp.Readings[0].Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(discardLogger(), svc)

		svc.EXPECT().ListCarTelemetryHistory(ctx, mock.Anything).Return(nil, errInternal)

		_, err := h.GetCarTelemetryHistory(ctx, &carsvc.GetCarTelemetryHistoryRequest{CarId: carID})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}
