package documentanalyzer

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Checker struct {
	conn *grpc.ClientConn
}

func NewChecker(conn *grpc.ClientConn) *Checker {
	return &Checker{conn: conn}
}

func (c *Checker) Ping(_ context.Context) error {
	state := c.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return fmt.Errorf("connection state: %s", state)
	}
	return nil
}
