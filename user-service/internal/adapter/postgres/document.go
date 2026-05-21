package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	pgdto "carsharing/user-service/internal/adapter/postgres/dto"
	"carsharing/user-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewDocumentRepository(log *slog.Logger, pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{
		log:  pkglog.WithComponent(log, "repo.DocumentRepository"),
		pool: pool,
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
	var errMsg *string
	var imageKey string

	if err := rs.Scan(
		&d.ID, &d.UserID, &imageType, &status, &errMsg, &imageKey,
		&d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return model.Document{}, err
	}

	d.ImageType = model.DocumentImageType(imageType)
	d.Status = model.DocumentStatus(status)
	d.Error = errMsg
	if imageKey != "" {
		d.Image = sharedmodel.Image{Key: imageKey}
	}

	return d, nil
}

func (r *DocumentRepository) Insert(ctx context.Context, doc model.Document) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	var imageKey *string
	if doc.Image.Key != "" {
		imageKey = &doc.Image.Key
	}

	var id string
	err := r.pool.QueryRow(ctx, `
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
		log.Error("unexpected postgres error", pkglog.Err(err))
		return "", model.ErrSql
	}

	return id, nil
}

func (r *DocumentRepository) FindByID(ctx context.Context, id string) (model.Document, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	doc, err := scanDocument(r.pool.QueryRow(ctx, documentSelect+" WHERE id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Document{}, model.ErrNotFound
		}
		log.Error("scanning document", pkglog.Err(err))
		return model.Document{}, model.ErrSql
	}

	return doc, nil
}

func (r *DocumentRepository) Find(ctx context.Context, filter model.DocumentFilter) ([]model.Document, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

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

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("querying documents", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		doc, err := scanDocument(rows)
		if err != nil {
			log.Error("scanning document row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		log.Error("iterating document rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return docs, nil
}

func (r *DocumentRepository) Update(ctx context.Context, id string, update model.DocumentUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, nextArg := pgdto.SetClausesFromDocumentUpdate(update)

	query := "UPDATE documents SET " + strings.Join(setClauses, ", ") +
		fmt.Sprintf(" WHERE id = $%d", nextArg)
	args = append(args, id)

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Error("unexpected postgres error", pkglog.Err(err))
		return model.ErrSql
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
