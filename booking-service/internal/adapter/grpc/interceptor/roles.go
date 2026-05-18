package interceptor

import (
	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	bookingpb "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
)

var bookingManagerRoles = []model.Role{model.RoleAdmin, model.RoleBookingManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// HealthService — public.
		bookingpb.HealthService_Health_FullMethodName: {public: true},

		// BookingService — ownership for user-scoped methods is assessed in the service layer.
		bookingpb.BookingService_CreateBooking_FullMethodName:           {},
		bookingpb.BookingService_GetBooking_FullMethodName:              {},
		bookingpb.BookingService_ListBookings_FullMethodName:            {},
		bookingpb.BookingService_CancelBooking_FullMethodName:           {},
		bookingpb.BookingService_UpdateBookingStatus_FullMethodName:     {allowedRoles: bookingManagerRoles},
		bookingpb.BookingService_GetBookingStatusHistory_FullMethodName: {allowedRoles: bookingManagerRoles},

		// PricingRuleService — reads open to any authenticated caller; writes restricted to booking managers.
		bookingpb.PricingRuleService_CreatePricingRule_FullMethodName: {allowedRoles: bookingManagerRoles},
		bookingpb.PricingRuleService_GetPricingRule_FullMethodName:    {},
		bookingpb.PricingRuleService_ListPricingRules_FullMethodName:  {},
		bookingpb.PricingRuleService_UpdatePricingRule_FullMethodName: {allowedRoles: bookingManagerRoles},
		bookingpb.PricingRuleService_DeletePricingRule_FullMethodName: {allowedRoles: bookingManagerRoles},
	}
}
