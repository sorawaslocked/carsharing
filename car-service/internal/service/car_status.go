package service

import "github.com/sorawaslocked/car-rental-car-service/internal/model"

var allowedCarStatusTransitions = map[model.CarStatus][]model.CarStatus{
	model.CarStatusAvailable:    {model.CarStatusReserved, model.CarStatusMaintenance, model.CarStatusOutOfService},
	model.CarStatusReserved:     {model.CarStatusInUse, model.CarStatusAvailable, model.CarStatusMaintenance, model.CarStatusOutOfService},
	model.CarStatusInUse:        {model.CarStatusAvailable, model.CarStatusMaintenance, model.CarStatusOutOfService},
	model.CarStatusMaintenance:  {model.CarStatusAvailable, model.CarStatusOutOfService},
	model.CarStatusOutOfService: nil,
}

func isCarStatusTransitionAllowed(from, to model.CarStatus) bool {
	allowed, exists := allowedCarStatusTransitions[from]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func transitionCarStatus(from, to model.CarStatus) error {
	if !isCarStatusTransitionAllowed(from, to) {
		return ErrInvalidStatusTransition{
			From: from,
			To:   to,
		}
	}

	return nil
}
