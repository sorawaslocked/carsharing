package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
