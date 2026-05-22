package handler

import (
	"context"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCarHandlerCreateCar(t *testing.T) {
	ctx := context.Background()

	t.Run("returns id from service", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

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
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().Create(ctx, mock.Anything).Return("", model.ErrInternalServerError)

		_, err := h.CreateCar(ctx, &carsvc.CreateCarRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerGetCar(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns populated car proto", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

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
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().Get(ctx, carID).Return(model.Car{}, model.ErrNotFound)

		_, err := h.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarHandlerListCars(t *testing.T) {
	ctx := context.Background()

	t.Run("returns car list", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().GetAll(ctx, mock.Anything).Return([]model.Car{
			{ID: "c-1"}, {ID: "c-2"},
		}, nil)

		resp, err := h.ListCars(ctx, &carsvc.ListCarsRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Cars, 2)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().GetAll(ctx, mock.Anything).Return(nil, model.ErrInternalServerError)

		_, err := h.ListCars(ctx, &carsvc.ListCarsRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerUpdateCar(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)

		resp, err := h.UpdateCar(ctx, &carsvc.UpdateCarRequest{Id: carID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().Update(ctx, carID, mock.Anything).Return(model.ErrNotFound)

		_, err := h.UpdateCar(ctx, &carsvc.UpdateCarRequest{Id: carID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarHandlerUpdateCarStatus(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().UpdateCarStatus(ctx, carID, mock.MatchedBy(func(in validation.CarStatusUpdate) bool {
			return in.Status == "reserved"
		})).Return(nil)

		resp, err := h.UpdateCarStatus(ctx, &carsvc.UpdateCarStatusRequest{Id: carID, Status: "reserved"})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().UpdateCarStatus(ctx, carID, mock.Anything).Return(model.ErrInternalServerError)

		_, err := h.UpdateCarStatus(ctx, &carsvc.UpdateCarStatusRequest{Id: carID, Status: "reserved"})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerDeleteCar(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().Delete(ctx, carID).Return(nil)

		resp, err := h.DeleteCar(ctx, &carsvc.DeleteCarRequest{Id: carID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().Delete(ctx, carID).Return(model.ErrNotFound)

		_, err := h.DeleteCar(ctx, &carsvc.DeleteCarRequest{Id: carID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarHandlerGetCarImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

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
		h := NewCarHandler(svc, discardLogger())

		mileage := int64(50_000)
		svc.EXPECT().UpdateCarTelemetry(ctx, carID, mock.MatchedBy(func(in model.CarTelematicsUpdateInput) bool {
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
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().GetCarStatusHistory(ctx, mock.Anything).Return([]model.CarStatusLogEntry{
			{CarID: carID, FromStatus: model.CarStatusAvailable, ToStatus: model.CarStatusReserved},
		}, nil)

		resp, err := h.GetCarStatusHistory(ctx, &carsvc.GetCarStatusHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.Equal(t, string(model.CarStatusAvailable), resp.Readings[0].FromStatus)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().GetCarStatusHistory(ctx, mock.Anything).Return(nil, model.ErrInternalServerError)

		_, err := h.GetCarStatusHistory(ctx, &carsvc.GetCarStatusHistoryRequest{CarId: carID})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarHandlerGetCarFuelHistory(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns fuel readings", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		fuel := float32(75.5)
		svc.EXPECT().GetCarFuelHistory(ctx, mock.Anything).Return([]model.CarTelematicsEvent{
			{CarID: carID, FuelLevel: &fuel},
		}, nil)

		resp, err := h.GetCarFuelHistory(ctx, &carsvc.GetCarFuelHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.InDelta(t, 75.5, resp.Readings[0].FuelPct, 0.01)
	})
}

func TestCarHandlerGetCarLocationHistory(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns location readings", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().GetCarLocationHistory(ctx, mock.Anything).Return([]model.CarTelematicsEvent{
			{ID: "evt-1", CarID: carID, Latitude: 51.5, Longitude: 0.1},
		}, nil)

		resp, err := h.GetCarLocationHistory(ctx, &carsvc.GetCarLocationHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.InDelta(t, 51.5, resp.Readings[0].Location.Latitude, 0.001)
	})
}

func TestCarHandlerGetCarBatteryHistory(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns battery readings", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		battery := float32(88.0)
		svc.EXPECT().GetCarBatteryHistory(ctx, mock.Anything).Return([]model.CarTelematicsEvent{
			{ID: "evt-1", CarID: carID, BatteryLevel: &battery},
		}, nil)

		resp, err := h.GetCarBatteryHistory(ctx, &carsvc.GetCarBatteryHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.InDelta(t, 88.0, resp.Readings[0].BatteryLevel, 0.01)
	})
}

func TestCarHandlerGetCarMileageHistory(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns mileage readings", func(t *testing.T) {
		svc := mocks.NewMockCarService(t)
		h := NewCarHandler(svc, discardLogger())

		svc.EXPECT().GetCarMileageHistory(ctx, mock.Anything).Return([]model.CarTelematicsEvent{
			{ID: "evt-1", CarID: carID, OdometerKM: 123_456},
		}, nil)

		resp, err := h.GetCarMileageHistory(ctx, &carsvc.GetCarMileageHistoryRequest{CarId: carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Readings, 1)
		assert.Equal(t, int64(123_456), resp.Readings[0].MileageKm)
	})
}
