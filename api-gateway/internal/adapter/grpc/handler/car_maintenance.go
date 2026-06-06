package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	carsvc "carsharing/protos/gen/service/car"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarMaintenanceHandler struct {
	client       carsvc.CarMaintenanceServiceClient
	streamClient carsvc.CarMaintenanceStreamServiceClient
	log          *slog.Logger
}

func NewCarMaintenanceHandler(client carsvc.CarMaintenanceServiceClient, streamClient carsvc.CarMaintenanceStreamServiceClient, logger *slog.Logger) *CarMaintenanceHandler {
	return &CarMaintenanceHandler{
		client:       client,
		streamClient: streamClient,
		log:          pkglog.WithComponent(logger, "grpc.CarMaintenanceHandler"),
	}
}

func (h *CarMaintenanceHandler) CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CreateTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	res, err := h.client.CreateMaintenanceTemplate(ctx, &carsvc.CreateMaintenanceTemplateRequest{
		Name:        data.Name,
		KmInterval:  data.KmInterval,
		DayInterval: data.DayInterval,
		IsMandatory: data.IsMandatory,
		WarnPct:     data.WarnPct,
		PullPct:     data.PullPct,
	})
	if err != nil {
		log.Warn("creating maintenance template", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	log.Debug("maintenance template created", slog.String("id", res.GetId()))

	return res.GetId(), nil
}

func (h *CarMaintenanceHandler) GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	res, err := h.client.GetMaintenanceTemplate(ctx, &carsvc.GetMaintenanceTemplateRequest{Id: id})
	if err != nil {
		log.Warn("getting maintenance template", pkglog.Err(err))

		return model.CarMaintenanceTemplate{}, dto.FromGrpcErr(err)
	}

	return dto.CarMaintenanceTemplateFromProto(res.GetTemplate()), nil
}

func (h *CarMaintenanceHandler) ListTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListTemplates"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	req := &carsvc.ListMaintenanceTemplatesRequest{}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListMaintenanceTemplates(ctx, req)
	if err != nil {
		log.Warn("listing maintenance templates", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	templates := make([]model.CarMaintenanceTemplate, len(res.GetTemplates()))
	for i, t := range res.GetTemplates() {
		templates[i] = dto.CarMaintenanceTemplateFromProto(t)
	}

	return templates, nil
}

func (h *CarMaintenanceHandler) UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

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
		log.Warn("updating maintenance template", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) DeleteTemplate(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "DeleteTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	_, err := h.client.DeleteMaintenanceTemplate(ctx, &carsvc.DeleteMaintenanceTemplateRequest{Id: id})
	if err != nil {
		log.Warn("deleting maintenance template", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) ListRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListRecords"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

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
		log.Warn("listing maintenance records", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	records := make([]model.CarMaintenanceRecord, len(res.GetRecords()))
	for i, r := range res.GetRecords() {
		records[i] = dto.CarMaintenanceRecordFromProto(r)
	}

	return records, nil
}

func (h *CarMaintenanceHandler) CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CompleteRecord"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	_, err := h.client.CompleteMaintenanceRecord(ctx, &carsvc.CompleteMaintenanceRecordRequest{
		RecordId:              recordID,
		MileageAtCompletionKm: data.MileageAtCompletionKM,
		CostTenge:             data.CostTenge,
		ReceiptImageKeys:      data.ReceiptImageKeys,
		Notes:                 data.Notes,
	})
	if err != nil {
		log.Warn("completing maintenance record", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) AssignTemplate(ctx context.Context, data model.CarMaintenanceTemplateAssign) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "AssignTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	req := &carsvc.AssignCarTemplateRequest{
		CarId:      data.CarID,
		TemplateId: data.TemplateID,
	}
	if data.InitialKM != nil {
		km := int64(*data.InitialKM)
		req.InitialKm = &km
	}
	if data.InitialDate != nil {
		req.InitialDate = timestamppb.New(*data.InitialDate)
	}

	_, err := h.client.AssignCarTemplate(ctx, req)
	if err != nil {
		log.Warn("assigning car template", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarMaintenanceHandler) StreamMaintenanceEvents(ctx context.Context, send func(model.CarMaintenanceEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamMaintenanceEvents"), utils.MetadataFromCtx(ctx))

	for {
		if ctx.Err() != nil {
			return nil
		}

		stream, err := h.streamClient.StreamMaintenanceEvents(ctx, &emptypb.Empty{})
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if isUnavailable(err) {
				log.Warn("transient error opening maintenance events stream, reconnecting", pkglog.Err(err))
				select {
				case <-time.After(streamReconnectDelay):
				case <-ctx.Done():
					return nil
				}
				continue
			}
			log.Warn("streaming maintenance events", pkglog.Err(err))
			return dto.FromGrpcErr(err)
		}

		for {
			msg, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				if isUnavailable(err) {
					log.Warn("maintenance events stream interrupted, reconnecting", pkglog.Err(err))
					select {
					case <-time.After(streamReconnectDelay):
					case <-ctx.Done():
						return nil
					}
					break
				}
				log.Warn("receiving maintenance events stream", pkglog.Err(err))
				return dto.FromGrpcErr(err)
			}

			if err = send(dto.CarMaintenanceEventFromProto(msg)); err != nil {
				return err
			}
		}
	}
}

func (h *CarMaintenanceHandler) GetReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetReceiptImageUploadData"), utils.MetadataFromCtx(ctx))
	log.Debug("calling car service")

	res, err := h.client.GetMaintenanceReceiptImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		log.Warn("getting maintenance receipt image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}
