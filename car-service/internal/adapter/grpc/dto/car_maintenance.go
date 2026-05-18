package dto

import (
	"carsharing/car-service/internal/model"

	basecar "github.com/sorawaslocked/car-rental-protos/gen/base/car"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateMaintenanceTemplateRequest(req *carsvc.CreateMaintenanceTemplateRequest) model.CarMaintenanceTemplateCreateInput {
	return model.CarMaintenanceTemplateCreateInput{
		Name:        req.Name,
		KmInterval:  req.KmInterval,
		DayInterval: req.DayInterval,
		IsMandatory: req.IsMandatory,
		WarnPct:     req.WarnPct,
		PullPct:     req.PullPct,
	}
}

func FromUpdateMaintenanceTemplateRequest(req *carsvc.UpdateMaintenanceTemplateRequest) model.CarMaintenanceTemplateUpdateInput {
	return model.CarMaintenanceTemplateUpdateInput{
		Name:        req.Name,
		KmInterval:  req.KmInterval,
		DayInterval: req.DayInterval,
		IsMandatory: req.IsMandatory,
		WarnPct:     req.WarnPct,
		PullPct:     req.PullPct,
	}
}

func FromListMaintenanceTemplatesRequest(req *carsvc.ListMaintenanceTemplatesRequest) model.CarMaintenanceTemplateFilterInput {
	filter := model.CarMaintenanceTemplateFilterInput{}
	if req.Pagination != nil {
		limit := req.Pagination.Limit
		offset := req.Pagination.Offset
		filter.PaginationInput = model.PaginationInput{
			Limit:  &limit,
			Offset: &offset,
		}
	}
	return filter
}

func FromListMaintenanceRecordsRequest(req *carsvc.ListMaintenanceRecordsRequest) model.CarMaintenanceRecordFilterInput {
	filter := model.CarMaintenanceRecordFilterInput{
		CarID:      req.CarId,
		TemplateID: req.TemplateId,
		Status:     req.Status,
	}
	if req.Pagination != nil {
		limit := req.Pagination.Limit
		offset := req.Pagination.Offset
		filter.PaginationInput = model.PaginationInput{
			Limit:  &limit,
			Offset: &offset,
		}
	}
	return filter
}

func FromCompleteMaintenanceRecordRequest(req *carsvc.CompleteMaintenanceRecordRequest) model.CarMaintenanceRecordCompleteInput {
	return model.CarMaintenanceRecordCompleteInput{
		CompletedKM:      req.OdometerAtCompletionKm,
		CostTenge:        req.CostTenge,
		Notes:            req.Notes,
		ReceiptImageKeys: req.ReceiptImageKeys,
	}
}

func ToCarMaintenanceTemplateProto(t model.CarMaintenanceTemplate) *basecar.CarMaintenanceTemplate {
	return &basecar.CarMaintenanceTemplate{
		Id:          t.ID,
		Name:        t.Name,
		KmInterval:  t.KmInterval,
		DayInterval: t.DayInterval,
		IsMandatory: t.IsMandatory,
		WarnPct:     t.WarnPct,
		PullPct:     t.PullPct,
	}
}

func ToCarMaintenanceTemplateProtos(templates []model.CarMaintenanceTemplate) []*basecar.CarMaintenanceTemplate {
	protos := make([]*basecar.CarMaintenanceTemplate, len(templates))
	for i, t := range templates {
		protos[i] = ToCarMaintenanceTemplateProto(t)
	}
	return protos
}

func ToCarMaintenanceRecordProto(r model.CarMaintenanceRecord) *basecar.CarMaintenanceRecord {
	proto := &basecar.CarMaintenanceRecord{
		Id:                     r.ID,
		CarId:                  r.CarID,
		TemplateId:             r.TemplateID,
		Status:                 string(r.Status),
		OdometerAtWarningKm:    r.OdometerAt,
		OdometerAtCompletionKm: r.CompletedKM,
		CostTenge:              r.CostTenge,
		AssignedTo:             r.AssignedTo,
		Notes:                  r.Notes,
		ReceiptImageUrls:       imageURLsFromImages(r.ReceiptImages),
		CreatedAt:              timestamppb.New(r.CreatedAt),
	}
	if r.DueBy != nil {
		proto.DueBy = timestamppb.New(*r.DueBy)
	}
	if r.CompletedAt != nil {
		proto.CompletedAt = timestamppb.New(*r.CompletedAt)
	}
	return proto
}

func ToCarMaintenanceRecordProtos(records []model.CarMaintenanceRecord) []*basecar.CarMaintenanceRecord {
	protos := make([]*basecar.CarMaintenanceRecord, len(records))
	for i, r := range records {
		protos[i] = ToCarMaintenanceRecordProto(r)
	}
	return protos
}
