package dto

import (
	eventuserpb "carsharing/protos/gen/event/user"
	"carsharing/user-service/internal/model"
)

func DocumentAnalyzedEventFromProto(e *eventuserpb.DocumentAnalyzedEvent) model.DocumentAnalyzedEvent {
	defects := make([]model.Defect, len(e.GetDefects()))
	for i, d := range e.GetDefects() {
		defects[i] = model.Defect{
			Type:        d.GetType(),
			Description: d.GetDescription(),
		}
	}

	return model.DocumentAnalyzedEvent{
		DocumentID: e.GetDocumentId(),
		UserID:     e.GetUserId(),
		Passed:     e.GetPassed(),
		Defects:    defects,
	}
}
