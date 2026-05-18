package postgres

import (
	"database/sql"
	"errors"
	"log/slog"

	"carsharing/trip-service/internal/model"
	pkglog "carsharing/trip-service/internal/pkg/log"
	"github.com/lib/pq"
)

func mapSQLError(log *slog.Logger, err error, msg string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if isUniqueViolation(err) {
		return model.ErrAlreadyExists
	}
	log.Error(msg, pkglog.Err(err))
	return model.ErrSQL
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
