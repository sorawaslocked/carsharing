package dto

import (
	"fmt"

	"carsharing/booking-service/internal/model"
)

func WhereClausesFromBookingListFilter(filter model.BookingListFilter, args []any, argNumber int) ([]string, []any, int) {
	var clauses []string
	if args == nil {
		args = []any{}
	}

	if filter.UserID != nil {
		clauses = append(clauses, fmt.Sprintf("user_id = $%d", argNumber))
		args = append(args, *filter.UserID)
		argNumber++
	}
	if filter.CarID != nil {
		clauses = append(clauses, fmt.Sprintf("car_id = $%d", argNumber))
		args = append(args, *filter.CarID)
		argNumber++
	}
	if filter.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argNumber))
		args = append(args, *filter.Status)
		argNumber++
	}
	if filter.PricingRuleID != nil {
		clauses = append(clauses, fmt.Sprintf("pricing_rule_id = $%d", argNumber))
		args = append(args, *filter.PricingRuleID)
		argNumber++
	}

	return clauses, args, argNumber
}

func WhereClausesFromStatusHistoryFilter(filter model.BookingStatusHistoryFilter, args []any, argNumber int) ([]string, []any, int) {
	if args == nil {
		args = []any{}
	}

	clauses := []string{fmt.Sprintf("booking_id = $%d", argNumber)}
	args = append(args, filter.BookingID)
	argNumber++

	if filter.TimeRange != nil {
		clauses = append(clauses, fmt.Sprintf("changed_at >= $%d", argNumber))
		args = append(args, filter.TimeRange.From)
		argNumber++

		clauses = append(clauses, fmt.Sprintf("changed_at <= $%d", argNumber))
		args = append(args, filter.TimeRange.To)
		argNumber++
	}

	return clauses, args, argNumber
}
