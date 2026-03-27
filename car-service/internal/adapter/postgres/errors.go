package postgres

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrSql           = errors.New("sql error")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
