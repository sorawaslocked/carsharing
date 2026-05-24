package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"carsharing/car-service/internal/adapter/postgres/dto"
	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarInsuranceRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarInsuranceRepository(log *slog.Logger, pool *pgxpool.Pool) *CarInsuranceRepository {
	return &CarInsuranceRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarInsuranceRepository"),
		pool: pool,
	}
}

func (r *CarInsuranceRepository) Insert(ctx context.Context, ins model.CarInsurance) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	args := []any{
		ins.CarID, string(ins.Type), ins.Provider, ins.PolicyNum,
		ins.StartsAt, ins.ExpiresAt, ins.CostTenge, string(ins.Status),
		ins.Notes, dto.ImagesToKeys(ins.Images), ins.CreatedAt, ins.UpdatedAt,
	}
	q := `INSERT INTO car_insurances
			(car_id, type, provider, policy_num, starts_at, expires_at,
			 cost_tenge, status, notes, image_keys, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		  RETURNING id`

	var id string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&id); err != nil {
		log.Error("failed to insert car insurance", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *CarInsuranceRepository) FindByID(ctx context.Context, id string) (model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, car_id, type, provider, policy_num, starts_at, expires_at,
		cost_tenge, status, notes, image_keys, created_at, updated_at
		FROM car_insurances WHERE id = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, q, id)

	ins, err := dto.ScanCarInsuranceRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.CarInsurance{}, model.ErrCarInsuranceNotFound
		}
		log.Error("failed to find car insurance by id", pkglog.Err(err))
		return model.CarInsurance{}, model.ErrSql
	}

	return ins, nil
}

func (r *CarInsuranceRepository) Find(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	whereClauses, args, n := dto.WhereClausesFromCarInsuranceFilter(filter, make([]any, 0), 0)
	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	q := `SELECT id, car_id, type, provider, policy_num, starts_at, expires_at,
		cost_tenge, status, notes, image_keys, created_at, updated_at
		FROM car_insurances` + where

	if filter.Pagination != nil {
		n++
		args = append(args, filter.Pagination.Limit)
		q += fmt.Sprintf(" LIMIT $%d", n)
		n++
		args = append(args, filter.Pagination.Offset)
		q += fmt.Sprintf(" OFFSET $%d", n)
	}

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		log.Error("failed to query car insurances", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.CarInsurance
	for rows.Next() {
		ins, err := dto.ScanCarInsuranceRow(rows)
		if err != nil {
			log.Error("failed to scan car insurance row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		result = append(result, ins)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}

func (r *CarInsuranceRepository) Update(ctx context.Context, id string, update model.CarInsuranceUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, n := dto.SetClausesFromCarInsuranceUpdate(update)
	if len(setClauses) <= 1 {
		return nil
	}

	n++
	args = append(args, id)
	q := `UPDATE car_insurances SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", n)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		log.Error("failed to update car insurance", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrCarInsuranceNotFound
	}

	return nil
}

func (r *CarInsuranceRepository) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM car_insurances WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete car insurance", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrCarInsuranceNotFound
	}

	return nil
}
