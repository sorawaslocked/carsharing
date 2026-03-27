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

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

func (r *CarRepository) Insert(ctx context.Context, car model.Car) (string, error) {
	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO cars
			(model_id, vin, license_plate, color, year_manufactured, status,
			 mileage_km, fuel_level, battery_level, latitude, longitude,
			 notes, last_seen_at, created_at, updated_at)
		VALUES
			(` + strings.Join([]string{
		b.Add(car.ModelID),
		b.Add(car.VIN),
		b.Add(car.LicensePlate),
		b.Add(car.Color),
		b.Add(car.YearManufactured),
		b.Add(string(car.Status)),
		b.Add(car.MileageKM),
		b.Add(dto.NullableFloat32(car.FuelLevel)),
		b.Add(dto.NullableFloat32(car.BatteryLevel)),
		b.Add(car.Location.Latitude),
		b.Add(car.Location.Longitude),
		b.Add(pq.StringArray(car.Notes)),
		b.Add(car.LastSeenAt),
		b.Add(car.CreatedAt),
		b.Add(car.UpdatedAt),
	}, ", ") + `)
		RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return "", ErrAlreadyExists
		}
		return "", ErrSql
	}

	return id, nil
}

func (r *CarRepository) FindOne(ctx context.Context, filter model.CarFilter) (model.Car, error) {
	b := &dto.ArgsBuilder{}

	join, where := buildCarWhere(filter, b)

	q := `SELECT c.* FROM cars c` + join + where + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	car, err := dto.ScanCarRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Car{}, ErrNotFound
		}

		return model.Car{}, ErrSql
	}

	return car, nil
}

func (r *CarRepository) Find(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	b := &dto.ArgsBuilder{}

	join, where := buildCarWhere(filter, b)

	q := `SELECT c.* FROM cars c` + join + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.Car
	for rows.Next() {
		car, err := dto.ScanCarRow(rows)
		if err != nil {
			return nil, ErrSql
		}

		result = append(result, car)
	}

	err = rows.Err()
	if err != nil {
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarRepository) Update(ctx context.Context, filter model.CarFilter, update model.CarUpdate) error {
	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildCarSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}
	_, where := buildCarWhere(filter, b)

	q := `UPDATE cars SET ` + strings.Join(setClauses, ", ") + where

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CarRepository) Delete(ctx context.Context, filter model.CarFilter) error {
	b := &dto.ArgsBuilder{}

	_, where := buildCarWhere(filter, b)

	q := `DELETE FROM cars` + where

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func buildCarWhere(f model.CarFilter, b *dto.ArgsBuilder) (join string, where string) {
	var joins []string
	var clauses []string

	if f.ID != nil {
		clauses = append(clauses, fmt.Sprintf("c.id = %s", b.Add(*f.ID)))
	}
	if f.Status != nil {
		clauses = append(clauses, fmt.Sprintf("c.status = %s", b.Add(string(*f.Status))))
	}

	if f.ModelFilter != nil {
		joins = append(joins, "JOIN car_models cm ON cm.id = c.model_id")
		// Delegate to the car-model where builder with the same b so all
		// placeholder indices continue from wherever the car clauses left off.
		subClauses := dto.BuildCarModelWhereClauses(b, *f.ModelFilter, "cm")
		clauses = append(clauses, subClauses...)
	}

	if f.LocationFilter != nil {
		lf := f.LocationFilter
		// earth_box performs a bounding-cube pre-filter that is satisfied by the
		// GiST index, keeping the query efficient for large fleets.
		p1 := b.Add(lf.Location.Latitude)
		p2 := b.Add(lf.Location.Longitude)
		p3 := b.Add(lf.RadiusKM * 1000) // km → metres
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

	return join, where
}
