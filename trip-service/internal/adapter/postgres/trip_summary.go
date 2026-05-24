package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/postgres/dto"
	"carsharing/trip-service/internal/model"
)

type TripSummaryRepo struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewTripSummaryRepo(log *slog.Logger, pool *pgxpool.Pool) *TripSummaryRepo {
	return &TripSummaryRepo{
		log:  pkglog.WithComponent(log, "adapter.postgres.TripSummaryRepository"),
		pool: pool,
	}
}

func (r *TripSummaryRepo) Create(ctx context.Context, s model.TripSummaryCreate) (model.TripSummary, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Create"), pkgutils.MetadataFromCtx(ctx))

	snapJSON, err := dto.MarshalPricingSnapshot(s.PricingSnapshot)
	if err != nil {
		log.Error("marshaling pricing snapshot", pkglog.Err(err))
		return model.TripSummary{}, model.ErrSQL
	}

	q := fmt.Sprintf(`
		INSERT INTO trip_summaries (
			trip_id, booking_id, started_at, ended_at,
			duration_seconds, distance_traveled_km, pricing_snapshot,
			base_cost_tenge, distance_cost_tenge, overtime_cost_tenge, total_cost_tenge
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING %s`, dto.TripSummaryColumns)

	summary, err := dto.ScanTripSummary(r.pool.QueryRow(ctx, q,
		s.TripID, s.BookingID, s.StartedAt, s.EndedAt,
		s.DurationSeconds, s.DistanceTraveledKM, snapJSON,
		s.BaseCostTenge, s.DistanceCostTenge, s.OvertimeCostTenge, s.TotalCostTenge,
	))
	if err != nil {
		return model.TripSummary{}, mapSQLError(log, err, "creating trip summary")
	}
	return summary, nil
}

func (r *TripSummaryRepo) GetByTripID(ctx context.Context, tripID string) (model.TripSummary, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "GetByTripID"), pkgutils.MetadataFromCtx(ctx))

	q := fmt.Sprintf(`SELECT %s FROM trip_summaries WHERE trip_id = $1`, dto.TripSummaryColumns)

	summary, err := dto.ScanTripSummary(r.pool.QueryRow(ctx, q, tripID))
	if err != nil {
		return model.TripSummary{}, mapSQLError(log, err, "getting trip summary")
	}
	return summary, nil
}
