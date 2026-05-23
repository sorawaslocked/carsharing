package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"carsharing/car-service/internal/adapter/postgres/dto"
	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// --- CarMaintenanceTemplateRepository ---

type CarMaintenanceTemplateRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarMaintenanceTemplateRepository(log *slog.Logger, pool *pgxpool.Pool) *CarMaintenanceTemplateRepository {
	return &CarMaintenanceTemplateRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarMaintenanceTemplateRepository"),
		pool: pool,
	}
}

func (r *CarMaintenanceTemplateRepository) Insert(ctx context.Context, t model.CarMaintenanceTemplate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	args := []any{
		t.Name, t.KmInterval, t.DayInterval, t.IsMandatory,
		t.WarnPct, t.PullPct, t.CreatedAt, t.UpdatedAt,
	}
	q := `INSERT INTO car_maintenance_templates
			(name, km_interval, day_interval, is_mandatory, warn_pct, pull_pct, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		  RETURNING id`

	var id string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&id); err != nil {
		log.Error("failed to insert maintenance template", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *CarMaintenanceTemplateRepository) FindByID(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, name, km_interval, day_interval, is_mandatory, warn_pct, pull_pct, created_at, updated_at
		FROM car_maintenance_templates WHERE id = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, q, id)

	tmpl, err := dto.ScanMaintenanceTemplateRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.CarMaintenanceTemplate{}, model.ErrNotFound
		}
		log.Error("failed to find maintenance template by id", pkglog.Err(err))
		return model.CarMaintenanceTemplate{}, model.ErrSql
	}

	return tmpl, nil
}

func (r *CarMaintenanceTemplateRepository) Find(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	whereClauses, args, n := dto.WhereClausesFromMaintenanceTemplateFilter(filter, make([]any, 0), 0)
	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	q := `SELECT id, name, km_interval, day_interval, is_mandatory, warn_pct, pull_pct, created_at, updated_at
		FROM car_maintenance_templates` + where

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
		log.Error("failed to query maintenance templates", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.CarMaintenanceTemplate
	for rows.Next() {
		tmpl, err := dto.ScanMaintenanceTemplateRow(rows)
		if err != nil {
			log.Error("failed to scan maintenance template row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		result = append(result, tmpl)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}

func (r *CarMaintenanceTemplateRepository) Update(ctx context.Context, id string, update model.CarMaintenanceTemplateUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, n := dto.SetClausesFromMaintenanceTemplateUpdate(update)
	if len(setClauses) <= 1 {
		return nil
	}

	n++
	args = append(args, id)
	q := `UPDATE car_maintenance_templates SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", n)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		log.Error("failed to update maintenance template", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *CarMaintenanceTemplateRepository) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM car_maintenance_templates WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete maintenance template", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

// --- CarMaintenanceRecordRepository ---

type CarMaintenanceRecordRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarMaintenanceRecordRepository(log *slog.Logger, pool *pgxpool.Pool) *CarMaintenanceRecordRepository {
	return &CarMaintenanceRecordRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarMaintenanceRecordRepository"),
		pool: pool,
	}
}

func (r *CarMaintenanceRecordRepository) Insert(ctx context.Context, rec model.CarMaintenanceRecord) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	args := []any{
		rec.CarID, rec.TemplateID, string(rec.Status),
		rec.OdometerAt, rec.CompletedKM, rec.CostTenge,
		rec.AssignedTo, rec.DueBy, rec.CompletedAt, rec.Notes,
		dto.ImagesToKeys(rec.ReceiptImages), rec.CreatedAt, rec.UpdatedAt,
	}
	q := `INSERT INTO car_maintenance_records
			(car_id, template_id, status, odometer_at, completed_km, cost_tenge,
			 assigned_to, due_by, completed_at, notes, receipt_image_keys, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		  RETURNING id`

	var id string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&id); err != nil {
		log.Error("failed to insert maintenance record", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *CarMaintenanceRecordRepository) FindByID(ctx context.Context, id string) (model.CarMaintenanceRecord, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, car_id, template_id, status, odometer_at, completed_km, cost_tenge,
		assigned_to, due_by, completed_at, notes, receipt_image_keys, created_at, updated_at
		FROM car_maintenance_records WHERE id = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, q, id)

	rec, err := dto.ScanMaintenanceRecordRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.CarMaintenanceRecord{}, model.ErrNotFound
		}
		log.Error("failed to find maintenance record by id", pkglog.Err(err))
		return model.CarMaintenanceRecord{}, model.ErrSql
	}

	return rec, nil
}

func (r *CarMaintenanceRecordRepository) Find(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	whereClauses, args, n := dto.WhereClausesFromMaintenanceRecordFilter(filter, make([]any, 0), 0)
	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	q := `SELECT id, car_id, template_id, status, odometer_at, completed_km, cost_tenge,
		assigned_to, due_by, completed_at, notes, receipt_image_keys, created_at, updated_at
		FROM car_maintenance_records` + where

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
		log.Error("failed to query maintenance records", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.CarMaintenanceRecord
	for rows.Next() {
		rec, err := dto.ScanMaintenanceRecordRow(rows)
		if err != nil {
			log.Error("failed to scan maintenance record row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		result = append(result, rec)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}

func (r *CarMaintenanceRecordRepository) Update(ctx context.Context, id string, update model.CarMaintenanceRecordUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, n := dto.SetClausesFromMaintenanceRecordUpdate(update)
	if len(setClauses) <= 1 {
		return nil
	}

	n++
	args = append(args, id)
	q := `UPDATE car_maintenance_records SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", n)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		log.Error("failed to update maintenance record", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

// --- CarServiceStateRepository ---

type CarServiceStateRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewCarServiceStateRepository(log *slog.Logger, pool *pgxpool.Pool) *CarServiceStateRepository {
	return &CarServiceStateRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.CarServiceStateRepository"),
		pool: pool,
	}
}

func (r *CarServiceStateRepository) Upsert(ctx context.Context, state model.CarServiceState) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Upsert"), utils.MetadataFromCtx(ctx))

	args := []any{
		state.CarID, state.TemplateID, state.LastKM,
		state.LastDate, state.NextDueKM, state.NextDueDate,
	}
	q := `INSERT INTO car_service_states (car_id, template_id, last_km, last_date, next_due_km, next_due_date)
		  VALUES ($1, $2, $3, $4, $5, $6)
		  ON CONFLICT (car_id, template_id) DO UPDATE SET
			last_km       = EXCLUDED.last_km,
			last_date     = EXCLUDED.last_date,
			next_due_km   = EXCLUDED.next_due_km,
			next_due_date = EXCLUDED.next_due_date`

	if _, err := r.pool.Exec(ctx, q, args...); err != nil {
		log.Error("failed to upsert car service state", pkglog.Err(err))
		return model.ErrSql
	}

	return nil
}

func (r *CarServiceStateRepository) FindAll(ctx context.Context, filter model.CarServiceStateFilter) ([]model.CarServiceState, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindAll"), utils.MetadataFromCtx(ctx))

	var clauses []string
	var args []any
	n := 0

	if filter.CarID != nil {
		n++
		args = append(args, *filter.CarID)
		clauses = append(clauses, fmt.Sprintf("car_id = $%d", n))
	}
	if filter.TemplateID != nil {
		n++
		args = append(args, *filter.TemplateID)
		clauses = append(clauses, fmt.Sprintf("template_id = $%d", n))
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT car_id, template_id, last_km, last_date, next_due_km, next_due_date
		FROM car_service_states` + where

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		log.Error("failed to query car service states", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.CarServiceState
	for rows.Next() {
		var s model.CarServiceState
		var lastDate *time.Time
		var nextDueKM *int32
		var nextDueDate *time.Time

		if err = rows.Scan(
			&s.CarID, &s.TemplateID, &s.LastKM,
			&lastDate, &nextDueKM, &nextDueDate,
		); err != nil {
			log.Error("failed to scan service state row", pkglog.Err(err))
			return nil, model.ErrSql
		}

		s.LastDate = lastDate
		s.NextDueKM = nextDueKM
		s.NextDueDate = nextDueDate

		result = append(result, s)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}
