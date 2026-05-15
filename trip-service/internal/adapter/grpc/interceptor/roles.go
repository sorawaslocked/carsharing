package interceptor

import (
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"
	"github.com/sorawaslocked/car-rental-trip-service/internal/model"
)

// tripManagerRoles — roles that bypass the owner check on ListTrips.
var tripManagerRoles = []model.Role{model.RoleAdmin, model.RoleBookingManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// HealthService — public.
		tripsvc.HealthService_Health_FullMethodName: {public: true},

		// TripService — authentication required for all methods.
		// ListTrips: managers see all trips; regular users must supply their own user_id.
		// Ownership enforcement for all other methods is delegated to the service layer.
		tripsvc.TripService_StartTrip_FullMethodName:            {},
		tripsvc.TripService_GetTrip_FullMethodName:              {},
		tripsvc.TripService_ListTrips_FullMethodName:            {allowedRoles: tripManagerRoles, ownerExtract: extractByUserID},
		tripsvc.TripService_EndTrip_FullMethodName:              {},
		tripsvc.TripService_CancelTrip_FullMethodName:           {},
		tripsvc.TripService_GetTripSummary_FullMethodName:       {},
		tripsvc.TripService_GetTripStatusHistory_FullMethodName: {},

		// TripStreamService — authentication required; ownership enforced by service layer.
		tripsvc.TripStreamService_StreamTripLiveFeed_FullMethodName: {},
	}
}
