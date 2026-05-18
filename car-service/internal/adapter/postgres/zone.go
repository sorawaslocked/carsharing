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
)

type ZoneRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewZoneRepository(db *sql.DB, log *slog.Logger) *ZoneRepository {
	return &ZoneRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.ZoneRepo"),
	}
}

func (r *ZoneRepository) Insert(ctx context.Context, zone model.Zone) (string, error) {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO zones (name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at)
		VALUES (` + strings.Join([]string{
		b.Add(zone.Name),
		b.Add(string(zone.Type)),
		b.Add(zone.BoundaryGeoJSON),
		b.Add(zone.FeeAdjustment),
		b.Add(zone.IsActive),
		b.Add(zone.CreatedAt),
		b.Add(zone.UpdatedAt),
	}, ", ") + `) RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		logger.Error("failed to insert zone", pkglog.Err(err))
		return "", ErrSql
	}

	return id, nil
}

func (r *ZoneRepository) FindByID(ctx context.Context, id string) (model.Zone, error) {
	logger := pkglog.WithMethod(r.log, "FindByID")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `SELECT id, name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at
		FROM zones WHERE id = ` + b.Add(id) + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	zone, err := dto.ScanZoneRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Zone{}, ErrNotFound
		}
		logger.Error("failed to find zone by id", pkglog.Err(err))
		return model.Zone{}, ErrSql
	}

	return zone, nil
}

func (r *ZoneRepository) Find(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	clauses := dto.BuildZoneWhereClauses(filter, b)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at
		FROM zones` + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query zones", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.Zone
	for rows.Next() {
		zone, err := dto.ScanZoneRow(rows)
		if err != nil {
			logger.Error("failed to scan zone row", pkglog.Err(err))
			return nil, ErrSql
		}
		result = append(result, zone)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}

func (r *ZoneRepository) Update(ctx context.Context, id string, update model.ZoneUpdate) error {
	logger := pkglog.WithMethod(r.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildZoneSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	q := `UPDATE zones SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to update zone", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *ZoneRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(r.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `DELETE FROM zones WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to delete zone", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}
