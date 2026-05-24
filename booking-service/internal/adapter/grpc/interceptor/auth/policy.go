package auth

import (
	bookingpb "carsharing/protos/gen/service/booking"
	sharedmodel "carsharing/shared/model"
)

type methodPolicy struct {
	public       bool
	allowedRoles []sharedmodel.Role
}

var bookingManagerRoles = []sharedmodel.Role{sharedmodel.RoleAdmin, sharedmodel.RoleBookingManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		bookingpb.HealthService_Health_FullMethodName: {public: true},

		bookingpb.BookingService_CreateBooking_FullMethodName:           {},
		bookingpb.BookingService_GetBooking_FullMethodName:              {},
		bookingpb.BookingService_ListBookings_FullMethodName:            {},
		bookingpb.BookingService_CancelBooking_FullMethodName:           {},
		bookingpb.BookingService_UpdateBookingStatus_FullMethodName:     {allowedRoles: bookingManagerRoles},
		bookingpb.BookingService_GetBookingStatusHistory_FullMethodName: {allowedRoles: bookingManagerRoles},

		bookingpb.PricingRuleService_CreatePricingRule_FullMethodName: {allowedRoles: bookingManagerRoles},
		bookingpb.PricingRuleService_GetPricingRule_FullMethodName:    {},
		bookingpb.PricingRuleService_ListPricingRules_FullMethodName:  {},
		bookingpb.PricingRuleService_UpdatePricingRule_FullMethodName: {allowedRoles: bookingManagerRoles},
		bookingpb.PricingRuleService_DeletePricingRule_FullMethodName: {allowedRoles: bookingManagerRoles},
	}
}
