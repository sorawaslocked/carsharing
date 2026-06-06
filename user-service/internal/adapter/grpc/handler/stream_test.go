package handler_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	usersvc "carsharing/protos/gen/service/user"
	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/adapter/grpc/handler"
	"carsharing/user-service/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- fakes ----------------------------------------------------------------

// fakeSubscriber captures the filter args passed to SubscribeStream and serves
// events from a pre-built channel.
type fakeSubscriber struct {
	ch             chan model.DocumentAnalyzedEvent
	capturedUserID *string
	capturedPassed *bool
}

func newFakeSub(events ...model.DocumentAnalyzedEvent) *fakeSubscriber {
	ch := make(chan model.DocumentAnalyzedEvent, len(events)+1)
	for _, e := range events {
		ch <- e
	}
	close(ch)
	return &fakeSubscriber{ch: ch}
}

func (f *fakeSubscriber) SubscribeStream(userID *string, passed *bool) (<-chan model.DocumentAnalyzedEvent, func()) {
	f.capturedUserID = userID
	f.capturedPassed = passed
	return f.ch, func() {}
}

// fakeDocStream implements grpc.ServerStreamingServer[usersvc.StreamDocumentAnalyzedResponse].
type fakeDocStream struct {
	ctx  context.Context
	sent []*usersvc.StreamDocumentAnalyzedResponse
}

func newFakeStream(ctx context.Context) *fakeDocStream {
	return &fakeDocStream{ctx: ctx}
}

func (s *fakeDocStream) Send(r *usersvc.StreamDocumentAnalyzedResponse) error {
	s.sent = append(s.sent, r)
	return nil
}
func (s *fakeDocStream) Context() context.Context     { return s.ctx }
func (s *fakeDocStream) SetHeader(metadata.MD) error  { return nil }
func (s *fakeDocStream) SendHeader(metadata.MD) error { return nil }
func (s *fakeDocStream) SetTrailer(metadata.MD)       {}
func (s *fakeDocStream) SendMsg(interface{}) error    { return nil }
func (s *fakeDocStream) RecvMsg(interface{}) error    { return nil }

// ---- helpers ---------------------------------------------------------------

func newStreamHandler(t *testing.T, sub handler.DocumentAnalyzedSubscriber) *handler.UserHandler {
	t.Helper()
	return handler.NewUserHandler(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		nil, // UserService not needed for stream tests
		sub,
	)
}

func ctxWithAdmin(userID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyUserRoles, []sharedmodel.Role{sharedmodel.RoleAdmin})
	return ctx
}

// ---- tests -----------------------------------------------------------------

func TestStreamDocumentAnalyzed_OwnEventsDelivered(t *testing.T) {
	event := model.DocumentAnalyzedEvent{
		DocumentID: testDocID,
		UserID:     testUserID,
		Passed:     true,
	}
	sub := newFakeSub(event)
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithUser(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{}, stream)

	require.NoError(t, err)
	require.Len(t, stream.sent, 1)
	assert.Equal(t, testDocID, stream.sent[0].DocumentId)
	assert.Equal(t, testUserID, stream.sent[0].UserId)
	assert.True(t, stream.sent[0].Passed)
}

func TestStreamDocumentAnalyzed_NilUserIDDefaultsToOwn(t *testing.T) {
	sub := newFakeSub()
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithUser(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{UserId: nil}, stream)

	require.NoError(t, err)
	require.NotNil(t, sub.capturedUserID, "userID filter must be set")
	assert.Equal(t, testUserID, *sub.capturedUserID)
}

func TestStreamDocumentAnalyzed_ForeignUserIDForbidden(t *testing.T) {
	otherUserID := "other-user-id"
	sub := newFakeSub()
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithUser(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{UserId: ptr(otherUserID)}, stream)

	assert.Equal(t, codes.PermissionDenied, grpcCode(err))
}

func TestStreamDocumentAnalyzed_NoAuthUserIDReturnsUnauthenticated(t *testing.T) {
	sub := newFakeSub()
	h := newStreamHandler(t, sub)

	// context has no user ID at all
	stream := newFakeStream(ctxAnon())
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{}, stream)

	assert.Equal(t, codes.Unauthenticated, grpcCode(err))
}

func TestStreamDocumentAnalyzed_PrivilegedUserCanFilterByOtherUserID(t *testing.T) {
	otherUserID := "other-user-id"
	event := model.DocumentAnalyzedEvent{DocumentID: testDocID, UserID: otherUserID, Passed: true}
	sub := newFakeSub(event)
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithAdmin(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{UserId: ptr(otherUserID)}, stream)

	require.NoError(t, err)
	require.NotNil(t, sub.capturedUserID)
	assert.Equal(t, otherUserID, *sub.capturedUserID)
	require.Len(t, stream.sent, 1)
	assert.Equal(t, testDocID, stream.sent[0].DocumentId)
}

func TestStreamDocumentAnalyzed_PassedFilterForwarded(t *testing.T) {
	sub := newFakeSub()
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithUser(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{Passed: ptr(true)}, stream)

	require.NoError(t, err)
	require.NotNil(t, sub.capturedPassed)
	assert.True(t, *sub.capturedPassed)
}

func TestStreamDocumentAnalyzed_MultipleEventsDeliveredInOrder(t *testing.T) {
	e1 := model.DocumentAnalyzedEvent{DocumentID: "doc-1", UserID: testUserID, Passed: true}
	e2 := model.DocumentAnalyzedEvent{DocumentID: "doc-2", UserID: testUserID, Passed: false}
	sub := newFakeSub(e1, e2)
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithUser(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{}, stream)

	require.NoError(t, err)
	require.Len(t, stream.sent, 2)
	assert.Equal(t, "doc-1", stream.sent[0].DocumentId)
	assert.Equal(t, "doc-2", stream.sent[1].DocumentId)
}

func TestStreamDocumentAnalyzed_ContextCancellationEndsStream(t *testing.T) {
	// channel stays open — only context cancellation should end the stream
	ch := make(chan model.DocumentAnalyzedEvent)
	sub := &fakeSubscriber{ch: ch}
	h := newStreamHandler(t, sub)

	ctx, cancel := context.WithCancel(ctxWithUser(testUserID))
	stream := newFakeStream(ctx)

	done := make(chan error, 1)
	go func() {
		done <- h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{}, stream)
	}()

	cancel()
	err := <-done

	assert.NoError(t, err)
	assert.Empty(t, stream.sent)
}

func TestStreamDocumentAnalyzed_DefectsAreForwarded(t *testing.T) {
	event := model.DocumentAnalyzedEvent{
		DocumentID: testDocID,
		UserID:     testUserID,
		Passed:     false,
		Defects: []model.Defect{
			{Type: "blur", Description: "image is blurry"},
		},
	}
	sub := newFakeSub(event)
	h := newStreamHandler(t, sub)

	stream := newFakeStream(ctxWithUser(testUserID))
	err := h.StreamDocumentAnalyzed(&usersvc.StreamDocumentAnalyzedRequest{}, stream)

	require.NoError(t, err)
	require.Len(t, stream.sent, 1)
	require.Len(t, stream.sent[0].Defects, 1)
	assert.Equal(t, "blur", stream.sent[0].Defects[0].Type)
	assert.Equal(t, "image is blurry", stream.sent[0].Defects[0].Description)
}
