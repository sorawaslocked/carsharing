package subscriber

import "context"

type TripEventHandler interface {
	Complete(ctx context.Context, bookingID string) error
}
