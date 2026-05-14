package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `yaml:"host"     env:"PG_HOST"     env-required:"true"`
	Port     int    `yaml:"port"     env:"PG_PORT"     env-default:"5432"`
	User     string `yaml:"user"     env:"PG_USER"     env-required:"true"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	Database string `yaml:"database" env:"PG_DATABASE" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"PG_SSL_MODE" env-default:"disable"`
}

func New(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return db, nil
}
