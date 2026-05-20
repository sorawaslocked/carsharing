package interceptor

import (
	sharedmodel "carsharing/shared/model"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"
)

// tripManagerRoles — roles that may view status history for any trip.
var tripManagerRoles = []sharedmodel.Role{sharedmodel.RoleAdmin, sharedmodel.RoleBookingManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// HealthService — public.
		tripsvc.HealthService_Health_FullMethodName: {public: true},

		// TripService — authentication required for all methods.
		// Ownership and access scoping are enforced by the service layer.
		tripsvc.TripService_StartTrip_FullMethodName:  {},
		tripsvc.TripService_GetTrip_FullMethodName:    {},
		tripsvc.TripService_ListTrips_FullMethodName:  {},
		tripsvc.TripService_EndTrip_FullMethodName:    {},
		tripsvc.TripService_CancelTrip_FullMethodName: {},

		tripsvc.TripService_GetTripSummary_FullMethodName: {},
		// Status history exposes internal audit data; restricted to managers.
		tripsvc.TripService_GetTripStatusHistory_FullMethodName: {allowedRoles: tripManagerRoles},

		// TripStreamService — authentication required; ownership enforced by service layer.
		tripsvc.TripStreamService_StreamTripLiveFeed_FullMethodName: {},
	}
}
