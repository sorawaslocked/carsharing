package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarStatusReadingRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarStatusReadingRepository(log *slog.Logger, pool *pgxpool.Pool) *CarStatusReadingRepository {
	return &CarStatusReadingRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarStatusReadingRepository"),
		pool: pool,
	}
}

func (r *CarStatusReadingRepository) Insert(ctx context.Context, entry model.CarStatusReading) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	var metadataJSON []byte
	if entry.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(entry.Metadata)
		if err != nil {
			log.Error("failed to marshal status reading metadata", pkglog.Err(err))
			return model.ErrSql
		}
	}

	args := []any{
		entry.CarID,
		string(entry.FromStatus), string(entry.ToStatus),
		string(entry.ActorType), entry.ActorID, entry.Reason,
		metadataJSON, entry.RecordedAt,
	}
	q := `INSERT INTO car_status_readings
			(car_id, from_status, to_status, actor_type, actor_id, reason, metadata, recorded_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	if _, err := r.pool.Exec(ctx, q, args...); err != nil {
		log.Error("failed to insert car status reading", pkglog.Err(err))
		return model.ErrSql
	}

	return nil
}

func (r *CarStatusReadingRepository) Find(ctx context.Context, filter model.CarStatusReadingFilter) ([]model.CarStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	var clauses []string
	var args []any
	n := 0

	n++
	args = append(args, filter.CarID)
	clauses = append(clauses, fmt.Sprintf("car_id = $%d", n))
	if filter.FromStatus != nil {
		n++
		args = append(args, string(*filter.FromStatus))
		clauses = append(clauses, fmt.Sprintf("from_status = $%d", n))
	}
	if filter.ToStatus != nil {
		n++
		args = append(args, string(*filter.ToStatus))
		clauses = append(clauses, fmt.Sprintf("to_status = $%d", n))
	}
	if filter.TimeRange != nil {
		if !filter.TimeRange.From.IsZero() {
			n++
			args = append(args, filter.TimeRange.From)
			clauses = append(clauses, fmt.Sprintf("recorded_at >= $%d", n))
		}
		if !filter.TimeRange.To.IsZero() {
			n++
			args = append(args, filter.TimeRange.To)
			clauses = append(clauses, fmt.Sprintf("recorded_at <= $%d", n))
		}
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, car_id, from_status, to_status, actor_type, actor_id, reason, metadata, recorded_at
		FROM car_status_readings` + where + ` ORDER BY recorded_at DESC`

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
		log.Error("failed to query car status readings", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.CarStatusReading
	for rows.Next() {
		var e model.CarStatusReading
		var fromStatus, toStatus, actorType string
		var metadataRaw []byte

		if err = rows.Scan(
			&e.ID, &e.CarID, &fromStatus, &toStatus,
			&actorType, &e.ActorID, &e.Reason,
			&metadataRaw, &e.RecordedAt,
		); err != nil {
			log.Error("failed to scan status reading row", pkglog.Err(err))
			return nil, model.ErrSql
		}

		e.FromStatus = model.CarStatus(fromStatus)
		e.ToStatus = model.CarStatus(toStatus)
		e.ActorType = sharedmodel.ActorType(actorType)

		if len(metadataRaw) > 0 {
			if err = json.Unmarshal(metadataRaw, &e.Metadata); err != nil {
				log.Error("failed to unmarshal status reading metadata", pkglog.Err(err))
				return nil, model.ErrSql
			}
		}

		result = append(result, e)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}
