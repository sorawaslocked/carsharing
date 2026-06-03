package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"carsharing/car-service/internal/adapter/postgres/dto"
	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ZoneRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewZoneRepository(log *slog.Logger, pool *pgxpool.Pool) *ZoneRepository {
	return &ZoneRepository{
		log:  pkglog.WithComponent(log, "adapter.postgres.ZoneRepository"),
		pool: pool,
	}
}

func (r *ZoneRepository) Insert(ctx context.Context, zone model.Zone) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	args := []any{
		zone.Name, string(zone.Type), zone.BoundaryGeoJSON,
		zone.FeeAdjustment, zone.IsActive, zone.CreatedAt, zone.UpdatedAt,
	}
	q := `INSERT INTO zones (name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7)
		  RETURNING id`

	var id string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&id); err != nil {
		log.Error("failed to insert zone", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *ZoneRepository) FindByID(ctx context.Context, id string) (model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at
		FROM zones WHERE id = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, q, id)

	zone, err := dto.ScanZoneRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Zone{}, model.ErrZoneNotFound
		}
		log.Error("failed to find zone by id", pkglog.Err(err))
		return model.Zone{}, model.ErrSql
	}

	return zone, nil
}

func (r *ZoneRepository) FindByLocation(ctx context.Context, lat, lng float64) (*model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByLocation"), utils.MetadataFromCtx(ctx))

	q := `SELECT id, name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at
		FROM zones
		WHERE is_active = true
		  AND ST_Contains(
		    ST_GeomFromGeoJSON(boundary_geo_json),
		    ST_SetSRID(ST_MakePoint($2, $1), 4326)
		  )
		ORDER BY CASE WHEN type = 'no_drop' THEN 0 ELSE 1 END
		LIMIT 1`

	row := r.pool.QueryRow(ctx, q, lat, lng)

	zone, err := dto.ScanZoneRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Error("failed to find zone by location", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return &zone, nil
}

func (r *ZoneRepository) Find(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	whereClauses, args, n := dto.WhereClausesFromZoneFilter(filter, make([]any, 0), 0)
	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	q := `SELECT id, name, type, boundary_geo_json, fee_adjustment, is_active, created_at, updated_at
		FROM zones` + where

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
		log.Error("failed to query zones", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var result []model.Zone
	for rows.Next() {
		zone, err := dto.ScanZoneRow(rows)
		if err != nil {
			log.Error("failed to scan zone row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		result = append(result, zone)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows iteration error", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return result, nil
}

func (r *ZoneRepository) Update(ctx context.Context, id string, update model.ZoneUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, n := dto.SetClausesFromZoneUpdate(update)
	if len(setClauses) <= 1 {
		return nil
	}

	n++
	args = append(args, id)
	q := `UPDATE zones SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", n)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		log.Error("failed to update zone", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrZoneNotFound
	}

	return nil
}

func (r *ZoneRepository) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM zones WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete zone", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrZoneNotFound
	}

	return nil
}
