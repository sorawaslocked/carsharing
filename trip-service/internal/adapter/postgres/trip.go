package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/postgres/dto"
	"carsharing/trip-service/internal/model"
)

type TripRepo struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewTripRepo(log *slog.Logger, pool *pgxpool.Pool) *TripRepo {
	return &TripRepo{
		log:  pkglog.WithComponent(log, "repo.TripRepo"),
		pool: pool,
	}
}

func (r *TripRepo) Create(ctx context.Context, trip model.TripCreate) (model.Trip, error) {
	log := pkglog.WithMethod(r.log, "Create")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	now := time.Now()

	q := fmt.Sprintf(`
		INSERT INTO trips (
			id, booking_id, user_id, car_id, status,
			started_at, start_latitude, start_longitude, start_mileage_km, start_fuel_level,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING %s`, dto.TripColumns)

	t, err := dto.ScanTrip(r.pool.QueryRow(ctx, q,
		trip.ID, trip.BookingID, trip.UserID, trip.CarID, trip.Status.String(),
		trip.StartedAt,
		trip.StartLocation.Latitude, trip.StartLocation.Longitude,
		trip.StartMileageKM, trip.StartFuelLevel,
		now, now,
	))
	if err != nil {
		return model.Trip{}, mapSQLError(log, err, "failed to create trip")
	}
	return t, nil
}

func (r *TripRepo) GetByID(ctx context.Context, id string) (model.Trip, error) {
	log := pkglog.WithMethod(r.log, "GetByID")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	q := fmt.Sprintf(`SELECT %s FROM trips WHERE id = $1`, dto.TripColumns)

	t, err := dto.ScanTrip(r.pool.QueryRow(ctx, q, id))
	if err != nil {
		return model.Trip{}, mapSQLError(log, err, "failed to get trip by id")
	}
	return t, nil
}

func (r *TripRepo) List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error) {
	log := pkglog.WithMethod(r.log, "List")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}
	where := dto.BuildTripWhereClause(dto.BuildTripWhereClauses(filter, b))
	pagination := dto.BuildPagination(filter.Pagination, b)

	q := fmt.Sprintf(
		`SELECT %s FROM trips %s ORDER BY created_at DESC%s`,
		dto.TripColumns, where, pagination,
	)

	rows, err := r.pool.Query(ctx, q, b.Args...)
	if err != nil {
		log.Error("failed to list trips", pkglog.Err(err))
		return nil, model.ErrSQL
	}
	defer rows.Close()

	var trips []model.Trip
	for rows.Next() {
		t, err := dto.ScanTrip(rows)
		if err != nil {
			log.Error("failed to scan trip row", pkglog.Err(err))
			return nil, model.ErrSQL
		}
		trips = append(trips, t)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSQL
	}

	return trips, nil
}

func (r *TripRepo) Update(ctx context.Context, id string, update model.TripUpdate) (model.Trip, error) {
	log := pkglog.WithMethod(r.log, "Update")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}
	setClauses := dto.BuildTripSetClauses(update, b)

	q := fmt.Sprintf(
		`UPDATE trips SET %s WHERE id = %s RETURNING %s`,
		strings.Join(setClauses, ", "),
		b.Add(id),
		dto.TripColumns,
	)

	t, err := dto.ScanTrip(r.pool.QueryRow(ctx, q, b.Args...))
	if err != nil {
		return model.Trip{}, mapSQLError(log, err, "failed to update trip")
	}
	return t, nil
}
