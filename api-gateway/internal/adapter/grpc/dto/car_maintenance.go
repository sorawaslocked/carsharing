package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "carsharing/protos/gen/base/car"
	carsvc "carsharing/protos/gen/service/car"
)

func CarMaintenanceTemplateFromProto(t *basecarpb.CarMaintenanceTemplate) model.CarMaintenanceTemplate {
	return model.CarMaintenanceTemplate{
		ID:          t.GetId(),
		Name:        t.GetName(),
		KmInterval:  t.KmInterval,
		DayInterval: t.DayInterval,
		IsMandatory: t.GetIsMandatory(),
		WarnPct:     t.GetWarnPct(),
		PullPct:     t.GetPullPct(),
	}
}

func CarMaintenanceEventFromProto(e *carsvc.MaintenanceEvent) model.CarMaintenanceEvent {
	event := model.CarMaintenanceEvent{
		CarID:      e.GetCarId(),
		TemplateID: e.GetTemplateId(),
		RecordID:   e.GetRecordId(),
		EventType:  e.GetEventType(),
	}
	if t := e.GetOccurredAt(); t != nil {
		event.OccurredAt = t.AsTime()
	}
	return event
}

func CarMaintenanceRecordFromProto(r *basecarpb.CarMaintenanceRecord) model.CarMaintenanceRecord {
	rec := model.CarMaintenanceRecord{
		ID:                    r.GetId(),
		CarID:                 r.GetCarId(),
		TemplateID:            r.GetTemplateId(),
		Status:                r.GetStatus(),
		MileageAtWarningKM:    r.GetMileageAtWarningKm(),
		MileageAtCompletionKM: r.MileageAtCompletionKm,
		CostTenge:             r.CostTenge,
		AssignedTo:            r.AssignedTo,
		ReceiptImages:         ImagesFromProto(r.GetReceiptImages()),
		Notes:                 r.Notes,
	}
	if r.GetDueBy() != nil {
		t := r.GetDueBy().AsTime()
		rec.DueBy = &t
	}
	if r.GetCompletedAt() != nil {
		t := r.GetCompletedAt().AsTime()
		rec.CompletedAt = &t
	}
	if r.GetCreatedAt() != nil {
		rec.CreatedAt = r.GetCreatedAt().AsTime()
	}
	return rec
}
