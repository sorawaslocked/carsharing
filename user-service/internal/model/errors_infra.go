package model

import "errors"

var (
	ErrNats           = errors.New("nats error")
	ErrRedis          = errors.New("redis error")
	ErrSqlTransaction = errors.New("sql transaction error")
	ErrSql            = errors.New("sql error")
	ErrBcrypt         = errors.New("bcrypt error")
)
