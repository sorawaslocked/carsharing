package dto

import (
	eventuserpb "github.com/sorawaslocked/car-rental-protos/gen/event/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
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
