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

type CarInsuranceRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCarInsuranceRepository(db *sql.DB, log *slog.Logger) *CarInsuranceRepository {
	return &CarInsuranceRepository{
		db:  db,
		log: pkglog.WithComponent(log, "repo.CarInsuranceRepo"),
	}
}

func (r *CarInsuranceRepository) Insert(ctx context.Context, ins model.CarInsurance) (string, error) {
	logger := pkglog.WithMethod(r.log, "Insert")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `
		INSERT INTO car_insurances
			(car_id, type, provider, policy_num, starts_at, expires_at,
			 cost_tenge, status, notes, image_keys, created_at, updated_at)
		VALUES (` + strings.Join([]string{
		b.Add(ins.CarID),
		b.Add(string(ins.Type)),
		b.Add(ins.Provider),
		b.Add(ins.PolicyNum),
		b.Add(ins.StartsAt),
		b.Add(ins.ExpiresAt),
		b.Add(ins.CostTenge),
		b.Add(string(ins.Status)),
		b.Add(ins.Notes),
		b.Add(pq.StringArray(dto.ImagesToKeys(ins.Images))),
		b.Add(ins.CreatedAt),
		b.Add(ins.UpdatedAt),
	}, ", ") + `) RETURNING id`

	var id string

	err := r.db.QueryRowContext(ctx, q, b.Args...).Scan(&id)
	if err != nil {
		logger.Error("failed to insert car insurance", pkglog.Err(err))
		return "", ErrSql
	}

	return id, nil
}

func (r *CarInsuranceRepository) FindByID(ctx context.Context, id string) (model.CarInsurance, error) {
	logger := pkglog.WithMethod(r.log, "FindByID")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `SELECT id, car_id, type, provider, policy_num, starts_at, expires_at,
		cost_tenge, status, notes, image_keys, created_at, updated_at
		FROM car_insurances WHERE id = ` + b.Add(id) + ` LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, b.Args...)

	ins, err := dto.ScanCarInsuranceRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CarInsurance{}, ErrNotFound
		}
		logger.Error("failed to find car insurance by id", pkglog.Err(err))
		return model.CarInsurance{}, ErrSql
	}

	return ins, nil
}

func (r *CarInsuranceRepository) Find(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	logger := pkglog.WithMethod(r.log, "Find")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	clauses := dto.BuildCarInsuranceWhereClauses(filter, b)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	q := `SELECT id, car_id, type, provider, policy_num, starts_at, expires_at,
		cost_tenge, status, notes, image_keys, created_at, updated_at
		FROM car_insurances` + where + dto.BuildPagination(b, filter.Pagination)

	rows, err := r.db.QueryContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to query car insurances", pkglog.Err(err))
		return nil, ErrSql
	}
	defer rows.Close()

	var result []model.CarInsurance
	for rows.Next() {
		ins, err := dto.ScanCarInsuranceRow(rows)
		if err != nil {
			logger.Error("failed to scan car insurance row", pkglog.Err(err))
			return nil, ErrSql
		}
		result = append(result, ins)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", pkglog.Err(err))
		return nil, ErrSql
	}

	return result, nil
}

func (r *CarInsuranceRepository) Update(ctx context.Context, id string, update model.CarInsuranceUpdate) error {
	logger := pkglog.WithMethod(r.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	setClauses := dto.BuildCarInsuranceSetClauses(update, b)
	if len(setClauses) <= 1 {
		return nil
	}

	q := `UPDATE car_insurances SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to update car insurance", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CarInsuranceRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(r.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	b := &dto.ArgsBuilder{}

	q := `DELETE FROM car_insurances WHERE id = ` + b.Add(id)

	res, err := r.db.ExecContext(ctx, q, b.Args...)
	if err != nil {
		logger.Error("failed to delete car insurance", pkglog.Err(err))
		return ErrSql
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}
