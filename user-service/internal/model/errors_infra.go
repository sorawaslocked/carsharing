package model

import "errors"

var (
	ErrBrevo          = errors.New("brevo error")
	ErrNats           = errors.New("nats error")
	ErrRedis          = errors.New("redis error")
	ErrSqlTransaction = errors.New("sql transaction error")
	ErrSql            = errors.New("sql error")
	ErrBcrypt         = errors.New("bcrypt error")
	ErrObjectStorage  = errors.New("object storage error")
)
