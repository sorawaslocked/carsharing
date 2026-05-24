package handler

import (
	"context"
	"sync"

	"carsharing/api-gateway/internal/model"
)

// CarStatusHub routes CarStatusUpdatedEvent events to waiting WebSocket connections
// keyed by car ID. Implements nats/handler.CarStatusEventHandler.
type CarStatusHub struct {
	mu   sync.RWMutex
	subs map[string][]chan model.CarStatusUpdatedEvent
}

func NewCarStatusHub() *CarStatusHub {
	return &CarStatusHub{
		subs: make(map[string][]chan model.CarStatusUpdatedEvent),
	}
}

func (h *CarStatusHub) Subscribe(carID string) (<-chan model.CarStatusUpdatedEvent, func()) {
	ch := make(chan model.CarStatusUpdatedEvent, 1)

	h.mu.Lock()
	h.subs[carID] = append(h.subs[carID], ch)
	h.mu.Unlock()

	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		chans := h.subs[carID]
		for i, c := range chans {
			if c == ch {
				h.subs[carID] = append(chans[:i], chans[i+1:]...)
				break
			}
		}
		if len(h.subs[carID]) == 0 {
			delete(h.subs, carID)
		}
		close(ch)
	}
}

// OnCarStatusUpdated implements nats/handler.CarStatusEventHandler.
func (h *CarStatusHub) OnCarStatusUpdated(_ context.Context, event model.CarStatusUpdatedEvent) error {
	h.mu.RLock()
	chans := h.subs[event.CarID]
	h.mu.RUnlock()

	for _, ch := range chans {
		select {
		case ch <- event:
		default:
		}
	}

	return nil
}
