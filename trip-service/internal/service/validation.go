package service

import (
	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/model"
)

func validateID(id string) error {
	if id == "" {
		return model.ValidationErrors{"id": model.ErrRequiredField}
	}
	return nil
}

func validateBookingID(bookingID string) error {
	if bookingID == "" {
		return model.ValidationErrors{"bookingID": model.ErrRequiredField}
	}
	return nil
}

func hasRole(roles []sharedmodel.Role, role sharedmodel.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func isPrivileged(roles []sharedmodel.Role) bool {
	return hasRole(roles, sharedmodel.RoleAdmin) || hasRole(roles, sharedmodel.RoleBookingManager)
}
