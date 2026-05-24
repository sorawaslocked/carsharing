package dto

import (
	"fmt"

	"carsharing/booking-service/internal/model"
)

func WhereClausesFromPricingRuleListFilter(filter model.PricingRuleListFilter, args []any, argNumber int) ([]string, []any, int) {
	var clauses []string
	if args == nil {
		args = []any{}
	}

	if filter.ModelID != nil {
		clauses = append(clauses, fmt.Sprintf("model_id = $%d", argNumber))
		args = append(args, *filter.ModelID)
		argNumber++
	}
	if filter.ZoneID != nil {
		clauses = append(clauses, fmt.Sprintf("zone_id = $%d", argNumber))
		args = append(args, *filter.ZoneID)
		argNumber++
	}
	if filter.Class != nil {
		clauses = append(clauses, fmt.Sprintf("class = $%d", argNumber))
		args = append(args, *filter.Class)
		argNumber++
	}
	if filter.Type != nil {
		clauses = append(clauses, fmt.Sprintf("type = $%d", argNumber))
		args = append(args, *filter.Type)
		argNumber++
	}
	if filter.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", argNumber))
		args = append(args, *filter.IsActive)
		argNumber++
	}

	return clauses, args, argNumber
}

func SetClausesFromPricingRuleUpdate(update model.PricingRuleUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	argNumber := 1

	if update.ModelID != nil {
		clauses = append(clauses, fmt.Sprintf("model_id = $%d", argNumber))
		args = append(args, *update.ModelID)
		argNumber++
	}
	if update.ZoneID != nil {
		clauses = append(clauses, fmt.Sprintf("zone_id = $%d", argNumber))
		args = append(args, *update.ZoneID)
		argNumber++
	}
	if update.Class != nil {
		clauses = append(clauses, fmt.Sprintf("class = $%d", argNumber))
		args = append(args, *update.Class)
		argNumber++
	}
	if update.Type != nil {
		clauses = append(clauses, fmt.Sprintf("type = $%d", argNumber))
		args = append(args, *update.Type)
		argNumber++
	}
	if update.RateTenge != nil {
		clauses = append(clauses, fmt.Sprintf("rate_tenge = $%d", argNumber))
		args = append(args, *update.RateTenge)
		argNumber++
	}
	if update.RatePerKMTenge != nil {
		clauses = append(clauses, fmt.Sprintf("rate_per_km_tenge = $%d", argNumber))
		args = append(args, *update.RatePerKMTenge)
		argNumber++
	}
	if update.FreeMinutes != nil {
		clauses = append(clauses, fmt.Sprintf("free_minutes = $%d", argNumber))
		args = append(args, *update.FreeMinutes)
		argNumber++
	}
	if update.MinChargeTenge != nil {
		clauses = append(clauses, fmt.Sprintf("min_charge_tenge = $%d", argNumber))
		args = append(args, *update.MinChargeTenge)
		argNumber++
	}
	if update.OvertimePolicy != nil {
		clauses = append(clauses, fmt.Sprintf("overtime_policy = $%d", argNumber))
		args = append(args, *update.OvertimePolicy)
		argNumber++
	}
	if update.OvertimeRateTenge != nil {
		clauses = append(clauses, fmt.Sprintf("overtime_rate_tenge = $%d", argNumber))
		args = append(args, *update.OvertimeRateTenge)
		argNumber++
	}
	if update.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", argNumber))
		args = append(args, *update.IsActive)
		argNumber++
	}

	clauses = append(clauses, "updated_at = NOW()")

	return clauses, args, argNumber
}
