package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lib/pq"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"
)

type CarModelRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarModelRepository(db *sql.DB, log *slog.Logger) *CarModelRepository {
	return &CarModelRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarModelRepo"),
	}
}

func (r *CarModelRepository) Insert(ctx context.Context, cm model.CarModel) (string, error) {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO car_models
			(brand, model, year, fuel_type, transmission, body_type, class,
			 seats, engine_volume, range_km, features, image_keys, created_at, updated_at)
		VALUES
			(` + strings.Join([]string{
		b.Add(cm.Brand),
		b.Add(cm.Model),
		b.Add(cm.Year),
		b.Add(string(cm.FuelType)),
		b.Add(string(cm.Transmission)),
		b.Add(string(cm.BodyType)),
		b.Add(string(cm.Class)),
		b.Add(cm.Seats),
		b.Add(dto.NullableFloat32(cm.EngineVolume)),
		b.Add(cm.RangeKM),
		b.Add(pq.StringArray(cm.Features)),
		b.Add(pq.StringArray(dto.ImagesToKeys(cm.Images))),
		b.Add(cm.CreatedAt),
		b.Add(cm.UpdatedAt),
	}, ", ") + `)
		RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		logger.Error("failed to insert car model", pkglog.Err(err))
		return "", ErrSql
	}

	return id, nil
}

func (r *CarModelRepository) FindByID(ctx context.Context, id string) (model.CarModel, error) {
	logger := pkglog.WithMethod(r.log, "FindByID")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `SELECT id, brand, model, year, fuel_type, transmission, body_type, class,
		seats, engine_volume, range_km, features, image_keys, created_at, updated_at
		FROM car_models WHERE id = ` + b.Add(id) + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	cm, err := dto.ScanCarModelRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CarModel{}, ErrNotFound
		}
		logger.Error("failed to find car model by id", pkglog.Err(err))
		return model.CarModel{}, ErrSql
	}

	return cm, nil
}

func (r *CarModelRepository) Find(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	whereClauses := dto.BuildCarModelWhereClauses(b, filter, "")
	where := ""
	if len(whereClauses) > 0 {
		where = fmt.Sprintf(" WHERE %s", strings.Join(whereClauses, " AND "))
	}

	q := `SELECT id, brand, model, year, fuel_type, transmission, body_type, class,
		seats, engine_volume, range_km, features, image_keys, created_at, updated_at
		FROM car_models` + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query car models", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarModel
	for rows.Next() {
		cm, err := dto.ScanCarModelRow(rows)
		if err != nil {
			logger.Error("failed to scan car model row", pkglog.Err(err))
			return nil, ErrSql
		}
		result = append(result, cm)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarModelRepository) Update(ctx context.Context, id string, update model.CarModelUpdate) error {
	logger := pkglog.WithMethod(r.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildCarModelSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	q := `UPDATE car_models SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to update car model", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CarModelRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(r.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `DELETE FROM car_models WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to delete car model", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}
