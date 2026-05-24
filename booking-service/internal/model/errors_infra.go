package model

import "errors"

var (
	ErrNats                = errors.New("nats error")
	ErrSqlTransaction      = errors.New("sql transaction error")
	ErrSql                 = errors.New("sql error")
	ErrInternalServerError = errors.New("internal server error")
)
