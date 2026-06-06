package subscriber

import (
	"context"
	"log/slog"
	"sync"

	pkglog "carsharing/shared/pkg/log"
	natsdto "carsharing/user-service/internal/adapter/nats/dto"
	"carsharing/user-service/internal/model"

	eventuserpb "carsharing/protos/gen/event/user"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const subjectDocumentAnalyzed = "document.analyzed"

type docSub struct {
	userID    *string
	passed    *bool
	ch        chan model.DocumentAnalyzedEvent
	closeOnce *sync.Once
}

type DocumentSubscriber struct {
	log     *slog.Logger
	conn    *nats.Conn
	service DocumentService

	mu   sync.RWMutex
	subs []docSub
}

func NewDocumentSubscriber(log *slog.Logger, conn *nats.Conn, service DocumentService) *DocumentSubscriber {
	return &DocumentSubscriber{
		log:     pkglog.WithComponent(log, "adapter.nats.subscriber.DocumentSubscriber"),
		conn:    conn,
		service: service,
	}
}

func (s *DocumentSubscriber) Subscribe() error {
	_, err := s.conn.Subscribe(subjectDocumentAnalyzed, s.handleDocumentAnalyzed)
	return err
}

// SubscribeStream registers a channel that receives DocumentAnalyzedEvents
// matching the optional userID and passed filters. The returned cancel func
// must be deferred by the caller to unregister and close the channel.
// If a subscription for the same userID already exists it is evicted, ensuring
// a reconnecting client only receives events from its new subscription.
func (s *DocumentSubscriber) SubscribeStream(userID *string, passed *bool) (<-chan model.DocumentAnalyzedEvent, func()) {
	ch := make(chan model.DocumentAnalyzedEvent)
	once := &sync.Once{}

	s.mu.Lock()
	// Evict any existing subscription for the same userID so a reconnecting
	// client does not accumulate duplicate, possibly stale-filtered subs.
	if userID != nil {
		kept := s.subs[:0]
		for _, sub := range s.subs {
			if sub.userID != nil && *sub.userID == *userID {
				sub.closeOnce.Do(func() { close(sub.ch) })
				continue
			}
			kept = append(kept, sub)
		}
		s.subs = kept
	}
	s.subs = append(s.subs, docSub{userID: userID, passed: passed, ch: ch, closeOnce: once})
	s.mu.Unlock()

	return ch, func() {
		s.mu.Lock()
		for i, sub := range s.subs {
			if sub.ch == ch {
				s.subs = append(s.subs[:i], s.subs[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		once.Do(func() { close(ch) })
	}
}

func (s *DocumentSubscriber) handleDocumentAnalyzed(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleDocumentAnalyzed")

	var pb eventuserpb.DocumentAnalyzedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling document analyzed event", pkglog.Err(err))
		return
	}

	event := natsdto.DocumentAnalyzedEventFromProto(&pb)

	if err := s.service.HandleDocumentAnalyzed(ctx, event); err != nil {
		log.Error("handling document analyzed event",
			slog.String("documentID", event.DocumentID),
			pkglog.Err(err),
		)
	}

	s.fanOut(event)
}

func (s *DocumentSubscriber) fanOut(event model.DocumentAnalyzedEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.subs {
		if sub.userID != nil && *sub.userID != event.UserID {
			continue
		}
		if sub.passed != nil && *sub.passed != event.Passed {
			continue
		}
		select {
		case sub.ch <- event:
		default:
		}
	}
}
