package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-trip-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-trip-service/internal/pkg/utils"
)

type TripSummaryRepo struct {
	log *slog.Logger
	db  *sql.DB
}

func NewTripSummaryRepo(log *slog.Logger, db *sql.DB) *TripSummaryRepo {
	return &TripSummaryRepo{
		log: pkglog.WithComponent(log, "repo.TripSummaryRepo"),
		db:  db,
	}
}

func (r *TripSummaryRepo) Create(ctx context.Context, s model.TripSummaryCreate) (model.TripSummary, error) {
	log := pkglog.WithMethod(r.log, "Create")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	snapJSON, err := dto.MarshalPricingSnapshot(s.PricingSnapshot)
	if err != nil {
		log.Error("failed to marshal pricing snapshot", pkglog.Err(err))
		return model.TripSummary{}, model.ErrSQL
	}

	q := fmt.Sprintf(`
		INSERT INTO trip_summaries (
			trip_id, booking_id, started_at, ended_at,
			duration_seconds, distance_traveled_km, pricing_snapshot,
			base_cost_tenge, distance_cost_tenge, overtime_cost_tenge, total_cost_tenge
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING %s`, dto.TripSummaryColumns)

	summary, err := dto.ScanTripSummary(r.db.QueryRowContext(ctx, q,
		s.TripID, s.BookingID, s.StartedAt, s.EndedAt,
		s.DurationSeconds, s.DistanceTraveledKM, snapJSON,
		s.BaseCostTenge, s.DistanceCostTenge, s.OvertimeCostTenge, s.TotalCostTenge,
	))
	if err != nil {
		return model.TripSummary{}, mapSQLError(log, err, "failed to create trip summary")
	}
	return summary, nil
}

func (r *TripSummaryRepo) GetByTripID(ctx context.Context, tripID string) (model.TripSummary, error) {
	log := pkglog.WithMethod(r.log, "GetByTripID")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	q := fmt.Sprintf(`SELECT %s FROM trip_summaries WHERE trip_id = $1`, dto.TripSummaryColumns)

	summary, err := dto.ScanTripSummary(r.db.QueryRowContext(ctx, q, tripID))
	if err != nil {
		return model.TripSummary{}, mapSQLError(log, err, "failed to get trip summary")
	}
	return summary, nil
}
