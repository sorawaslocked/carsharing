package dto

import (
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedvalidation "carsharing/shared/validation"

	basecar "carsharing/protos/gen/base/car"
	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateMaintenanceTemplateRequest(req *carsvc.CreateMaintenanceTemplateRequest) validation.CarMaintenanceTemplateCreate {
	return validation.CarMaintenanceTemplateCreate{
		Name:        req.Name,
		KmInterval:  req.KmInterval,
		DayInterval: req.DayInterval,
		IsMandatory: req.IsMandatory,
		WarnPct:     req.WarnPct,
		PullPct:     req.PullPct,
	}
}

func FromUpdateMaintenanceTemplateRequest(req *carsvc.UpdateMaintenanceTemplateRequest) validation.CarMaintenanceTemplateUpdate {
	return validation.CarMaintenanceTemplateUpdate{
		Name:        req.Name,
		KmInterval:  req.KmInterval,
		DayInterval: req.DayInterval,
		IsMandatory: req.IsMandatory,
		WarnPct:     req.WarnPct,
		PullPct:     req.PullPct,
	}
}

func FromListMaintenanceTemplatesRequest(req *carsvc.ListMaintenanceTemplatesRequest) validation.CarMaintenanceTemplateFilter {
	filter := validation.CarMaintenanceTemplateFilter{}
	if req.Pagination != nil {
		filter.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func FromListMaintenanceRecordsRequest(req *carsvc.ListMaintenanceRecordsRequest) validation.CarMaintenanceRecordFilter {
	filter := validation.CarMaintenanceRecordFilter{
		CarID:      req.CarId,
		TemplateID: req.TemplateId,
		Status:     req.Status,
	}
	if req.Pagination != nil {
		filter.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func FromCompleteMaintenanceRecordRequest(req *carsvc.CompleteMaintenanceRecordRequest) validation.CarMaintenanceRecordComplete {
	return validation.CarMaintenanceRecordComplete{
		CompletedKM:      req.MileageAtCompletionKm,
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
		Id:                    r.ID,
		CarId:                 r.CarID,
		TemplateId:            r.TemplateID,
		Status:                string(r.Status),
		MileageAtWarningKm:    r.OdometerAt,
		MileageAtCompletionKm: r.CompletedKM,
		CostTenge:             r.CostTenge,
		AssignedTo:            r.AssignedTo,
		Notes:                 r.Notes,
		ReceiptImageUrls:      imageURLsFromImages(r.ReceiptImages),
		CreatedAt:             timestamppb.New(r.CreatedAt),
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
