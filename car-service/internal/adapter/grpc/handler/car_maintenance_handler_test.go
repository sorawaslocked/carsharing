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

func TestCarMaintenanceHandlerCreateMaintenanceTemplate(t *testing.T) {
	ctx := context.Background()

	t.Run("returns id from service", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().CreateTemplate(ctx, mock.MatchedBy(func(in validation.CarMaintenanceTemplateCreate) bool {
			return in.Name == "Oil Change" && in.IsMandatory
		})).Return("tmpl-123", nil)

		resp, err := h.CreateMaintenanceTemplate(ctx, &carsvc.CreateMaintenanceTemplateRequest{
			Name: "Oil Change", IsMandatory: true, WarnPct: 0.8, PullPct: 1.0,
		})
		assert.NoError(t, err)
		assert.Equal(t, "tmpl-123", resp.Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().CreateTemplate(ctx, mock.Anything).Return("", model.ErrInternalServerError)

		_, err := h.CreateMaintenanceTemplate(ctx, &carsvc.CreateMaintenanceTemplateRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarMaintenanceHandlerGetMaintenanceTemplate(t *testing.T) {
	ctx := context.Background()
	tmplID := "tmpl-123"

	t.Run("returns populated template proto", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		kmInterval := int32(5_000)
		svc.EXPECT().GetTemplate(ctx, tmplID).Return(model.CarMaintenanceTemplate{
			ID: tmplID, Name: "Tire Rotation", KmInterval: &kmInterval, IsMandatory: true,
		}, nil)

		resp, err := h.GetMaintenanceTemplate(ctx, &carsvc.GetMaintenanceTemplateRequest{Id: tmplID})
		assert.NoError(t, err)
		assert.Equal(t, tmplID, resp.Template.Id)
		assert.Equal(t, "Tire Rotation", resp.Template.Name)
		assert.True(t, resp.Template.IsMandatory)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().GetTemplate(ctx, tmplID).Return(model.CarMaintenanceTemplate{}, model.ErrNotFound)

		_, err := h.GetMaintenanceTemplate(ctx, &carsvc.GetMaintenanceTemplateRequest{Id: tmplID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarMaintenanceHandlerListMaintenanceTemplates(t *testing.T) {
	ctx := context.Background()

	t.Run("returns template list", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().GetAllTemplates(ctx, mock.Anything).Return([]model.CarMaintenanceTemplate{
			{ID: "t-1", Name: "Oil Change"},
			{ID: "t-2", Name: "Brake Check"},
		}, nil)

		resp, err := h.ListMaintenanceTemplates(ctx, &carsvc.ListMaintenanceTemplatesRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Templates, 2)
		assert.Equal(t, "t-1", resp.Templates[0].Id)
	})
}

func TestCarMaintenanceHandlerUpdateMaintenanceTemplate(t *testing.T) {
	ctx := context.Background()
	tmplID := "tmpl-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().UpdateTemplate(ctx, tmplID, mock.Anything).Return(nil)

		resp, err := h.UpdateMaintenanceTemplate(ctx, &carsvc.UpdateMaintenanceTemplateRequest{Id: tmplID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().UpdateTemplate(ctx, tmplID, mock.Anything).Return(model.ErrNotFound)

		_, err := h.UpdateMaintenanceTemplate(ctx, &carsvc.UpdateMaintenanceTemplateRequest{Id: tmplID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarMaintenanceHandlerDeleteMaintenanceTemplate(t *testing.T) {
	ctx := context.Background()
	tmplID := "tmpl-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().DeleteTemplate(ctx, tmplID).Return(nil)

		resp, err := h.DeleteMaintenanceTemplate(ctx, &carsvc.DeleteMaintenanceTemplateRequest{Id: tmplID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().DeleteTemplate(ctx, tmplID).Return(model.ErrNotFound)

		_, err := h.DeleteMaintenanceTemplate(ctx, &carsvc.DeleteMaintenanceTemplateRequest{Id: tmplID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarMaintenanceHandlerListMaintenanceRecords(t *testing.T) {
	ctx := context.Background()
	carID := "car-123"

	t.Run("returns record list filtered by car", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().GetRecords(ctx, mock.MatchedBy(func(in validation.CarMaintenanceRecordFilter) bool {
			return in.CarID != nil && *in.CarID == carID
		})).Return([]model.CarMaintenanceRecord{
			{ID: "rec-1", CarID: carID, Status: model.MaintenanceRecordStatusPending},
		}, nil)

		resp, err := h.ListMaintenanceRecords(ctx, &carsvc.ListMaintenanceRecordsRequest{CarId: &carID})
		assert.NoError(t, err)
		assert.Len(t, resp.Records, 1)
		assert.Equal(t, "rec-1", resp.Records[0].Id)
		assert.Equal(t, string(model.MaintenanceRecordStatusPending), resp.Records[0].Status)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().GetRecords(ctx, mock.Anything).Return(nil, model.ErrInternalServerError)

		_, err := h.ListMaintenanceRecords(ctx, &carsvc.ListMaintenanceRecordsRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestCarMaintenanceHandlerCompleteMaintenanceRecord(t *testing.T) {
	ctx := context.Background()
	recordID := "rec-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().CompleteRecord(ctx, recordID, mock.MatchedBy(func(in validation.CarMaintenanceRecordComplete) bool {
			return in.CompletedKM == 50_000 && in.CostTenge == 15_000
		})).Return(nil)

		resp, err := h.CompleteMaintenanceRecord(ctx, &carsvc.CompleteMaintenanceRecordRequest{
			RecordId:               recordID,
			OdometerAtCompletionKm: 50_000,
			CostTenge:              15_000,
		})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().CompleteRecord(ctx, recordID, mock.Anything).Return(model.ErrNotFound)

		_, err := h.CompleteMaintenanceRecord(ctx, &carsvc.CompleteMaintenanceRecordRequest{RecordId: recordID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestCarMaintenanceHandlerGetMaintenanceReceiptImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().GetReceiptImageUploadData(ctx).Return(sharedmodel.ImageUploadData{
			PresignedPutURL: "https://upload.example.com/receipt",
			ObjectKey:       "receipts/invoice.pdf",
		}, nil)

		resp, err := h.GetMaintenanceReceiptImageUploadData(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "https://upload.example.com/receipt", resp.UploadData.PresignedPutUrl)
		assert.Equal(t, "receipts/invoice.pdf", resp.UploadData.ObjectKey)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockCarMaintenanceService(t)
		h := NewCarMaintenanceHandler(svc, discardLogger())

		svc.EXPECT().GetReceiptImageUploadData(ctx).Return(sharedmodel.ImageUploadData{}, model.ErrInternalServerError)

		_, err := h.GetMaintenanceReceiptImageUploadData(ctx, &emptypb.Empty{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}
