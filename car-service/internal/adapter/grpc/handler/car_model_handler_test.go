package handler

import (
	"context"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"carsharing/car-service/internal/model"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCarModelHandlerCreateCarModel(t *testing.T) {
	ctx := context.Background()

	t.Run("returns id from service", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Create(ctx, mock.Anything).Return("model-123", nil)

		resp, err := h.CreateCarModel(ctx, &carsvc.CreateCarModelRequest{
			Brand: "Toyota", Model: "Camry", Year: 2020,
			FuelType: "petrol", Transmission: "auto", BodyType: "sedan", Class: "comfort",
			Seats: 5,
		})
		assert.NoError(t, err)
		assert.Equal(t, "model-123", resp.Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Create(ctx, mock.Anything).Return("", model.ErrInternalServerError)

		_, err := h.CreateCarModel(ctx, &carsvc.CreateCarModelRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarModelHandlerGetCarModel(t *testing.T) {
	ctx := context.Background()
	modelID := "model-123"

	t.Run("returns populated car model proto", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Get(ctx, modelID).Return(model.CarModel{
			ID: modelID, Brand: "BMW", Model: "X5",
		}, nil)

		resp, err := h.GetCarModel(ctx, &carsvc.GetCarModelRequest{Id: modelID})
		assert.NoError(t, err)
		assert.Equal(t, modelID, resp.CarModel.Id)
		assert.Equal(t, "BMW", resp.CarModel.Brand)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Get(ctx, modelID).Return(model.CarModel{}, model.ErrNotFound)

		_, err := h.GetCarModel(ctx, &carsvc.GetCarModelRequest{Id: modelID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarModelHandlerListCarModels(t *testing.T) {
	ctx := context.Background()

	t.Run("returns car model list", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().GetAll(ctx, mock.Anything).Return([]model.CarModel{
			{ID: "m-1", Brand: "Toyota"},
			{ID: "m-2", Brand: "Honda"},
		}, nil)

		resp, err := h.ListCarModels(ctx, &carsvc.ListCarModelsRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.CarModels, 2)
		assert.Equal(t, "m-1", resp.CarModels[0].Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().GetAll(ctx, mock.Anything).Return(nil, model.ErrInternalServerError)

		_, err := h.ListCarModels(ctx, &carsvc.ListCarModelsRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarModelHandlerUpdateCarModel(t *testing.T) {
	ctx := context.Background()
	modelID := "model-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Update(ctx, modelID, mock.Anything).Return(nil)

		resp, err := h.UpdateCarModel(ctx, &carsvc.UpdateCarModelRequest{Id: modelID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Update(ctx, modelID, mock.Anything).Return(model.ErrNotFound)

		_, err := h.UpdateCarModel(ctx, &carsvc.UpdateCarModelRequest{Id: modelID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarModelHandlerDeleteCarModel(t *testing.T) {
	ctx := context.Background()
	modelID := "model-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Delete(ctx, modelID).Return(nil)

		resp, err := h.DeleteCarModel(ctx, &carsvc.DeleteCarModelRequest{Id: modelID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().Delete(ctx, modelID).Return(model.ErrNotFound)

		_, err := h.DeleteCarModel(ctx, &carsvc.DeleteCarModelRequest{Id: modelID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarModelHandlerGetCarModelImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns presigned URL and object key", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().GetImageUploadData(ctx).Return(model.ImageUploadData{
			URL:       "https://upload.example.com/put",
			ObjectKey: "models/abc.jpg",
		}, nil)

		resp, err := h.GetCarModelImageUploadData(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "https://upload.example.com/put", resp.UploadData.PresignedPutUrl)
		assert.Equal(t, "models/abc.jpg", resp.UploadData.ObjectKey)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarModelService(t)
		h := NewCarModelHandler(svc, discardLogger())

		svc.EXPECT().GetImageUploadData(ctx).Return(model.ImageUploadData{}, model.ErrInternalServerError)

		_, err := h.GetCarModelImageUploadData(ctx, &emptypb.Empty{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}
