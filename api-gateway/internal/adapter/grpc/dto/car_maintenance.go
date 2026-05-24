package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "carsharing/protos/gen/base/car"
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
		ReceiptImageURLs:      r.GetReceiptImageUrls(),
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
