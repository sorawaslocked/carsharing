package dto

import (
	"carsharing/api-gateway/internal/model"
	eventuserpb "carsharing/protos/gen/event/user"
)

func DocumentAnalyzedEventFromProto(e *eventuserpb.DocumentAnalyzedEvent) model.DocumentAnalyzedEvent {
	defects := make([]model.DocumentDefect, len(e.GetDefects()))
	for i, d := range e.GetDefects() {
		defects[i] = model.DocumentDefect{
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
