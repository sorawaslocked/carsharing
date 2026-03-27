package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"strings"

	"github.com/lib/pq"
)

type CarModelRepository struct {
	db *sql.DB
}

func NewCarModelRepository(db *sql.DB) *CarModelRepository {
	return &CarModelRepository{db: db}
}

func (r *CarModelRepository) Insert(ctx context.Context, cm model.CarModel) (string, error) {
	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO car_models
			(brand, model, year, fuel_type, transmission, body_type, class,
			 seats, engine_volume, range_km, features, created_at, updated_at)
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
		b.Add(cm.CreatedAt),
		b.Add(cm.UpdatedAt),
	}, ", ") + `)
		RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		return "", ErrSql
	}

	return id, nil
}

func (r *CarModelRepository) FindOne(ctx context.Context, filter model.CarModelFilter) (model.CarModel, error) {
	b := &dto.ArgsBuilder{}

	whereClauses := dto.BuildCarModelWhereClauses(b, filter, "")
	var where string

	if len(whereClauses) == 0 {
		where = ""
	} else {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereClauses, " AND "))
	}

	q := `SELECT * FROM car_models` + where + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	cm, err := dto.ScanCarModelRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CarModel{}, ErrNotFound
		}

		return model.CarModel{}, ErrSql
	}

	return cm, nil
}

func (r *CarModelRepository) Find(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	b := &dto.ArgsBuilder{}
	whereClauses := dto.BuildCarModelWhereClauses(b, filter, "")
	var where string

	if len(whereClauses) == 0 {
		where = ""
	} else {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereClauses, " AND "))
	}

	q := `SELECT * FROM car_models` + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarModel
	for rows.Next() {
		cm, err := dto.ScanCarModelRow(rows)
		if err != nil {
			return nil, ErrSql
		}

		result = append(result, cm)
	}

	err = rows.Err()
	if err != nil {
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarModelRepository) Update(ctx context.Context, filter model.CarModelFilter, update model.CarModelUpdate) error {
	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildCarModelSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	whereClauses := dto.BuildCarModelWhereClauses(b, filter, "")
	var where string

	if len(whereClauses) == 0 {
		where = ""
	} else {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereClauses, " AND "))
	}

	q := `UPDATE car_models SET ` + strings.Join(setClauses, ", ") + where

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CarModelRepository) Delete(ctx context.Context, filter model.CarModelFilter) error {
	b := &dto.ArgsBuilder{}

	whereClauses := dto.BuildCarModelWhereClauses(b, filter, "")
	var where string

	if len(whereClauses) == 0 {
		where = ""
	} else {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereClauses, " AND "))
	}

	q := `DELETE FROM car_models` + where

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}
