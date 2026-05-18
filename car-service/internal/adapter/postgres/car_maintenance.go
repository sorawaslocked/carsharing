package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"carsharing/car-service/internal/adapter/postgres/dto"
	"carsharing/car-service/internal/model"
	pkglog "carsharing/car-service/internal/pkg/log"
	"carsharing/car-service/internal/pkg/utils"
	"github.com/lib/pq"
)

// --- CarMaintenanceTemplateRepository ---

type CarMaintenanceTemplateRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarMaintenanceTemplateRepository(db *sql.DB, log *slog.Logger) *CarMaintenanceTemplateRepository {
	return &CarMaintenanceTemplateRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarMaintenanceTemplateRepo"),
	}
}

func (r *CarMaintenanceTemplateRepository) Insert(ctx context.Context, t model.CarMaintenanceTemplate) (string, error) {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO car_maintenance_templates
			(name, km_interval, day_interval, is_mandatory, warn_pct, pull_pct, created_at, updated_at)
		VALUES (` + strings.Join([]string{
		b.Add(t.Name),
		b.Add(t.KmInterval),
		b.Add(t.DayInterval),
		b.Add(t.IsMandatory),
		b.Add(t.WarnPct),
		b.Add(t.PullPct),
		b.Add(t.CreatedAt),
		b.Add(t.UpdatedAt),
	}, ", ") + `) RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		logger.Error("failed to insert maintenance template", pkglog.Err(err))
		return "", ErrSql
	}

	return id, nil
}

func (r *CarMaintenanceTemplateRepository) FindByID(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	logger := pkglog.WithMethod(r.log, "FindByID")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `SELECT id, name, km_interval, day_interval, is_mandatory, warn_pct, pull_pct, created_at, updated_at
		FROM car_maintenance_templates WHERE id = ` + b.Add(id) + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	tmpl, err := dto.ScanMaintenanceTemplateRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CarMaintenanceTemplate{}, ErrNotFound
		}
		logger.Error("failed to find maintenance template by id", pkglog.Err(err))
		return model.CarMaintenanceTemplate{}, ErrSql
	}

	return tmpl, nil
}

func (r *CarMaintenanceTemplateRepository) Find(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	clauses := dto.BuildMaintenanceTemplateWhereClauses(filter, b)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, name, km_interval, day_interval, is_mandatory, warn_pct, pull_pct, created_at, updated_at
		FROM car_maintenance_templates` + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query maintenance templates", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarMaintenanceTemplate
	for rows.Next() {
		tmpl, err := dto.ScanMaintenanceTemplateRow(rows)
		if err != nil {
			logger.Error("failed to scan maintenance template row", pkglog.Err(err))
			return nil, ErrSql
		}
		result = append(result, tmpl)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarMaintenanceTemplateRepository) Update(ctx context.Context, id string, update model.CarMaintenanceTemplateUpdate) error {
	logger := pkglog.WithMethod(r.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildMaintenanceTemplateSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	q := `UPDATE car_maintenance_templates SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to update maintenance template", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CarMaintenanceTemplateRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(r.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `DELETE FROM car_maintenance_templates WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to delete maintenance template", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

// --- CarMaintenanceRecordRepository ---

type CarMaintenanceRecordRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarMaintenanceRecordRepository(db *sql.DB, log *slog.Logger) *CarMaintenanceRecordRepository {
	return &CarMaintenanceRecordRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarMaintenanceRecordRepo"),
	}
}

func (r *CarMaintenanceRecordRepository) Insert(ctx context.Context, rec model.CarMaintenanceRecord) (string, error) {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO car_maintenance_records
			(car_id, template_id, status, odometer_at, completed_km, cost_tenge,
			 assigned_to, due_by, completed_at, notes, receipt_image_keys, created_at, updated_at)
		VALUES (` + strings.Join([]string{
		b.Add(rec.CarID),
		b.Add(rec.TemplateID),
		b.Add(string(rec.Status)),
		b.Add(rec.OdometerAt),
		b.Add(rec.CompletedKM),
		b.Add(rec.CostTenge),
		b.Add(rec.AssignedTo),
		b.Add(rec.DueBy),
		b.Add(rec.CompletedAt),
		b.Add(rec.Notes),
		b.Add(pq.StringArray(dto.ImagesToKeys(rec.ReceiptImages))),
		b.Add(rec.CreatedAt),
		b.Add(rec.UpdatedAt),
	}, ", ") + `) RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		logger.Error("failed to insert maintenance record", pkglog.Err(err))
		return "", ErrSql
	}

	return id, nil
}

func (r *CarMaintenanceRecordRepository) FindByID(ctx context.Context, id string) (model.CarMaintenanceRecord, error) {
	logger := pkglog.WithMethod(r.log, "FindByID")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `SELECT id, car_id, template_id, status, odometer_at, completed_km, cost_tenge,
		assigned_to, due_by, completed_at, notes, receipt_image_keys, created_at, updated_at
		FROM car_maintenance_records WHERE id = ` + b.Add(id) + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	rec, err := dto.ScanMaintenanceRecordRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CarMaintenanceRecord{}, ErrNotFound
		}
		logger.Error("failed to find maintenance record by id", pkglog.Err(err))
		return model.CarMaintenanceRecord{}, ErrSql
	}

	return rec, nil
}

func (r *CarMaintenanceRecordRepository) Find(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	clauses := dto.BuildMaintenanceRecordWhereClauses(filter, b)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, car_id, template_id, status, odometer_at, completed_km, cost_tenge,
		assigned_to, due_by, completed_at, notes, receipt_image_keys, created_at, updated_at
		FROM car_maintenance_records` + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query maintenance records", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarMaintenanceRecord
	for rows.Next() {
		rec, err := dto.ScanMaintenanceRecordRow(rows)
		if err != nil {
			logger.Error("failed to scan maintenance record row", pkglog.Err(err))
			return nil, ErrSql
		}
		result = append(result, rec)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarMaintenanceRecordRepository) Update(ctx context.Context, id string, update model.CarMaintenanceRecordUpdate) error {
	logger := pkglog.WithMethod(r.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildMaintenanceRecordSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	q := `UPDATE car_maintenance_records SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to update maintenance record", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

// --- CarServiceStateRepository ---

type CarServiceStateRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarServiceStateRepository(db *sql.DB, log *slog.Logger) *CarServiceStateRepository {
	return &CarServiceStateRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarServiceStateRepo"),
	}
}

func (r *CarServiceStateRepository) Upsert(ctx context.Context, state model.CarServiceState) error {
	logger := pkglog.WithMethod(r.log, "Upsert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO car_service_states (car_id, template_id, last_km, last_date, next_due_km, next_due_date)
		VALUES (` + strings.Join([]string{
		b.Add(state.CarID),
		b.Add(state.TemplateID),
		b.Add(state.LastKM),
		b.Add(state.LastDate),
		b.Add(state.NextDueKM),
		b.Add(state.NextDueDate),
	}, ", ") + `)
		ON CONFLICT (car_id, template_id) DO UPDATE SET
			last_km      = EXCLUDED.last_km,
			last_date    = EXCLUDED.last_date,
			next_due_km  = EXCLUDED.next_due_km,
			next_due_date = EXCLUDED.next_due_date`

	_, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to upsert car service state", pkglog.Err(err))
		return ErrSql
	}

	return nil
}

func (r *CarServiceStateRepository) FindAll(ctx context.Context, filter model.CarServiceStateFilter) ([]model.CarServiceState, error) {
	logger := pkglog.WithMethod(r.log, "FindAll")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	var clauses []string
	if filter.CarID != nil {
		clauses = append(clauses, "car_id = "+b.Add(*filter.CarID))
	}
	if filter.TemplateID != nil {
		clauses = append(clauses, "template_id = "+b.Add(*filter.TemplateID))
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT car_id, template_id, last_km, last_date, next_due_km, next_due_date
		FROM car_service_states` + where

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query car service states", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarServiceState
	for rows.Next() {
		var s model.CarServiceState
		var lastDate sql.NullTime
		var nextDueKM sql.NullInt32
		var nextDueDate sql.NullTime

		if err = rows.Scan(
			&s.CarID, &s.TemplateID, &s.LastKM,
			&lastDate, &nextDueKM, &nextDueDate,
		); err != nil {
			logger.Error("failed to scan service state row", pkglog.Err(err))
			return nil, ErrSql
		}

		if lastDate.Valid {
			s.LastDate = &lastDate.Time
		}
		if nextDueKM.Valid {
			v := nextDueKM.Int32
			s.NextDueKM = &v
		}
		if nextDueDate.Valid {
			s.NextDueDate = &nextDueDate.Time
		}

		result = append(result, s)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}
