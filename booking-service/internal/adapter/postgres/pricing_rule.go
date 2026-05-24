package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	pgdto "carsharing/booking-service/internal/adapter/postgres/dto"
	"carsharing/booking-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PricingRuleRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewPricingRuleRepository(log *slog.Logger, pool *pgxpool.Pool) *PricingRuleRepository {
	return &PricingRuleRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.PricingRuleRepository"),
		pool: pool,
	}
}

const pricingRuleSelect = `
    SELECT id, model_id, zone_id, class, type, rate_tenge, rate_per_km_tenge,
           free_minutes, min_charge_tenge, overtime_policy, overtime_rate_tenge,
           is_active, created_at, updated_at
    FROM pricing_rules`

func scanPricingRule(rs rowScanner) (model.PricingRule, error) {
	var rule model.PricingRule
	var modelID, zoneID, class, overtimePolicy *string
	var ratePerKM, freeMinutes, minCharge, overtimeRate *int32

	if err := rs.Scan(
		&rule.ID, &modelID, &zoneID, &class, &rule.Type, &rule.RateTenge,
		&ratePerKM, &freeMinutes, &minCharge, &overtimePolicy, &overtimeRate,
		&rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
	); err != nil {
		return model.PricingRule{}, err
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

func (r *PricingRuleRepository) Create(ctx context.Context, data model.PricingRuleCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Create"), utils.MetadataFromCtx(ctx))

	var id string
	err := r.pool.QueryRow(ctx, `
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
		log.Error("inserting pricing rule", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *PricingRuleRepository) GetByID(ctx context.Context, id string) (model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "GetByID"), utils.MetadataFromCtx(ctx))

	rule, err := scanPricingRule(r.pool.QueryRow(ctx, pricingRuleSelect+" WHERE id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PricingRule{}, model.ErrPricingRuleNotFound
		}
		log.Error("scanning pricing rule", pkglog.Err(err))
		return model.PricingRule{}, model.ErrSql
	}
	return rule, nil
}

func (r *PricingRuleRepository) List(ctx context.Context, filter model.PricingRuleListFilter) ([]model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "List"), utils.MetadataFromCtx(ctx))

	clauses, args, nextArg := pgdto.WhereClausesFromPricingRuleListFilter(filter, nil, 1)

	query := pricingRuleSelect
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", nextArg, nextArg+1)
	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("querying pricing rules", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var rules []model.PricingRule
	for rows.Next() {
		rule, err := scanPricingRule(rows)
		if err != nil {
			log.Error("scanning pricing rule row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		log.Error("iterating pricing rule rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return rules, nil
}

func (r *PricingRuleRepository) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, nextArg := pgdto.SetClausesFromPricingRuleUpdate(data)
	query := "UPDATE pricing_rules SET " + strings.Join(setClauses, ", ") +
		fmt.Sprintf(" WHERE id = $%d", nextArg)
	args = append(args, id)

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Error("updating pricing rule", pkglog.Err(err))
		return model.ErrSql
	}
	if tag.RowsAffected() == 0 {
		return model.ErrPricingRuleNotFound
	}

	return nil
}

func (r *PricingRuleRepository) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM pricing_rules WHERE id = $1`, id)
	if err != nil {
		log.Error("deleting pricing rule", pkglog.Err(err))
		return model.ErrSql
	}
	if tag.RowsAffected() == 0 {
		return model.ErrPricingRuleNotFound
	}

	return nil
}
