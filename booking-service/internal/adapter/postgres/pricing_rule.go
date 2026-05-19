package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"carsharing/booking-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PricingRuleRepo struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewPricingRuleRepo(log *slog.Logger, db *pgxpool.Pool) *PricingRuleRepo {
	return &PricingRuleRepo{
		log: pkglog.WithComponent(log, "repo.PricingRuleRepo"),
		db:  db,
	}
}

func (r *PricingRuleRepo) Create(ctx context.Context, data model.PricingRuleCreate) (string, error) {
	log := pkglog.WithMethod(r.log, "Create")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO pricing_rules
			(model_id, zone_id, class, type, rate_tenge, rate_per_km_tenge, free_minutes, min_charge_tenge, overtime_policy, overtime_rate_tenge)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		data.ModelID, data.ZoneID, data.Class, data.Type, data.RateTenge,
		data.RatePerKMTenge, data.FreeMinutes, data.MinChargeTenge,
		data.OvertimePolicy, data.OvertimeRateTenge,
	).Scan(&id)
	if err != nil {
		log.Error("failed to insert pricing rule", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	return id, nil
}

func (r *PricingRuleRepo) GetByID(ctx context.Context, id string) (model.PricingRule, error) {
	log := pkglog.WithMethod(r.log, "GetByID")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	var rule model.PricingRule
	var modelID, zoneID, class, overtimePolicy *string
	var ratePerKM, freeMinutes, minCharge, overtimeRate *int32

	err := r.db.QueryRow(ctx, `
		SELECT id, model_id, zone_id, class, type, rate_tenge, rate_per_km_tenge,
		       free_minutes, min_charge_tenge, overtime_policy, overtime_rate_tenge,
		       is_active, created_at, updated_at
		FROM pricing_rules WHERE id = $1
	`, id).Scan(
		&rule.ID, &modelID, &zoneID, &class, &rule.Type, &rule.RateTenge,
		&ratePerKM, &freeMinutes, &minCharge, &overtimePolicy, &overtimeRate,
		&rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.PricingRule{}, model.ErrNotFound
	}
	if err != nil {
		log.Error("failed to query pricing rule", pkglog.Err(err))
		return model.PricingRule{}, model.ErrInternalServerError
	}

	rule.ModelID = modelID
	rule.ZoneID = zoneID
	rule.Class = class
	rule.OvertimePolicy = overtimePolicy
	rule.RatePerKMTenge = ratePerKM
	rule.FreeMinutes = freeMinutes
	rule.MinChargeTenge = minCharge
	rule.OvertimeRateTenge = overtimeRate

	return rule, nil
}

func (r *PricingRuleRepo) List(ctx context.Context, filter model.PricingRuleListFilter) ([]model.PricingRule, error) {
	log := pkglog.WithMethod(r.log, "List")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if filter.ModelID != nil {
		where = append(where, fmt.Sprintf("model_id = $%d", idx))
		args = append(args, *filter.ModelID)
		idx++
	}
	if filter.ZoneID != nil {
		where = append(where, fmt.Sprintf("zone_id = $%d", idx))
		args = append(args, *filter.ZoneID)
		idx++
	}
	if filter.Class != nil {
		where = append(where, fmt.Sprintf("class = $%d", idx))
		args = append(args, *filter.Class)
		idx++
	}
	if filter.Type != nil {
		where = append(where, fmt.Sprintf("type = $%d", idx))
		args = append(args, *filter.Type)
		idx++
	}
	if filter.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *filter.IsActive)
		idx++
	}

	query := fmt.Sprintf(`
		SELECT id, model_id, zone_id, class, type, rate_tenge, rate_per_km_tenge,
		       free_minutes, min_charge_tenge, overtime_policy, overtime_rate_tenge,
		       is_active, created_at, updated_at
		FROM pricing_rules WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), idx, idx+1)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to list pricing rules", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}
	defer rows.Close()

	var rules []model.PricingRule
	for rows.Next() {
		var rule model.PricingRule
		var modelID, zoneID, class, overtimePolicy *string
		var ratePerKM, freeMinutes, minCharge, overtimeRate *int32

		if err := rows.Scan(
			&rule.ID, &modelID, &zoneID, &class, &rule.Type, &rule.RateTenge,
			&ratePerKM, &freeMinutes, &minCharge, &overtimePolicy, &overtimeRate,
			&rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
		); err != nil {
			log.Error("failed to scan pricing rule row", pkglog.Err(err))
			return nil, model.ErrInternalServerError
		}

		rule.ModelID = modelID
		rule.ZoneID = zoneID
		rule.Class = class
		rule.OvertimePolicy = overtimePolicy
		rule.RatePerKMTenge = ratePerKM
		rule.FreeMinutes = freeMinutes
		rule.MinChargeTenge = minCharge
		rule.OvertimeRateTenge = overtimeRate

		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}

	return rules, nil
}

func (r *PricingRuleRepo) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	log := pkglog.WithMethod(r.log, "Update")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	sets := []string{"updated_at = NOW()"}
	args := []any{}
	idx := 1

	if data.ModelID != nil {
		sets = append(sets, fmt.Sprintf("model_id = $%d", idx))
		args = append(args, *data.ModelID)
		idx++
	}
	if data.ZoneID != nil {
		sets = append(sets, fmt.Sprintf("zone_id = $%d", idx))
		args = append(args, *data.ZoneID)
		idx++
	}
	if data.Class != nil {
		sets = append(sets, fmt.Sprintf("class = $%d", idx))
		args = append(args, *data.Class)
		idx++
	}
	if data.Type != nil {
		sets = append(sets, fmt.Sprintf("type = $%d", idx))
		args = append(args, *data.Type)
		idx++
	}
	if data.RateTenge != nil {
		sets = append(sets, fmt.Sprintf("rate_tenge = $%d", idx))
		args = append(args, *data.RateTenge)
		idx++
	}
	if data.RatePerKMTenge != nil {
		sets = append(sets, fmt.Sprintf("rate_per_km_tenge = $%d", idx))
		args = append(args, *data.RatePerKMTenge)
		idx++
	}
	if data.FreeMinutes != nil {
		sets = append(sets, fmt.Sprintf("free_minutes = $%d", idx))
		args = append(args, *data.FreeMinutes)
		idx++
	}
	if data.MinChargeTenge != nil {
		sets = append(sets, fmt.Sprintf("min_charge_tenge = $%d", idx))
		args = append(args, *data.MinChargeTenge)
		idx++
	}
	if data.OvertimePolicy != nil {
		sets = append(sets, fmt.Sprintf("overtime_policy = $%d", idx))
		args = append(args, *data.OvertimePolicy)
		idx++
	}
	if data.OvertimeRateTenge != nil {
		sets = append(sets, fmt.Sprintf("overtime_rate_tenge = $%d", idx))
		args = append(args, *data.OvertimeRateTenge)
		idx++
	}
	if data.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *data.IsActive)
		idx++
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE pricing_rules SET %s WHERE id = $%d`, strings.Join(sets, ", "), idx)

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to update pricing rule", pkglog.Err(err))
		return model.ErrInternalServerError
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *PricingRuleRepo) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMethod(r.log, "Delete")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	tag, err := r.db.Exec(ctx, `DELETE FROM pricing_rules WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete pricing rule", pkglog.Err(err))
		return model.ErrInternalServerError
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
