package auth

import (
	tripsvc "carsharing/protos/gen/service/trip"
	sharedmodel "carsharing/shared/model"
)

var tripManagerRoles = []sharedmodel.Role{sharedmodel.RoleAdmin, sharedmodel.RoleBookingManager}

type methodPolicy struct {
	public       bool
	allowedRoles []sharedmodel.Role
}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// Public — no authentication required.
		tripsvc.HealthService_Health_FullMethodName: {public: true},

		// Any authenticated user — ownership enforced by the service layer.
		tripsvc.TripService_StartTrip_FullMethodName:                {},
		tripsvc.TripService_GetTrip_FullMethodName:                  {},
		tripsvc.TripService_ListTrips_FullMethodName:                {},
		tripsvc.TripService_EndTrip_FullMethodName:                  {},
		tripsvc.TripService_CancelTrip_FullMethodName:               {},
		tripsvc.TripService_GetTripSummary_FullMethodName:           {},
		tripsvc.TripStreamService_StreamTripLiveFeed_FullMethodName: {},

		// Manager roles only — exposes internal audit data.
		tripsvc.TripService_GetTripStatusHistory_FullMethodName: {allowedRoles: tripManagerRoles},
	}
}
