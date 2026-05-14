package nats

import (
	"fmt"

	natsgo "github.com/nats-io/nats.go"
	"github.com/sorawaslocked/car-rental-booking-service/internal/config"
)

func New(cfg config.NATS) (*natsgo.Conn, error) {
	nc, err := natsgo.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	return nc, nil
}
