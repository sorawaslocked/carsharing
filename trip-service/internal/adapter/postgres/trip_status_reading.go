package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"carsharing/trip-service/internal/adapter/postgres/dto"
	"carsharing/trip-service/internal/model"
	pkglog "carsharing/trip-service/internal/pkg/log"
	"carsharing/trip-service/internal/pkg/utils"
	"github.com/google/uuid"
)

type TripStatusReadingRepo struct {
	log *slog.Logger
	db  *sql.DB
}

func NewTripStatusReadingRepo(log *slog.Logger, db *sql.DB) *TripStatusReadingRepo {
	return &TripStatusReadingRepo{
		log: pkglog.WithComponent(log, "repo.TripStatusReadingRepo"),
		db:  db,
	}
}

func (r *TripStatusReadingRepo) Create(ctx context.Context, reading model.TripStatusReadingCreate) (model.TripStatusReading, error) {
	log := pkglog.WithMethod(r.log, "Create")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	q := `
		INSERT INTO trip_status_readings (
			id, trip_id, from_status, to_status, actor_type, actor_id, reason, changed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, trip_id, from_status, to_status, actor_type, actor_id, reason, changed_at`

	result, err := dto.ScanTripStatusReading(r.db.QueryRowContext(ctx, q,
		uuid.New().String(),
		reading.TripID,
		reading.FromStatus.String(), reading.ToStatus.String(),
		reading.ActorType.String(),
		dto.NullableString(reading.ActorID),
		dto.NullableString(reading.Reason),
		reading.ChangedAt,
	))
	if err != nil {
		return model.TripStatusReading{}, mapSQLError(log, err, "failed to create status reading")
	}
	return result, nil
}

func (r *TripStatusReadingRepo) List(ctx context.Context, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	log := pkglog.WithMethod(r.log, "List")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

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

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		log.Error("failed to list status readings", pkglog.Err(err))
		return nil, model.ErrSQL
	}
	defer rows.Close()

	var readings []model.TripStatusReading
	for rows.Next() {
		reading, err := dto.ScanTripStatusReading(rows)
		if err != nil {
			log.Error("failed to scan status reading row", pkglog.Err(err))
			return nil, model.ErrSQL
		}
		readings = append(readings, reading)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSQL
	}

	return readings, nil
}
