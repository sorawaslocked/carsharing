package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Host               string        `yaml:"host"                    env:"POSTGRES_HOST"                 env-required:"true"`
	Port               int           `yaml:"port"                    env:"POSTGRES_PORT"                 env-default:"5432"`
	User               string        `yaml:"user"                    env:"POSTGRES_USER"                 env-required:"true"`
	Password           string        `yaml:"password"                env:"POSTGRES_PASSWORD"             env-required:"true"`
	Database           string        `yaml:"database"                env:"POSTGRES_DATABASE"             env-required:"true"`
	SSLMode            string        `yaml:"ssl_mode"                env:"POSTGRES_SSL_MODE"             env-default:"disable"`
	MaxOpenConnections int           `yaml:"max_open_connections"    env:"POSTGRES_MAX_OPEN_CONNECTIONS" env-default:"25"`
	MaxIdleConnections int           `yaml:"max_idle_connections"    env:"POSTGRES_MAX_IDLE_CONNECTIONS" env-default:"25"`
	MaxIdleTime        time.Duration `yaml:"max_idle_time"           env:"POSTGRES_MAX_IDLE_TIME"        env-default:"15m"`
}

func NewDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
