//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"carsharing/user-service/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Insert ---

func TestDocumentRepo_Insert_ReturnsID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	userID := mustInsertUser(t, testUser())

	id, err := newDocRepo().Insert(ctx, testDoc(userID))

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestDocumentRepo_Insert_DefaultStatusIsPending(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	userID := mustInsertUser(t, testUser())
	r := newDocRepo()

	id, err := r.Insert(ctx, testDoc(userID))
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, model.DocumentStatusPending, found.Status)
}

// --- FindByID ---

func TestDocumentRepo_FindByID_Found(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	userID := mustInsertUser(t, testUser())
	r := newDocRepo()

	doc := testDoc(userID)
	id, err := r.Insert(ctx, doc)
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, found.ID)
	assert.Equal(t, userID, found.UserID)
	assert.Equal(t, doc.ImageType, found.ImageType)
	assert.Equal(t, doc.Image.Key, found.Image.Key)
	assert.Nil(t, found.Error)
}

func TestDocumentRepo_FindByID_NotFound(t *testing.T) {
	truncate(t)

	_, err := newDocRepo().FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- Find ---

func TestDocumentRepo_Find_ByUserID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()

	userID1 := mustInsertUser(t, testUser())
	userID2 := mustInsertUser(t, testUserWithEmail("other@example.com"))

	_, err := r.Insert(ctx, testDoc(userID1))
	require.NoError(t, err)
	_, err = r.Insert(ctx, testDoc(userID2))
	require.NoError(t, err)

	docs, err := r.Find(ctx, model.DocumentFilter{UserID: &userID1})

	require.NoError(t, err)
	require.Len(t, docs, 1)
	assert.Equal(t, userID1, docs[0].UserID)
}

func TestDocumentRepo_Find_ExcludeStatus(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()
	userID := mustInsertUser(t, testUser())

	// pending document
	_, err := r.Insert(ctx, testDoc(userID))
	require.NoError(t, err)

	// approved document (different image type to avoid conflicts)
	idApproved, err := r.Insert(ctx, testDocOfType(userID, model.DocumentImageTypeIDBack))
	require.NoError(t, err)
	approved := model.DocumentStatusApproved
	require.NoError(t, r.Update(ctx, idApproved, model.DocumentUpdate{Status: &approved, UpdatedAt: time.Now()}))

	docs, err := r.Find(ctx, model.DocumentFilter{
		UserID:        &userID,
		ExcludeStatus: ptr(model.DocumentStatusPending),
	})

	require.NoError(t, err)
	require.Len(t, docs, 1)
	assert.Equal(t, model.DocumentStatusApproved, docs[0].Status)
}

func TestDocumentRepo_Find_Empty(t *testing.T) {
	truncate(t)
	userID := mustInsertUser(t, testUser())

	docs, err := newDocRepo().Find(context.Background(), model.DocumentFilter{UserID: &userID})

	require.NoError(t, err)
	assert.Empty(t, docs)
}

func TestDocumentRepo_Find_LatestPerType_ReturnsNewest(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()
	userID := mustInsertUser(t, testUser())

	now := time.Now()
	older := now.Add(-time.Hour)
	newer := now

	// Two id_front documents inserted out of order by timestamp
	_, err := r.Insert(ctx, testDocAt(userID, model.DocumentImageTypeIDFront, older))
	require.NoError(t, err)
	newerID, err := r.Insert(ctx, testDocAt(userID, model.DocumentImageTypeIDFront, newer))
	require.NoError(t, err)

	// One id_back document
	_, err = r.Insert(ctx, testDocAt(userID, model.DocumentImageTypeIDBack, now))
	require.NoError(t, err)

	docs, err := r.Find(ctx, model.DocumentFilter{UserID: &userID, LatestPerType: true})

	require.NoError(t, err)
	require.Len(t, docs, 2) // one per type

	for _, d := range docs {
		if d.ImageType == model.DocumentImageTypeIDFront {
			assert.Equal(t, newerID, d.ID, "expected newest id_front document")
		}
	}
}

func TestDocumentRepo_Find_LatestPerType_MultipleTypes(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()
	userID := mustInsertUser(t, testUser())

	allTypes := model.AllDocumentImageTypes()
	for _, it := range allTypes {
		_, err := r.Insert(ctx, testDocOfType(userID, it))
		require.NoError(t, err)
	}

	docs, err := r.Find(ctx, model.DocumentFilter{UserID: &userID, LatestPerType: true})

	require.NoError(t, err)
	assert.Len(t, docs, len(allTypes))
}

// --- Update ---

func TestDocumentRepo_Update_Status(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()
	userID := mustInsertUser(t, testUser())

	id, err := r.Insert(ctx, testDoc(userID))
	require.NoError(t, err)

	approved := model.DocumentStatusApproved
	err = r.Update(ctx, id, model.DocumentUpdate{Status: &approved, UpdatedAt: time.Now()})
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, model.DocumentStatusApproved, found.Status)
}

func TestDocumentRepo_Update_StatusAndError(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()
	userID := mustInsertUser(t, testUser())

	id, err := r.Insert(ctx, testDoc(userID))
	require.NoError(t, err)

	rejected := model.DocumentStatusRejected
	errMsg := "blurry image"
	err = r.Update(ctx, id, model.DocumentUpdate{
		Status:    &rejected,
		Error:     &errMsg,
		UpdatedAt: time.Now(),
	})
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, model.DocumentStatusRejected, found.Status)
	require.NotNil(t, found.Error)
	assert.Equal(t, errMsg, *found.Error)
}

func TestDocumentRepo_Update_NotFound(t *testing.T) {
	truncate(t)
	approved := model.DocumentStatusApproved

	err := newDocRepo().Update(context.Background(), "00000000-0000-0000-0000-000000000000", model.DocumentUpdate{
		Status:    &approved,
		UpdatedAt: time.Now(),
	})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestDocumentRepo_Delete_CascadesWithUser(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newDocRepo()
	userID := mustInsertUser(t, testUser())

	id, err := r.Insert(ctx, testDoc(userID))
	require.NoError(t, err)

	// Deleting the owning user should cascade to documents
	require.NoError(t, newUserRepo().Delete(ctx, userID))

	_, err = r.FindByID(ctx, id)
	assert.ErrorIs(t, err, model.ErrNotFound)
}
