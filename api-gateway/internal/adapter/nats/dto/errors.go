package dto

import "errors"

var ErrSubscribeFailed = errors.New("nats: failed to subscribe to subject")
