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

type CarModelRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarModelRepository(log *slog.Logger, pool *pgxpool.Pool) *CarModelRepository {
	return &CarModelRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarModelRepository"),
		pool: pool,
	}
}

func (r *CarModelRepository) Insert(ctx context.Context, cm model.CarModel) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	args := []any{
		cm.Brand, cm.Model, cm.Year,
		string(cm.FuelType), string(cm.Transmission), string(cm.BodyType), string(cm.Class),
		cm.Seats, cm.EngineVolume, cm.RangeKM,
		cm.Features, dto.ImagesToKeys(cm.Images),
		cm.CreatedAt, cm.UpdatedAt,
	}
	q := `INSERT INTO car_models
			(brand, model, year, fuel_type, transmission, body_type, class,
			 seats, engine_volume, range_km, features, image_keys, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		  RETURNING id`

	var id string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&id); err != nil {
		log.Error("failed to insert car model", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *CarModelRepository) FindByID(ctx context.Context, id string) (model.CarModel, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, brand, model, year, fuel_type, transmission, body_type, class,
		seats, engine_volume, range_km, features, image_keys, created_at, updated_at
		FROM car_models WHERE id = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, q, id)

	cm, err := dto.ScanCarModelRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.CarModel{}, model.ErrNotFound
		}
		log.Error("failed to find car model by id", pkglog.Err(err))
		return model.CarModel{}, model.ErrSql
	}

	return cm, nil
}

func (r *CarModelRepository) Find(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	whereClauses, args, n := dto.WhereClausesFromCarModelFilter(filter, make([]any, 0), 0, "")
	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	q := `SELECT id, brand, model, year, fuel_type, transmission, body_type, class,
		seats, engine_volume, range_km, features, image_keys, created_at, updated_at
		FROM car_models` + where

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
		log.Error("failed to query car models", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.CarModel
	for rows.Next() {
		cm, err := dto.ScanCarModelRow(rows)
		if err != nil {
			log.Error("failed to scan car model row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		result = append(result, cm)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}

func (r *CarModelRepository) Update(ctx context.Context, id string, update model.CarModelUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, n := dto.SetClausesFromCarModelUpdate(update)
	if len(setClauses) <= 1 {
		return nil
	}

	n++
	args = append(args, id)
	q := `UPDATE car_models SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", n)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		log.Error("failed to update car model", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *CarModelRepository) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM car_models WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete car model", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
