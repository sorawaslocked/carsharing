package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strings"

	"carsharing/car-service/internal/adapter/postgres/dto"
	"carsharing/car-service/internal/model"
	pkglog "carsharing/car-service/internal/pkg/log"
	"carsharing/car-service/internal/pkg/utils"
)

type CarStatusLogRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarStatusLogRepository(db *sql.DB, log *slog.Logger) *CarStatusLogRepository {
	return &CarStatusLogRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarStatusLogRepo"),
	}
}

func (r *CarStatusLogRepository) Insert(ctx context.Context, entry model.CarStatusLogEntry) error {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	var metadataJSON []byte
	if entry.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(entry.Metadata)
		if err != nil {
			logger.Error("failed to marshal status log metadata", pkglog.Err(err))
			return ErrSql
		}
	}

	q := `
		INSERT INTO car_status_logs
			(car_id, from_status, to_status, actor_type, actor_id, reason, metadata, changed_at)
		VALUES (` + strings.Join([]string{
		b.Add(entry.CarID),
		b.Add(string(entry.FromStatus)),
		b.Add(string(entry.ToStatus)),
		b.Add(string(entry.ActorType)),
		b.Add(entry.ActorID),
		b.Add(entry.Reason),
		b.Add(metadataJSON),
		b.Add(entry.ChangedAt),
	}, ", ") + `)`

	_, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to insert car status log entry", pkglog.Err(err))
		return ErrSql
	}

	return nil
}

func (r *CarStatusLogRepository) Find(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	var clauses []string
	if filter.CarID != nil {
		clauses = append(clauses, "car_id = "+b.Add(*filter.CarID))
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, car_id, from_status, to_status, actor_type, actor_id, reason, metadata, changed_at
		FROM car_status_logs` + where + ` ORDER BY changed_at DESC` + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query car status logs", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarStatusLogEntry
	for rows.Next() {
		var e model.CarStatusLogEntry
		var fromStatus, toStatus, actorType string
		var actorID, reason sql.NullString
		var metadataRaw []byte

		if err = rows.Scan(
			&e.ID, &e.CarID, &fromStatus, &toStatus, &actorType,
			&actorID, &reason, &metadataRaw, &e.ChangedAt,
		); err != nil {
			logger.Error("failed to scan status log row", pkglog.Err(err))
			return nil, ErrSql
		}

		e.FromStatus = model.CarStatus(fromStatus)
		e.ToStatus = model.CarStatus(toStatus)
		e.ActorType = model.CarStatusActor(actorType)
		if actorID.Valid {
			e.ActorID = &actorID.String
		}
		if reason.Valid {
			e.Reason = &reason.String
		}
		if len(metadataRaw) > 0 {
			if err = json.Unmarshal(metadataRaw, &e.Metadata); err != nil {
				logger.Error("failed to unmarshal status log metadata", pkglog.Err(err))
				return nil, ErrSql
			}
		}

		result = append(result, e)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}
