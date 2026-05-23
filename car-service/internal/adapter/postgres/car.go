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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarRepository(log *slog.Logger, pool *pgxpool.Pool) *CarRepository {
	return &CarRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarRepository"),
		pool: pool,
	}
}

func (r *CarRepository) handlePGErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return model.ErrAlreadyExists
	}
	return model.ErrSql
}

func (r *CarRepository) Insert(ctx context.Context, car model.Car) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	args := []any{
		car.ModelID, car.VIN, car.LicensePlate, car.Color,
		car.YearManufactured, string(car.Status),
		car.MileageKM, car.FuelLevel, car.BatteryLevel,
		car.Location.Latitude, car.Location.Longitude,
		car.Notes, dto.ImagesToKeys(car.Images), car.LastSeenAt,
		car.CreatedAt, car.UpdatedAt,
	}
	q := `INSERT INTO cars
			(model_id, vin, license_plate, color, year_manufactured, status,
			 mileage_km, fuel_level, battery_level, latitude, longitude,
			 notes, image_keys, last_seen_at, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		  RETURNING id`

	var id string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&id); err != nil {
		log.Error("failed to insert car", pkglog.Err(err))
		return "", r.handlePGErr(err)
	}

	return id, nil
}

func (r *CarRepository) FindByID(ctx context.Context, id string) (model.Car, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, model_id, vin, license_plate, color, year_manufactured, status,
		mileage_km, fuel_level, battery_level, latitude, longitude, notes, image_keys,
		last_seen_at, created_at, updated_at FROM cars WHERE id = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, q, id)

	car, err := dto.ScanCarRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Car{}, model.ErrNotFound
		}
		log.Error("failed to find car by id", pkglog.Err(err))
		return model.Car{}, model.ErrSql
	}

	return car, nil
}

func (r *CarRepository) Find(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	join, where, args, n := buildCarFilter(filter, make([]any, 0), 0)

	q := `SELECT c.id, c.model_id, c.vin, c.license_plate, c.color, c.year_manufactured, c.status,
		c.mileage_km, c.fuel_level, c.battery_level, c.latitude, c.longitude, c.notes, c.image_keys,
		c.last_seen_at, c.created_at, c.updated_at FROM cars c` + join + where

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
		log.Error("failed to query cars", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.Car
	for rows.Next() {
		car, err := dto.ScanCarRow(rows)
		if err != nil {
			log.Error("failed to scan car row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		result = append(result, car)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}

func (r *CarRepository) Update(ctx context.Context, id string, update model.CarUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, n := dto.SetClausesFromCarUpdate(update)
	if len(setClauses) <= 1 {
		return nil
	}

	n++
	args = append(args, id)
	q := `UPDATE cars SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", n)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		log.Error("failed to update car", pkglog.Err(err))
		return r.handlePGErr(err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *CarRepository) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM cars WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete car", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func buildCarFilter(f model.CarFilter, args []any, n int) (join, where string, outArgs []any, outN int) {
	var joins []string
	var clauses []string

	if f.ID != nil {
		n++
		args = append(args, *f.ID)
		clauses = append(clauses, fmt.Sprintf("c.id = $%d", n))
	}
	if f.Status != nil {
		n++
		args = append(args, string(*f.Status))
		clauses = append(clauses, fmt.Sprintf("c.status = $%d", n))
	}

	if f.ModelFilter != nil {
		joins = append(joins, "JOIN car_models cm ON cm.id = c.model_id")
		var modelClauses []string
		modelClauses, args, n = dto.WhereClausesFromCarModelFilter(*f.ModelFilter, args, n, "cm")
		clauses = append(clauses, modelClauses...)
	}

	if f.LocationFilter != nil {
		lf := f.LocationFilter
		n++
		args = append(args, lf.Location.Latitude)
		p1 := fmt.Sprintf("$%d", n)
		n++
		args = append(args, lf.Location.Longitude)
		p2 := fmt.Sprintf("$%d", n)
		n++
		args = append(args, lf.RadiusKM*1000)
		p3 := fmt.Sprintf("$%d", n)
		clauses = append(clauses, fmt.Sprintf(
			"earth_box(ll_to_earth(%s, %s), %s) @> ll_to_earth(c.latitude, c.longitude)",
			p1, p2, p3,
		))
	}

	if len(joins) > 0 {
		join = " " + strings.Join(joins, " ")
	}
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	return join, where, args, n
}
