package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	pgdto "github.com/sorawaslocked/car-rental-user-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/utils"
)

type DocumentRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewDocumentRepository(log *slog.Logger, db *sql.DB) *DocumentRepository {
	return &DocumentRepository{
		log: pkglog.WithComponent(log, "repo.DocumentRepository"),
		db:  db,
	}
}

const documentSelect = `
    SELECT id, user_id, image_type, status, error, image_key, created_at, updated_at
    FROM documents`

// DISTINCT ON keeps the first row per image_type; ORDER BY image_type, created_at DESC
// makes the most recent document the first for each type.
const documentSelectLatestPerType = `
    SELECT DISTINCT ON (image_type) id, user_id, image_type, status, error, image_key, created_at, updated_at
    FROM documents`

func scanDocument(rs rowScanner) (model.Document, error) {
	var d model.Document
	var imageType string
	var status string
	var errMsg sql.NullString
	var imageKey string

	if err := rs.Scan(
		&d.ID, &d.UserID, &imageType, &status, &errMsg, &imageKey,
		&d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return model.Document{}, err
	}

	d.ImageType = model.ImageType(imageType)
	d.Status = model.DocumentStatus(status)

	if errMsg.Valid {
		d.Error = &errMsg.String
	}
	if imageKey != "" {
		d.Image = &model.Image{Key: imageKey}
	}

	return d, nil
}

func (r *DocumentRepository) Insert(ctx context.Context, doc model.Document) (string, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	var imageKey *string
	if doc.Image != nil {
		imageKey = &doc.Image.Key
	}

	var id string
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO documents (user_id, image_type, status, error, image_key, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		doc.UserID,
		doc.ImageType.String(),
		doc.Status.String(),
		doc.Error,
		imageKey,
		doc.CreatedAt,
		doc.UpdatedAt,
	).Scan(&id)
	if err != nil {
		logger.Error("unexpected sql error", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *DocumentRepository) FindByID(ctx context.Context, id string) (model.Document, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	doc, err := scanDocument(r.db.QueryRowContext(ctx, documentSelect+" WHERE id = $1", id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Document{}, model.ErrNotFound
		}
		logger.Error("scanning document", pkglog.Err(err))
		return model.Document{}, model.ErrSql
	}

	return doc, nil
}

func (r *DocumentRepository) Find(ctx context.Context, filter model.DocumentFilter) ([]model.Document, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	base := documentSelect
	if filter.LatestPerType {
		base = documentSelectLatestPerType
	}

	query := base
	clauses, args, _ := pgdto.WhereClausesFromDocumentFilter(filter, nil, 1)
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	if filter.LatestPerType {
		query += " ORDER BY image_type, created_at DESC"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("querying documents", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		doc, err := scanDocument(rows)
		if err != nil {
			logger.Error("scanning document row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		logger.Error("iterating document rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return docs, nil
}

func (r *DocumentRepository) Update(ctx context.Context, id string, update model.DocumentUpdate) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, nextArg := pgdto.SetClausesFromDocumentUpdate(update)

	query := "UPDATE documents SET " + strings.Join(setClauses, ", ") +
		fmt.Sprintf(" WHERE id = $%d", nextArg)
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error("unexpected sql error", pkglog.Err(err))
		return model.ErrSql
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}
