package postgres

import (
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/trip-service/internal/model"
)

func mapSQLError(log *slog.Logger, err error, msg string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrNotFound
	}
	if isUniqueViolation(err) {
		return model.ErrAlreadyExists
	}
	log.Error(msg, pkglog.Err(err))
	return model.ErrSQL
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
