package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"
)

type TelematicsRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewTelematicsRepository(db *sql.DB, log *slog.Logger) *TelematicsRepository {
	return &TelematicsRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.TelematicsRepo"),
	}
}

func (r *TelematicsRepository) InsertEvent(ctx context.Context, event model.CarTelematicsEvent) error {
	logger := pkglog.WithMethod(r.log, "InsertEvent")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	var metadataJSON []byte
	if event.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(event.Metadata)
		if err != nil {
			logger.Error("failed to marshal telematics event metadata", pkglog.Err(err))
			return ErrSql
		}
	}

	q := `
		INSERT INTO car_telematics_events
			(car_id, latitude, longitude, fuel_level, battery_level,
			 odometer_km, actor_id, actor_type, metadata, recorded_at, received_at)
		VALUES (` + strings.Join([]string{
		b.Add(event.CarID),
		b.Add(event.Latitude),
		b.Add(event.Longitude),
		b.Add(dto.NullableFloat32(event.FuelLevel)),
		b.Add(dto.NullableFloat32(event.BatteryLevel)),
		b.Add(event.OdometerKM),
		b.Add(event.ActorID),
		b.Add(event.ActorType),
		b.Add(metadataJSON),
		b.Add(event.RecordedAt),
		b.Add(event.ReceivedAt),
	}, ", ") + `)`

	_, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to insert telematics event", pkglog.Err(err))
		return ErrSql
	}

	return nil
}

func (r *TelematicsRepository) FindEvents(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error) {
	logger := pkglog.WithMethod(r.log, "FindEvents")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	var clauses []string
	if filter.CarID != nil {
		clauses = append(clauses, "car_id = "+b.Add(*filter.CarID))
	}
	if filter.From != nil {
		clauses = append(clauses, "recorded_at >= "+b.Add(*filter.From))
	}
	if filter.To != nil {
		clauses = append(clauses, "recorded_at <= "+b.Add(*filter.To))
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, car_id, latitude, longitude, fuel_level, battery_level,
		odometer_km, actor_id, actor_type, metadata, recorded_at, received_at
		FROM car_telematics_events` + where + ` ORDER BY recorded_at DESC` + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query telematics events", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarTelematicsEvent
	for rows.Next() {
		var e model.CarTelematicsEvent
		var fuelLevel, batteryLevel sql.NullFloat64
		var actorID sql.NullString
		var metadataRaw []byte

		if err = rows.Scan(
			&e.ID, &e.CarID, &e.Latitude, &e.Longitude,
			&fuelLevel, &batteryLevel, &e.OdometerKM,
			&actorID, &e.ActorType, &metadataRaw,
			&e.RecordedAt, &e.ReceivedAt,
		); err != nil {
			logger.Error("failed to scan telematics event row", pkglog.Err(err))
			return nil, ErrSql
		}

		if fuelLevel.Valid {
			f := float32(fuelLevel.Float64)
			e.FuelLevel = &f
		}
		if batteryLevel.Valid {
			bv := float32(batteryLevel.Float64)
			e.BatteryLevel = &bv
		}
		if actorID.Valid {
			e.ActorID = &actorID.String
		}
		if len(metadataRaw) > 0 {
			if err = json.Unmarshal(metadataRaw, &e.Metadata); err != nil {
				logger.Error("failed to unmarshal telematics event metadata", pkglog.Err(err))
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
