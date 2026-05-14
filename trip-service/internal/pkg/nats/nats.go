package nats

import (
	"time"

	natsio "github.com/nats-io/nats.go"
)

func NewConn(url string) (*natsio.Conn, error) {
	return natsio.Connect(url,
		natsio.MaxReconnects(-1),
		natsio.ReconnectWait(2*time.Second),
	)
}
