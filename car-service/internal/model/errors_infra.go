package model

import "errors"

var (
	ErrNats          = errors.New("nats error")
	ErrSql           = errors.New("sql error")
	ErrObjectStorage = errors.New("object storage error")
)
