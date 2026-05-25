package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/postgres/dto"
	"carsharing/trip-service/internal/model"
)

type TripStatusReadingRepo struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewTripStatusReadingRepo(log *slog.Logger, pool *pgxpool.Pool) *TripStatusReadingRepo {
	return &TripStatusReadingRepo{
		log:  pkglog.WithComponent(log, "adapter.postgres.TripStatusReadingRepository"),
		pool: pool,
	}
}

func (r *TripStatusReadingRepo) Create(ctx context.Context, reading model.TripStatusReadingCreate) (model.TripStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Create"), pkgutils.MetadataFromCtx(ctx))

	q := `
		INSERT INTO trip_status_readings (
			trip_id, from_status, to_status, actor_type, actor_id, reason, changed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, trip_id, from_status, to_status, actor_type, actor_id, reason, changed_at`

	result, err := dto.ScanTripStatusReading(dbFromCtx(ctx, r.pool).QueryRow(ctx, q,
		reading.TripID,
		reading.FromStatus.String(), reading.ToStatus.String(),
		string(reading.ActorType),
		reading.ActorID,
		reading.Reason,
		reading.ChangedAt,
	))
	if err != nil {
		return model.TripStatusReading{}, mapSQLError(log, err, "creating status reading")
	}
	return result, nil
}

func (r *TripStatusReadingRepo) List(ctx context.Context, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "List"), pkgutils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}
	whereClauses := dto.BuildStatusReadingWhereClauses(filter, b)
	pagination := dto.BuildPagination(filter.Pagination, b)

	q := fmt.Sprintf(`
		SELECT id, trip_id, from_status, to_status, actor_type, actor_id, reason, changed_at
		FROM trip_status_readings
		WHERE %s
		ORDER BY changed_at ASC%s`,
		strings.Join(whereClauses, " AND "),
		pagination,
	)

	rows, err := dbFromCtx(ctx, r.pool).Query(ctx, q, b.Args...)
	if err != nil {
		log.Error("listing status readings", pkglog.Err(err))
		return nil, model.ErrSQL
	}
	defer rows.Close()

	readings := []model.TripStatusReading{}
	for rows.Next() {
		reading, err := dto.ScanTripStatusReading(rows)
		if err != nil {
			log.Error("scanning status reading row", pkglog.Err(err))
			return nil, model.ErrSQL
		}
		readings = append(readings, reading)
	}
	if err := rows.Err(); err != nil {
		log.Error("iterating status reading rows", pkglog.Err(err))
		return nil, model.ErrSQL
	}

	return readings, nil
}
