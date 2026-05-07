package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CarMaintenanceHandler struct {
	client carsvc.CarMaintenanceServiceClient
	log    *slog.Logger
}

func NewCarMaintenanceHandler(client carsvc.CarMaintenanceServiceClient, logger *slog.Logger) *CarMaintenanceHandler {
	return &CarMaintenanceHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.CarMaintenanceHandler"),
	}
}

func (h *CarMaintenanceHandler) CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "CreateTemplate")

	res, err := h.client.CreateMaintenanceTemplate(ctx, &carsvc.CreateMaintenanceTemplateRequest{
		Name:        data.Name,
		KmInterval:  data.KmInterval,
		DayInterval: data.DayInterval,
		IsMandatory: data.IsMandatory,
		WarnPct:     data.WarnPct,
		PullPct:     data.PullPct,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *CarMaintenanceHandler) GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	logger := pkglog.WithMethod(h.log, "GetTemplate")

	res, err := h.client.GetMaintenanceTemplate(ctx, &carsvc.GetMaintenanceTemplateRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.CarMaintenanceTemplate{}, dto.FromGrpcErr(err)
	}

	return dto.CarMaintenanceTemplateFromProto(res.GetTemplate()), nil
}

func (h *CarMaintenanceHandler) ListTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	logger := pkglog.WithMethod(h.log, "ListTemplates")

	req := &carsvc.ListMaintenanceTemplatesRequest{}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListMaintenanceTemplates(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	templates := make([]model.CarMaintenanceTemplate, len(res.GetTemplates()))
	for i, t := range res.GetTemplates() {
		templates[i] = dto.CarMaintenanceTemplateFromProto(t)
	}

	return templates, nil
}

func (h *CarMaintenanceHandler) UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error {
	logger := pkglog.WithMethod(h.log, "UpdateTemplate")

	_, err := h.client.UpdateMaintenanceTemplate(ctx, &carsvc.UpdateMaintenanceTemplateRequest{
		Id:          id,
		Name:        data.Name,
		KmInterval:  data.KmInterval,
		DayInterval: data.DayInterval,
		IsMandatory: data.IsMandatory,
		WarnPct:     data.WarnPct,
		PullPct:     data.PullPct,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) DeleteTemplate(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "DeleteTemplate")

	_, err := h.client.DeleteMaintenanceTemplate(ctx, &carsvc.DeleteMaintenanceTemplateRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) ListRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	logger := pkglog.WithMethod(h.log, "ListRecords")

	req := &carsvc.ListMaintenanceRecordsRequest{
		CarId:      filter.CarID,
		TemplateId: filter.TemplateID,
		Status:     filter.Status,
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListMaintenanceRecords(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	records := make([]model.CarMaintenanceRecord, len(res.GetRecords()))
	for i, r := range res.GetRecords() {
		records[i] = dto.CarMaintenanceRecordFromProto(r)
	}

	return records, nil
}

func (h *CarMaintenanceHandler) CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error {
	logger := pkglog.WithMethod(h.log, "CompleteRecord")

	_, err := h.client.CompleteMaintenanceRecord(ctx, &carsvc.CompleteMaintenanceRecordRequest{
		RecordId:               recordID,
		OdometerAtCompletionKm: data.OdometerAtCompletionKM,
		CostTenge:              data.CostTenge,
		ReceiptImageKeys:       data.ReceiptImageKeys,
		Notes:                  data.Notes,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	logger := pkglog.WithMethod(h.log, "GetReceiptImageUploadData")

	res, err := h.client.GetMaintenanceReceiptImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}
