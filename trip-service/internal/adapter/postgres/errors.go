package postgres

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/sorawaslocked/car-rental-trip-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
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
