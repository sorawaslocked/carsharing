package handler

import (
	"context"
	"testing"
	"time"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	carsvc "carsharing/protos/gen/service/car"
	sharedmodel "carsharing/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCarInsuranceHandlerCreateCarInsurance(t *testing.T) {
	ctx := context.Background()

	req := &carsvc.CreateCarInsuranceRequest{
		CarId:     "00000000-0000-0000-0000-000000000001",
		Type:      "osago",
		Provider:  "InsureCo",
		PolicyNum: "POL-001",
		StartsAt:  timestamppb.New(time.Now()),
		ExpiresAt: timestamppb.New(time.Now().Add(365 * 24 * time.Hour)),
	}

	t.Run("returns id from service", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.MatchedBy(func(in validation.CarInsuranceCreate) bool {
			return in.CarID == req.CarId && in.Provider == "InsureCo"
		})).Return("ins-123", nil)

		resp, err := h.CreateCarInsurance(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "ins-123", resp.Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.Anything).Return("", errInternal)

		_, err := h.CreateCarInsurance(ctx, req)
		assert.Equal(t, codes.Internal, grpcCode(err))
	})

	t.Run("validation error maps to gRPC InvalidArgument", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.Anything).Return("", validation.Errors{"car_id": validation.ErrRequiredField})

		_, err := h.CreateCarInsurance(ctx, req)
		assert.Equal(t, codes.InvalidArgument, grpcCode(err))
	})
}

func TestCarInsuranceHandlerGetCarInsurance(t *testing.T) {
	ctx := context.Background()
	insID := "ins-123"

	t.Run("returns populated insurance proto", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Get(ctx, insID).Return(model.CarInsurance{
			ID:       insID,
			CarID:    "car-456",
			Type:     model.InsuranceTypeOSAGO,
			Provider: "InsureCo",
			Status:   model.InsuranceStatusActive,
		}, nil)

		resp, err := h.GetCarInsurance(ctx, &carsvc.GetCarInsuranceRequest{Id: insID})
		assert.NoError(t, err)
		assert.Equal(t, insID, resp.CarInsurance.Id)
		assert.Equal(t, "car-456", resp.CarInsurance.CarId)
		assert.Equal(t, string(model.InsuranceStatusActive), resp.CarInsurance.Status)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Get(ctx, insID).Return(model.CarInsurance{}, model.ErrCarInsuranceNotFound)

		_, err := h.GetCarInsurance(ctx, &carsvc.GetCarInsuranceRequest{Id: insID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarInsuranceHandlerListCarInsurances(t *testing.T) {
	ctx := context.Background()

	t.Run("returns insurance list", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().List(ctx, mock.Anything).Return([]model.CarInsurance{
			{ID: "i-1"}, {ID: "i-2"},
		}, nil)

		resp, err := h.ListCarInsurances(ctx, &carsvc.ListCarInsurancesRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.CarInsurances, 2)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().List(ctx, mock.Anything).Return(nil, errInternal)

		_, err := h.ListCarInsurances(ctx, &carsvc.ListCarInsurancesRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarInsuranceHandlerUpdateCarInsurance(t *testing.T) {
	ctx := context.Background()
	insID := "ins-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Update(ctx, insID, mock.Anything).Return(nil)

		resp, err := h.UpdateCarInsurance(ctx, &carsvc.UpdateCarInsuranceRequest{Id: insID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Update(ctx, insID, mock.Anything).Return(model.ErrCarInsuranceNotFound)

		_, err := h.UpdateCarInsurance(ctx, &carsvc.UpdateCarInsuranceRequest{Id: insID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarInsuranceHandlerDeleteCarInsurance(t *testing.T) {
	ctx := context.Background()
	insID := "ins-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Delete(ctx, insID).Return(nil)

		resp, err := h.DeleteCarInsurance(ctx, &carsvc.DeleteCarInsuranceRequest{Id: insID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().Delete(ctx, insID).Return(model.ErrCarInsuranceNotFound)

		_, err := h.DeleteCarInsurance(ctx, &carsvc.DeleteCarInsuranceRequest{Id: insID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarInsuranceHandlerGetCarInsuranceImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data", func(t *testing.T) {
		svc := mocks.NewMockCarInsuranceService(t)
		h := NewCarInsuranceHandler(discardLogger(), svc)

		svc.EXPECT().GetImageUploadData(ctx).Return(sharedmodel.ImageUploadData{
			PresignedPutURL: "https://upload.example.com/ins",
			ObjectKey:       "insurance/policy.pdf",
		}, nil)

		resp, err := h.GetCarInsuranceImageUploadData(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "https://upload.example.com/ins", resp.UploadData.PresignedPutUrl)
		assert.Equal(t, "insurance/policy.pdf", resp.UploadData.ObjectKey)
	})
}
