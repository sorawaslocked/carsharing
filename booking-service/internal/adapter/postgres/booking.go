package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pgdto "carsharing/booking-service/internal/adapter/postgres/dto"
	"carsharing/booking-service/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pricingSnapshotJSON struct {
	RateTenge         int32   `json:"rate_tenge"`
	RatePerKMTenge    *int32  `json:"rate_per_km_tenge,omitempty"`
	FreeMinutes       *int32  `json:"free_minutes,omitempty"`
	MinChargeTenge    *int32  `json:"min_charge_tenge,omitempty"`
	OvertimePolicy    *string `json:"overtime_policy,omitempty"`
	OvertimeRateTenge *int32  `json:"overtime_rate_tenge,omitempty"`
}

type rowScanner interface {
	Scan(dest ...any) error
}

type BookingRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewBookingRepository(log *slog.Logger, pool *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.BookingRepository"),
		pool: pool,
	}
}

const bookingSelect = `
    SELECT id, user_id, car_id, committed_periods, status, pricing_rule_id, pricing_snapshot, expires_at, created_at, updated_at
    FROM bookings`

func scanBooking(rs rowScanner) (model.Booking, error) {
	var b model.Booking
	var statusStr string
	var snapJSON []byte
	var cp *int32

	if err := rs.Scan(
		&b.ID, &b.UserID, &b.CarID, &cp,
		&statusStr, &b.PricingRuleID, &snapJSON, &b.ExpiresAt, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return model.Booking{}, err
	}

	b.Status = model.BookingStatus(statusStr)
	b.CommittedPeriods = cp

	var snap pricingSnapshotJSON
	if err := json.Unmarshal(snapJSON, &snap); err != nil {
		return model.Booking{}, err
	}
	b.PricingSnapshot = model.PricingSnapshot{
		RateTenge:         snap.RateTenge,
		RatePerKMTenge:    snap.RatePerKMTenge,
		FreeMinutes:       snap.FreeMinutes,
		MinChargeTenge:    snap.MinChargeTenge,
		OvertimePolicy:    snap.OvertimePolicy,
		OvertimeRateTenge: snap.OvertimeRateTenge,
	}

	return b, nil
}

func (r *BookingRepository) Create(ctx context.Context, data model.BookingCreate, expiresAt time.Time) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Create"), utils.MetadataFromCtx(ctx))

	var snap pricingSnapshotJSON
	var ratePerKM, freeMinutes, minCharge, overtimeRate *int32
	var overtimePolicy *string

	err := r.pool.QueryRow(ctx, `
        SELECT rate_tenge, rate_per_km_tenge, free_minutes, min_charge_tenge, overtime_policy, overtime_rate_tenge
        FROM pricing_rules WHERE id = $1 AND is_active = TRUE
    `, data.PricingRuleID).Scan(
		&snap.RateTenge, &ratePerKM, &freeMinutes,
		&minCharge, &overtimePolicy, &overtimeRate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.ErrPricingRuleNotFound
	}
	if err != nil {
		log.Error("fetching pricing rule for snapshot", pkglog.Err(err))
		return "", model.ErrSql
	}
	snap.RatePerKMTenge = ratePerKM
	snap.FreeMinutes = freeMinutes
	snap.MinChargeTenge = minCharge
	snap.OvertimePolicy = overtimePolicy
	snap.OvertimeRateTenge = overtimeRate

	snapJSON, err := json.Marshal(snap)
	if err != nil {
		log.Error("marshaling pricing snapshot", pkglog.Err(err))
		return "", model.ErrSql
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		log.Error("beginning transaction", pkglog.Err(err))
		return "", model.ErrSqlTransaction
	}
	defer tx.Rollback(ctx)

	var id string
	err = tx.QueryRow(ctx, `
        INSERT INTO bookings (user_id, car_id, committed_periods, pricing_rule_id, pricing_snapshot, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, data.UserID, data.CarID, data.CommittedPeriods, data.PricingRuleID, snapJSON, expiresAt).Scan(&id)
	if err != nil {
		log.Error("inserting booking", pkglog.Err(err))
		return "", model.ErrSql
	}

	if _, err = tx.Exec(ctx, `
        INSERT INTO booking_status_history (booking_id, from_status, to_status, actor_type)
        VALUES ($1, '', 'created', 'system')
    `, id); err != nil {
		log.Error("inserting initial status history", pkglog.Err(err))
		return "", model.ErrSql
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("committing transaction", pkglog.Err(err))
		return "", model.ErrSqlTransaction
	}

	return id, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "GetByID"), utils.MetadataFromCtx(ctx))

	booking, err := scanBooking(r.pool.QueryRow(ctx, bookingSelect+" WHERE id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Booking{}, model.ErrBookingNotFound
		}
		log.Error("scanning booking", pkglog.Err(err))
		return model.Booking{}, model.ErrSql
	}
	return booking, nil
}

func (r *BookingRepository) List(ctx context.Context, filter model.BookingListFilter) ([]model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "List"), utils.MetadataFromCtx(ctx))

	clauses, args, nextArg := pgdto.WhereClausesFromBookingListFilter(filter, nil, 1)

	query := bookingSelect
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", nextArg, nextArg+1)
	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("querying bookings", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		booking, err := scanBooking(rows)
		if err != nil {
			log.Error("scanning booking row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		log.Error("iterating booking rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return bookings, nil
}

func (r *BookingRepository) ListCreatedExpired(ctx context.Context, now time.Time) ([]model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "ListCreatedExpired"), utils.MetadataFromCtx(ctx))

	rows, err := r.pool.Query(ctx, bookingSelect+" WHERE status = 'created' AND expires_at <= $1", now)
	if err != nil {
		log.Error("querying expired bookings", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		booking, err := scanBooking(rows)
		if err != nil {
			log.Error("scanning expired booking row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		log.Error("iterating expired booking rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return bookings, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id string, status model.BookingStatus, actorType sharedmodel.ActorType, actorID *string, reason *string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "UpdateStatus"), utils.MetadataFromCtx(ctx))

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		log.Error("beginning transaction", pkglog.Err(err))
		return model.ErrSqlTransaction
	}
	defer tx.Rollback(ctx)

	var fromStatus string
	err = tx.QueryRow(ctx, `SELECT status FROM bookings WHERE id = $1 FOR UPDATE`, id).Scan(&fromStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrBookingNotFound
	}
	if err != nil {
		log.Error("locking booking row", pkglog.Err(err))
		return model.ErrSql
	}

	if _, err = tx.Exec(ctx, `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`, status, id); err != nil {
		log.Error("updating booking status", pkglog.Err(err))
		return model.ErrSql
	}

	if _, err = tx.Exec(ctx, `
        INSERT INTO booking_status_history (booking_id, from_status, to_status, actor_type, actor_id, reason)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, id, fromStatus, status, actorType, actorID, reason); err != nil {
		log.Error("inserting status history entry", pkglog.Err(err))
		return model.ErrSql
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("committing transaction", pkglog.Err(err))
		return model.ErrSqlTransaction
	}

	return nil
}

func (r *BookingRepository) GetStatusHistory(ctx context.Context, filter model.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "GetStatusHistory"), utils.MetadataFromCtx(ctx))

	clauses, args, nextArg := pgdto.WhereClausesFromStatusHistoryFilter(filter, nil, 1)

	query := fmt.Sprintf(`
        SELECT id, booking_id, from_status, to_status, actor_type, actor_id, reason, changed_at
        FROM booking_status_history WHERE %s ORDER BY changed_at ASC LIMIT $%d OFFSET $%d
    `, strings.Join(clauses, " AND "), nextArg, nextArg+1)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("querying status history", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var history []model.BookingStatusReading
	for rows.Next() {
		var reading model.BookingStatusReading
		var actorType string

		if err := rows.Scan(
			&reading.ID, &reading.BookingID, &reading.FromStatus, &reading.ToStatus,
			&actorType, &reading.ActorID, &reading.Reason, &reading.ChangedAt,
		); err != nil {
			log.Error("scanning status history row", pkglog.Err(err))
			return nil, model.ErrSql
		}

		reading.ActorType = sharedmodel.ActorType(actorType)
		history = append(history, reading)
	}

	if err := rows.Err(); err != nil {
		log.Error("iterating status history rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return history, nil
}
