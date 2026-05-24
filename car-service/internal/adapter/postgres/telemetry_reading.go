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

type TelemetryReadingRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewTelemetryReadingRepository(log *slog.Logger, pool *pgxpool.Pool) *TelemetryReadingRepository {
	return &TelemetryReadingRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.TelemetryReadingRepository"),
		pool: pool,
	}
}

func (r *TelemetryReadingRepository) Insert(ctx context.Context, reading model.TelemetryReading) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	var lat, lon *float64
	if reading.Location != nil {
		lat = &reading.Location.Latitude
		lon = &reading.Location.Longitude
	}

	var metadataJSON []byte
	if reading.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(reading.Metadata)
		if err != nil {
			log.Error("failed to marshal telemetry reading metadata", pkglog.Err(err))
			return model.ErrSql
		}
	}

	args := []any{
		reading.CarID,
		lat, lon,
		reading.FuelPct, reading.FuelRawPct, reading.BatteryLevel, reading.MileageKM,
		reading.ActorID, string(reading.ActorType), reading.Reason,
		metadataJSON, reading.RecordedAt,
	}
	q := `INSERT INTO car_telemetry_readings
			(car_id, latitude, longitude, fuel_pct, fuel_raw_pct,
			 battery_level, mileage_km, actor_id, actor_type, reason, metadata, recorded_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	if _, err := r.pool.Exec(ctx, q, args...); err != nil {
		log.Error("failed to insert telemetry reading", pkglog.Err(err))
		return model.ErrSql
	}

	return nil
}

func (r *TelemetryReadingRepository) Find(ctx context.Context, filter model.TelemetryReadingFilter) ([]model.TelemetryReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	var clauses []string
	var args []any
	n := 0

	n++
	args = append(args, filter.CarID)
	clauses = append(clauses, fmt.Sprintf("car_id = $%d", n))
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

	q := `SELECT id, car_id, latitude, longitude, fuel_pct, fuel_raw_pct,
		battery_level, mileage_km, actor_id, actor_type, reason, metadata, recorded_at
		FROM car_telemetry_readings` + where + ` ORDER BY recorded_at DESC`

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
		log.Error("failed to query telemetry readings", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.TelemetryReading
	for rows.Next() {
		var t model.TelemetryReading
		var lat, lon *float64
		var actorType string
		var metadataRaw []byte

		if err = rows.Scan(
			&t.ID, &t.CarID, &lat, &lon,
			&t.FuelPct, &t.FuelRawPct, &t.BatteryLevel, &t.MileageKM,
			&t.ActorID, &actorType, &t.Reason,
			&metadataRaw, &t.RecordedAt,
		); err != nil {
			log.Error("failed to scan telemetry reading row", pkglog.Err(err))
			return nil, model.ErrSql
		}

		t.ActorType = sharedmodel.ActorType(actorType)

		if lat != nil && lon != nil {
			t.Location = &sharedmodel.Location{Latitude: *lat, Longitude: *lon}
		}

		if len(metadataRaw) > 0 {
			if err = json.Unmarshal(metadataRaw, &t.Metadata); err != nil {
				log.Error("failed to unmarshal telemetry reading metadata", pkglog.Err(err))
				return nil, model.ErrSql
			}
		}

		result = append(result, t)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}
