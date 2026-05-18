package service_test

import (
	"testing"
	"time"

	"carsharing/user-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- CreateDocument ---

func TestCreateDocument_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.docRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(doc model.Document) bool {
		return doc.UserID == testUserID &&
			doc.ImageType == model.ImageTypeIDFront &&
			doc.Status == model.DocumentStatusPending &&
			doc.Image != nil && doc.Image.Key == testObjKey
	})).Return(testDocID, nil)
	d.analyzer.EXPECT().Analyze(ctx, testDocID, testObjKey)

	id, err := svc.CreateDocument(ctx, testObjKey, model.ImageTypeIDFront)

	require.NoError(t, err)
	assert.Equal(t, testDocID, id)
}

func TestCreateDocument_NoUserIDInCtx(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	_, err := svc.CreateDocument(ctxAnon(), testObjKey, model.ImageTypeIDFront)

	assert.ErrorIs(t, err, model.ErrUnauthenticated)
}

// --- GetDocumentImageUploadData ---

func TestGetDocumentImageUploadData_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	expected := model.ImageUploadData{PresignedPutURL: "https://minio/put", ObjectKey: "documents/id_front/key"}

	d.storage.EXPECT().GetDocumentImageUploadData(ctx, testImgType).Return(expected, nil)

	got, err := svc.GetDocumentImageUploadData(ctx, testImgType)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

// --- GetProcessedDocumentsForUser ---

func TestGetProcessedDocumentsForUser_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	pending := model.DocumentStatusPending
	docs := []model.Document{
		{ID: testDocID, UserID: testUserID, Status: model.DocumentStatusApproved, ImageType: model.ImageTypeIDFront},
	}

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(baseUser(), nil)
	d.docRepo.EXPECT().Find(ctx, model.DocumentFilter{
		UserID:        ptr(testUserID),
		ExcludeStatus: &pending,
		LatestPerType: true,
	}).Return(docs, nil)

	got, err := svc.GetProcessedDocumentsForUser(ctx, testUserID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, testDocID, got[0].ID)
}

func TestGetProcessedDocumentsForUser_UserNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(model.User{}, model.ErrNotFound)

	_, err := svc.GetProcessedDocumentsForUser(ctx, testUserID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- CheckDocument ---

func TestCheckDocument_Rejected(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	rejected := model.DocumentStatusRejected
	errMsg := "blurry image"
	doc := model.Document{ID: testDocID, UserID: testUserID, Status: model.DocumentStatusPending, ImageType: model.ImageTypeIDFront}

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(doc, nil)
	d.docRepo.EXPECT().Update(ctx, testDocID, mock.MatchedBy(func(u model.DocumentUpdate) bool {
		return u.Status != nil && *u.Status == rejected && u.Error != nil && *u.Error == errMsg
	})).Return(nil)

	err := svc.CheckDocument(ctx, testDocID, rejected, ptr(errMsg))

	require.NoError(t, err)
}

func TestCheckDocument_Approved_NotAllTypesPresent(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	approved := model.DocumentStatusApproved
	doc := model.Document{ID: testDocID, UserID: testUserID, Status: model.DocumentStatusPending, ImageType: model.ImageTypeIDFront}

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(doc, nil)
	d.docRepo.EXPECT().Update(ctx, testDocID, mock.MatchedBy(func(u model.DocumentUpdate) bool {
		return u.Status != nil && *u.Status == approved
	})).Return(nil)
	// Only one document type returned — not enough to set IsDocumentVerified.
	d.docRepo.EXPECT().Find(ctx, model.DocumentFilter{UserID: ptr(testUserID), LatestPerType: true}).
		Return([]model.Document{doc}, nil)

	err := svc.CheckDocument(ctx, testDocID, approved, nil)

	require.NoError(t, err)
}

func TestCheckDocument_Approved_AllTypesApproved_SetsVerified(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	approved := model.DocumentStatusApproved

	doc := model.Document{ID: testDocID, UserID: testUserID, Status: model.DocumentStatusPending, ImageType: model.ImageTypeIDFront}

	// All four document types present and approved.
	allApproved := make([]model.Document, len(model.AllImageTypes()))
	for i, it := range model.AllImageTypes() {
		allApproved[i] = model.Document{UserID: testUserID, ImageType: it, Status: model.DocumentStatusApproved}
	}

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(doc, nil)
	d.docRepo.EXPECT().Update(ctx, testDocID, mock.MatchedBy(func(u model.DocumentUpdate) bool {
		return u.Status != nil && *u.Status == approved
	})).Return(nil)
	d.docRepo.EXPECT().Find(ctx, model.DocumentFilter{UserID: ptr(testUserID), LatestPerType: true}).
		Return(allApproved, nil)
	d.userRepo.EXPECT().Update(ctx, testUserID, mock.MatchedBy(func(u model.UserRepoUpdate) bool {
		return u.IsDocumentVerified != nil && *u.IsDocumentVerified
	})).Return(nil)
	d.publisher.EXPECT().PublishUserUpdated(ctx, testUserID, false).Return(nil)

	err := svc.CheckDocument(ctx, testDocID, approved, nil)

	require.NoError(t, err)
}

func TestCheckDocument_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(model.Document{}, model.ErrNotFound)

	err := svc.CheckDocument(ctx, testDocID, model.DocumentStatusApproved, nil)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- HandleDocumentAnalyzed ---

func TestHandleDocumentAnalyzed_Passed(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	doc := model.Document{ID: testDocID, UserID: testUserID, ImageType: model.ImageTypeIDFront, Status: model.DocumentStatusPending}
	event := model.DocumentAnalyzedEvent{DocumentID: testDocID, Passed: true}

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(doc, nil)
	d.docRepo.EXPECT().Update(ctx, testDocID, mock.MatchedBy(func(u model.DocumentUpdate) bool {
		return u.Status != nil && *u.Status == model.DocumentStatusApproved && u.Error == nil
	})).Return(nil)
	// Only one doc type — won't trigger verification.
	d.docRepo.EXPECT().Find(ctx, model.DocumentFilter{UserID: ptr(testUserID), LatestPerType: true}).
		Return([]model.Document{doc}, nil)

	err := svc.HandleDocumentAnalyzed(ctx, event)

	require.NoError(t, err)
}

func TestHandleDocumentAnalyzed_Failed_WithDefects(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	doc := model.Document{ID: testDocID, UserID: testUserID, ImageType: model.ImageTypeIDFront, Status: model.DocumentStatusPending}
	event := model.DocumentAnalyzedEvent{
		DocumentID: testDocID,
		Passed:     false,
		Defects:    []model.Defect{{Type: "blur", Description: "image is blurry"}},
	}

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(doc, nil)
	d.docRepo.EXPECT().Update(ctx, testDocID, mock.MatchedBy(func(u model.DocumentUpdate) bool {
		return u.Status != nil && *u.Status == model.DocumentStatusRejected &&
			u.Error != nil && *u.Error != ""
	})).Return(nil)

	err := svc.HandleDocumentAnalyzed(ctx, event)

	require.NoError(t, err)
}

func TestHandleDocumentAnalyzed_DocumentNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(model.Document{}, model.ErrNotFound)

	err := svc.HandleDocumentAnalyzed(ctx, model.DocumentAnalyzedEvent{DocumentID: testDocID, Passed: true})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// Ensure timestamps in DocumentUpdate are not zero.
func TestCheckDocument_UpdateTimestampSet(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	rejected := model.DocumentStatusRejected
	doc := model.Document{ID: testDocID, UserID: testUserID, Status: model.DocumentStatusPending, ImageType: model.ImageTypeIDFront}

	d.docRepo.EXPECT().FindByID(ctx, testDocID).Return(doc, nil)
	d.docRepo.EXPECT().Update(ctx, testDocID, mock.MatchedBy(func(u model.DocumentUpdate) bool {
		return !u.UpdatedAt.IsZero() && u.UpdatedAt.Before(time.Now().Add(time.Second))
	})).Return(nil)

	err := svc.CheckDocument(ctx, testDocID, rejected, nil)

	require.NoError(t, err)
}
