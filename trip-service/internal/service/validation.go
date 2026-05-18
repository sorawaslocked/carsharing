package service

import "carsharing/trip-service/internal/model"

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

func hasRole(roles []model.Role, role model.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func isPrivileged(roles []model.Role) bool {
	return hasRole(roles, model.RoleAdmin) || hasRole(roles, model.RoleBookingManager)
}
