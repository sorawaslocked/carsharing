package dto

import (
	"carsharing/api-gateway/internal/model"
	eventcarpb "carsharing/protos/gen/event/car"
)

func CarStatusUpdatedFromProto(e *eventcarpb.CarStatusUpdatedEvent) model.CarStatusUpdatedEvent {
	return model.CarStatusUpdatedEvent{
		CarID:      e.GetCarId(),
		FromStatus: e.GetFromStatus(),
		ToStatus:   e.GetToStatus(),
	}
}
