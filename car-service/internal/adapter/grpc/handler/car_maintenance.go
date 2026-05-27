package handler

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"
	pkglog "carsharing/shared/pkg/log"

	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CarMaintenanceHandler struct {
	log                *slog.Logger
	maintenanceService CarMaintenanceService

	carsvc.UnimplementedCarMaintenanceServiceServer
}

func NewCarMaintenanceHandler(log *slog.Logger, maintenanceService CarMaintenanceService) *CarMaintenanceHandler {
	return &CarMaintenanceHandler{
		log:                pkglog.WithComponent(log, "adapter.grpc.handler.CarMaintenanceHandler"),
		maintenanceService: maintenanceService,
	}
}

func (h *CarMaintenanceHandler) CreateMaintenanceTemplate(ctx context.Context, req *carsvc.CreateMaintenanceTemplateRequest) (*carsvc.CreateMaintenanceTemplateResponse, error) {
	createInput := dto.FromCreateMaintenanceTemplateRequest(req)

	id, err := h.maintenanceService.CreateTemplate(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateMaintenanceTemplateResponse{Id: id}, nil
}

func (h *CarMaintenanceHandler) GetMaintenanceTemplate(ctx context.Context, req *carsvc.GetMaintenanceTemplateRequest) (*carsvc.GetMaintenanceTemplateResponse, error) {
	template, err := h.maintenanceService.GetTemplate(ctx, req.Id)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetMaintenanceTemplateResponse{Template: dto.ToCarMaintenanceTemplateProto(template)}, nil
}

func (h *CarMaintenanceHandler) ListMaintenanceTemplates(ctx context.Context, req *carsvc.ListMaintenanceTemplatesRequest) (*carsvc.ListMaintenanceTemplatesResponse, error) {
	filterInput := dto.FromListMaintenanceTemplatesRequest(req)

	templates, err := h.maintenanceService.ListTemplates(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.ListMaintenanceTemplatesResponse{Templates: dto.ToCarMaintenanceTemplateProtos(templates)}, nil
}

func (h *CarMaintenanceHandler) UpdateMaintenanceTemplate(ctx context.Context, req *carsvc.UpdateMaintenanceTemplateRequest) (*emptypb.Empty, error) {
	updateInput := dto.FromUpdateMaintenanceTemplateRequest(req)

	if err := h.maintenanceService.UpdateTemplate(ctx, req.Id, updateInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarMaintenanceHandler) DeleteMaintenanceTemplate(ctx context.Context, req *carsvc.DeleteMaintenanceTemplateRequest) (*emptypb.Empty, error) {
	if err := h.maintenanceService.DeleteTemplate(ctx, req.Id); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarMaintenanceHandler) AssignCarTemplate(ctx context.Context, req *carsvc.AssignCarTemplateRequest) (*emptypb.Empty, error) {
	if err := h.maintenanceService.AssignCarTemplate(ctx, dto.FromAssignCarTemplateRequest(req)); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarMaintenanceHandler) ListMaintenanceRecords(ctx context.Context, req *carsvc.ListMaintenanceRecordsRequest) (*carsvc.ListMaintenanceRecordsResponse, error) {
	filterInput := dto.FromListMaintenanceRecordsRequest(req)

	records, err := h.maintenanceService.ListRecords(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.ListMaintenanceRecordsResponse{Records: dto.ToCarMaintenanceRecordProtos(records)}, nil
}

func (h *CarMaintenanceHandler) CompleteMaintenanceRecord(ctx context.Context, req *carsvc.CompleteMaintenanceRecordRequest) (*emptypb.Empty, error) {
	completeInput := dto.FromCompleteMaintenanceRecordRequest(req)

	if err := h.maintenanceService.CompleteRecord(ctx, req.RecordId, completeInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarMaintenanceHandler) GetMaintenanceReceiptImageUploadData(ctx context.Context, _ *emptypb.Empty) (*carsvc.GetMaintenanceReceiptImageUploadDataResponse, error) {
	data, err := h.maintenanceService.GetReceiptImageUploadData(ctx)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetMaintenanceReceiptImageUploadDataResponse{
		UploadData: dto.ToImageUploadData(data),
	}, nil
}
