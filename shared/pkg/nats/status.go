package nats

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
)

const (
	StatusConnected    = "connected"
	StatusDisconnected = "disconnected"
)

func convertStatus(st nats.Status) string {
	switch st {
	case nats.CONNECTED, nats.CONNECTING:
		return StatusConnected
	default:
		return StatusDisconnected
	}
}

func Ping(ctx context.Context, log *slog.Logger, nc *nats.Conn) error {
	log = pkglog.WithMethod(log, "nats.Ping")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	ncStatus := nc.Status()

	if convertStatus(ncStatus) != StatusConnected {
		log.Error(
			"pinging nats",
			slog.String("addr", nc.ConnectedAddr()),
			slog.String("serverName", nc.ConnectedServerName()),
			slog.String("status", ncStatus.String()),
		)

		return ErrFailedConnection
	}

	return nil
}
