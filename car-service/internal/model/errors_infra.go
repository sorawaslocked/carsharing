package model

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNats                        = errors.New("nats error")
	ErrSqlTransaction              = errors.New("sql transaction error")
	ErrSql                         = errors.New("sql error")
	ErrObjectStorage               = errors.New("object storage error")
	ErrTelemetryNoStreamsConnected = errors.New("telemetry: no streams connected")
)

type ErrTelemetryAllStreamsDisconnected struct {
	LastActivity time.Duration
}

func (e ErrTelemetryAllStreamsDisconnected) Error() string {
	return fmt.Sprintf("telemetry: all streams disconnected, last activity %s ago", e.LastActivity)
}

type ErrTelemetryStreamStale struct {
	SinceLastUpdate time.Duration
}

func (e ErrTelemetryStreamStale) Error() string {
	return fmt.Sprintf("telemetry: stream stale, no updates for %s", e.SinceLastUpdate)
}
