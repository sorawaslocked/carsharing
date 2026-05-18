package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"carsharing/car-service/internal/adapter/postgres/dto"
	"carsharing/car-service/internal/model"
	pkglog "carsharing/car-service/internal/pkg/log"
	"carsharing/car-service/internal/pkg/utils"
	"github.com/lib/pq"
)

type CarRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarRepository(db *sql.DB, log *slog.Logger) *CarRepository {
	return &CarRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarRepo"),
	}
}

func (r *CarRepository) Insert(ctx context.Context, car model.Car) (string, error) {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO cars
			(model_id, vin, license_plate, color, year_manufactured, status,
			 mileage_km, fuel_level, battery_level, latitude, longitude,
			 notes, image_keys, last_seen_at, created_at, updated_at)
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
		b.Add(pq.StringArray(dto.ImagesToKeys(car.Images))),
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
		logger.Error("failed to insert car", pkglog.Err(err))
		return "", ErrSql
	}

	return id, nil
}

func (r *CarRepository) FindByID(ctx context.Context, id string) (model.Car, error) {
	logger := pkglog.WithMethod(r.log, "FindByID")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `SELECT id, model_id, vin, license_plate, color, year_manufactured, status,
		mileage_km, fuel_level, battery_level, latitude, longitude, notes, image_keys,
		last_seen_at, created_at, updated_at FROM cars WHERE id = ` + b.Add(id) + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	car, err := dto.ScanCarRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Car{}, ErrNotFound
		}
		logger.Error("failed to find car by id", pkglog.Err(err))
		return model.Car{}, ErrSql
	}

	return car, nil
}

func (r *CarRepository) Find(ctx context.Context, filter model.CarFilter) ([]model.Car, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	join, where := buildCarWhere(filter, b)

	q := `SELECT c.id, c.model_id, c.vin, c.license_plate, c.color, c.year_manufactured, c.status,
		c.mileage_km, c.fuel_level, c.battery_level, c.latitude, c.longitude, c.notes, c.image_keys,
		c.last_seen_at, c.created_at, c.updated_at FROM cars c` + join + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query cars", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.Car
	for rows.Next() {
		car, err := dto.ScanCarRow(rows)
		if err != nil {
			logger.Error("failed to scan car row", pkglog.Err(err))
			return nil, ErrSql
		}
		result = append(result, car)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarRepository) Update(ctx context.Context, id string, update model.CarUpdate) error {
	logger := pkglog.WithMethod(r.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildCarSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	q := `UPDATE cars SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		logger.Error("failed to update car", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CarRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(r.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `DELETE FROM cars WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to delete car", pkglog.Err(err))
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
		subClauses := dto.BuildCarModelWhereClauses(b, *f.ModelFilter, "cm")
		clauses = append(clauses, subClauses...)
	}

	if f.LocationFilter != nil {
		lf := f.LocationFilter
		// earth_box performs a bounding-cube pre-filter satisfied by the GiST index.
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
