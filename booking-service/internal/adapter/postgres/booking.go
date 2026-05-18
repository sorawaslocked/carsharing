package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-booking-service/internal/pkg/utils"
)

type pricingSnapshotJSON struct {
	RateTenge         int32   `json:"rate_tenge"`
	RatePerKMTenge    *int32  `json:"rate_per_km_tenge,omitempty"`
	FreeMinutes       *int32  `json:"free_minutes,omitempty"`
	MinChargeTenge    *int32  `json:"min_charge_tenge,omitempty"`
	OvertimePolicy    *string `json:"overtime_policy,omitempty"`
	OvertimeRateTenge *int32  `json:"overtime_rate_tenge,omitempty"`
}

type BookingRepo struct {
	log *slog.Logger
	db  *sql.DB
}

func NewBookingRepo(log *slog.Logger, db *sql.DB) *BookingRepo {
	return &BookingRepo{
		log: pkglog.WithComponent(log, "repo.BookingRepo"),
		db:  db,
	}
}

func (r *BookingRepo) Create(ctx context.Context, data model.BookingCreate, expiresAt time.Time) (string, error) {
	log := pkglog.WithMethod(r.log, "Create")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	var snap pricingSnapshotJSON
	err := r.db.QueryRowContext(ctx, `
		SELECT rate_tenge, rate_per_km_tenge, free_minutes, min_charge_tenge, overtime_policy, overtime_rate_tenge
		FROM pricing_rules WHERE id = $1 AND is_active = TRUE
	`, data.PricingRuleID).Scan(
		&snap.RateTenge, &snap.RatePerKMTenge, &snap.FreeMinutes,
		&snap.MinChargeTenge, &snap.OvertimePolicy, &snap.OvertimeRateTenge,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return "", model.ErrNotFound
	}
	if err != nil {
		log.Error("failed to fetch pricing rule for snapshot", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	snapJSON, err := json.Marshal(snap)
	if err != nil {
		log.Error("failed to marshal pricing snapshot", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}
	defer tx.Rollback()

	var id string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO bookings (user_id, car_id, committed_periods, pricing_rule_id, pricing_snapshot, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, data.UserID, data.CarID, data.CommittedPeriods, data.PricingRuleID, snapJSON, expiresAt).Scan(&id)
	if err != nil {
		log.Error("failed to insert booking", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO booking_status_history (booking_id, from_status, to_status, actor_type)
		VALUES ($1, '', 'created', 'system')
	`, id)
	if err != nil {
		log.Error("failed to insert initial status history", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction", pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	return id, nil
}

func (r *BookingRepo) GetByID(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMethod(r.log, "GetByID")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	var b model.Booking
	var statusStr string
	var snapJSON []byte
	var cp sql.NullInt32

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, car_id, committed_periods, status, pricing_rule_id, pricing_snapshot, expires_at, created_at, updated_at
		FROM bookings WHERE id = $1
	`, id).Scan(
		&b.ID, &b.UserID, &b.CarID, &cp,
		&statusStr, &b.PricingRuleID, &snapJSON, &b.ExpiresAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Booking{}, model.ErrNotFound
	}
	if err != nil {
		log.Error("failed to query booking", pkglog.Err(err))
		return model.Booking{}, model.ErrInternalServerError
	}

	b.Status = model.BookingStatus(statusStr)
	if cp.Valid {
		v := cp.Int32
		b.CommittedPeriods = &v
	}

	var snap pricingSnapshotJSON
	if err := json.Unmarshal(snapJSON, &snap); err != nil {
		log.Error("failed to unmarshal pricing snapshot", pkglog.Err(err))
		return model.Booking{}, model.ErrInternalServerError
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

func (r *BookingRepo) List(ctx context.Context, filter model.BookingListFilter) ([]model.Booking, error) {
	log := pkglog.WithMethod(r.log, "List")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if filter.UserID != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", idx))
		args = append(args, *filter.UserID)
		idx++
	}
	if filter.CarID != nil {
		where = append(where, fmt.Sprintf("car_id = $%d", idx))
		args = append(args, *filter.CarID)
		idx++
	}
	if filter.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, *filter.Status)
		idx++
	}
	if filter.PricingRuleID != nil {
		where = append(where, fmt.Sprintf("pricing_rule_id = $%d", idx))
		args = append(args, *filter.PricingRuleID)
		idx++
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, car_id, committed_periods, status, pricing_rule_id, pricing_snapshot, expires_at, created_at, updated_at
		FROM bookings WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), idx, idx+1)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to list bookings", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		var statusStr string
		var snapJSON []byte
		var cp sql.NullInt32

		if err := rows.Scan(
			&b.ID, &b.UserID, &b.CarID, &cp,
			&statusStr, &b.PricingRuleID, &snapJSON, &b.ExpiresAt, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			log.Error("failed to scan booking row", pkglog.Err(err))
			return nil, model.ErrInternalServerError
		}

		b.Status = model.BookingStatus(statusStr)
		if cp.Valid {
			v := cp.Int32
			b.CommittedPeriods = &v
		}

		var snap pricingSnapshotJSON
		if err := json.Unmarshal(snapJSON, &snap); err != nil {
			log.Error("failed to unmarshal pricing snapshot", pkglog.Err(err))
			return nil, model.ErrInternalServerError
		}
		b.PricingSnapshot = model.PricingSnapshot{
			RateTenge:         snap.RateTenge,
			RatePerKMTenge:    snap.RatePerKMTenge,
			FreeMinutes:       snap.FreeMinutes,
			MinChargeTenge:    snap.MinChargeTenge,
			OvertimePolicy:    snap.OvertimePolicy,
			OvertimeRateTenge: snap.OvertimeRateTenge,
		}
		bookings = append(bookings, b)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}

	return bookings, nil
}

func (r *BookingRepo) ListCreatedExpired(ctx context.Context, now time.Time) ([]model.Booking, error) {
	log := pkglog.WithMethod(r.log, "ListCreatedExpired")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, car_id, committed_periods, status, pricing_rule_id, pricing_snapshot, expires_at, created_at, updated_at
		FROM bookings WHERE status = 'created' AND expires_at <= $1
	`, now)
	if err != nil {
		log.Error("failed to query expired bookings", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		var statusStr string
		var snapJSON []byte
		var cp sql.NullInt32

		if err := rows.Scan(
			&b.ID, &b.UserID, &b.CarID, &cp,
			&statusStr, &b.PricingRuleID, &snapJSON, &b.ExpiresAt, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			log.Error("failed to scan expired booking row", pkglog.Err(err))
			return nil, model.ErrInternalServerError
		}

		b.Status = model.BookingStatus(statusStr)
		if cp.Valid {
			v := cp.Int32
			b.CommittedPeriods = &v
		}

		var snap pricingSnapshotJSON
		if err := json.Unmarshal(snapJSON, &snap); err != nil {
			log.Error("failed to unmarshal pricing snapshot", pkglog.Err(err))
			return nil, model.ErrInternalServerError
		}
		b.PricingSnapshot = model.PricingSnapshot{
			RateTenge:         snap.RateTenge,
			RatePerKMTenge:    snap.RatePerKMTenge,
			FreeMinutes:       snap.FreeMinutes,
			MinChargeTenge:    snap.MinChargeTenge,
			OvertimePolicy:    snap.OvertimePolicy,
			OvertimeRateTenge: snap.OvertimeRateTenge,
		}
		bookings = append(bookings, b)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}

	return bookings, nil
}

func (r *BookingRepo) UpdateStatus(ctx context.Context, id, toStatus, actorType string, actorID, reason *string) error {
	log := pkglog.WithMethod(r.log, "UpdateStatus")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", pkglog.Err(err))
		return model.ErrInternalServerError
	}
	defer tx.Rollback()

	var fromStatus string
	err = tx.QueryRowContext(ctx, `SELECT status FROM bookings WHERE id = $1 FOR UPDATE`, id).Scan(&fromStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		log.Error("failed to lock booking row", pkglog.Err(err))
		return model.ErrInternalServerError
	}

	if _, err = tx.ExecContext(ctx, `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`, toStatus, id); err != nil {
		log.Error("failed to update booking status", pkglog.Err(err))
		return model.ErrInternalServerError
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO booking_status_history (booking_id, from_status, to_status, actor_type, actor_id, reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, fromStatus, toStatus, actorType, actorID, reason); err != nil {
		log.Error("failed to insert status history entry", pkglog.Err(err))
		return model.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction", pkglog.Err(err))
		return model.ErrInternalServerError
	}

	return nil
}

func (r *BookingRepo) GetStatusHistory(ctx context.Context, filter model.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error) {
	log := pkglog.WithMethod(r.log, "GetStatusHistory")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	where := []string{"booking_id = $1"}
	args := []any{filter.BookingID}
	idx := 2

	if filter.From != nil {
		where = append(where, fmt.Sprintf("changed_at >= $%d", idx))
		args = append(args, *filter.From)
		idx++
	}
	if filter.To != nil {
		where = append(where, fmt.Sprintf("changed_at <= $%d", idx))
		args = append(args, *filter.To)
		idx++
	}

	query := fmt.Sprintf(`
		SELECT id, booking_id, from_status, to_status, actor_type, actor_id, reason, changed_at
		FROM booking_status_history WHERE %s ORDER BY changed_at ASC LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), idx, idx+1)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to query status history", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}
	defer rows.Close()

	var history []model.BookingStatusReading
	for rows.Next() {
		var reading model.BookingStatusReading
		var actorID, reason sql.NullString

		if err := rows.Scan(
			&reading.ID, &reading.BookingID, &reading.FromStatus, &reading.ToStatus,
			&reading.ActorType, &actorID, &reason, &reading.ChangedAt,
		); err != nil {
			log.Error("failed to scan status history row", pkglog.Err(err))
			return nil, model.ErrInternalServerError
		}

		if actorID.Valid {
			reading.ActorID = &actorID.String
		}
		if reason.Valid {
			reading.Reason = &reason.String
		}

		history = append(history, reading)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrInternalServerError
	}

	return history, nil
}
