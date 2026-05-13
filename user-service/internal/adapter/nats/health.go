package nats

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
)

type Checker struct {
	conn *nats.Conn
}

func NewChecker(conn *nats.Conn) *Checker {
	return &Checker{conn: conn}
}

func (c *Checker) Ping(_ context.Context) error {
	if !c.conn.IsConnected() {
		return errors.New("not connected")
	}
	return nil
}
