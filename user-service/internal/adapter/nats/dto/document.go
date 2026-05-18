package dto

import (
	"carsharing/user-service/internal/model"
	eventuserpb "github.com/sorawaslocked/car-rental-protos/gen/event/user"
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
		Passed:     e.GetPassed(),
		Defects:    defects,
	}
}
