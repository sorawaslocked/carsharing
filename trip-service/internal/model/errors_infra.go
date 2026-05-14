package model

import "errors"

var (
	ErrSQL  = errors.New("sql error")
	ErrNATS = errors.New("nats error")
)
