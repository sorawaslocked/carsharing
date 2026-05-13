package handler

import (
	"context"
	"sync"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

// DocumentHub routes DocumentAnalyzedEvent events to waiting WebSocket connections
// keyed by document ID. Implements nats/handler.DocumentEventHandler.
type DocumentHub struct {
	mu   sync.RWMutex
	subs map[string][]chan model.DocumentAnalyzedEvent
}

func NewDocumentHub() *DocumentHub {
	return &DocumentHub{
		subs: make(map[string][]chan model.DocumentAnalyzedEvent),
	}
}

// Subscribe registers a receive channel for the given document ID and returns
// an unsubscribe function the caller must defer.
func (h *DocumentHub) Subscribe(docID string) (<-chan model.DocumentAnalyzedEvent, func()) {
	ch := make(chan model.DocumentAnalyzedEvent, 1)

	h.mu.Lock()
	h.subs[docID] = append(h.subs[docID], ch)
	h.mu.Unlock()

	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		chans := h.subs[docID]
		for i, c := range chans {
			if c == ch {
				h.subs[docID] = append(chans[:i], chans[i+1:]...)
				break
			}
		}
		if len(h.subs[docID]) == 0 {
			delete(h.subs, docID)
		}
		close(ch)
	}
}

// OnDocumentAnalyzed implements nats/handler.DocumentEventHandler.
func (h *DocumentHub) OnDocumentAnalyzed(_ context.Context, event model.DocumentAnalyzedEvent) error {
	h.mu.RLock()
	chans := h.subs[event.DocumentID]
	h.mu.RUnlock()

	for _, ch := range chans {
		select {
		case ch <- event:
		default:
		}
	}

	return nil
}
