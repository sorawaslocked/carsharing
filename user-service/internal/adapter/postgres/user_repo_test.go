//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Insert ---

func TestUserRepo_Insert_ReturnsID(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	id, err := newUserRepo().Insert(ctx, testUser())

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestUserRepo_Insert_RolesPersistedCorrectly(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	u := testUser()
	u.Roles = []sharedmodel.Role{sharedmodel.RoleUser, sharedmodel.RoleAdmin}
	id, err := r.Insert(ctx, u)
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.ElementsMatch(t, u.Roles, found.Roles)
}

func TestUserRepo_Insert_NilPhoneNumber(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	u := testUser()
	u.PhoneNumber = nil
	id, err := r.Insert(ctx, u)
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, found.PhoneNumber)
}

func TestUserRepo_Insert_DuplicateEmail(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	_, err := r.Insert(ctx, testUser())
	require.NoError(t, err)

	u2 := testUser()
	u2.PhoneNumber = nil
	_, err = r.Insert(ctx, u2)

	assert.ErrorIs(t, err, model.ErrDuplicateEmail)
}

func TestUserRepo_Insert_DuplicatePhone(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	_, err := r.Insert(ctx, testUser())
	require.NoError(t, err)

	u2 := testUser()
	u2.Email = "other@example.com"
	_, err = r.Insert(ctx, u2)

	assert.ErrorIs(t, err, model.ErrDuplicatePhone)
}

// --- FindByID ---

func TestUserRepo_FindByID_Found(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	u := testUser()
	id, err := r.Insert(ctx, u)
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, found.ID)
	assert.Equal(t, u.Email, found.Email)
	assert.Equal(t, u.FirstName, found.FirstName)
	assert.Equal(t, u.LastName, found.LastName)
	assert.Equal(t, u.PhoneNumber, found.PhoneNumber)
	assert.Equal(t, u.BirthDate, found.BirthDate)
}

func TestUserRepo_FindByID_NotFound(t *testing.T) {
	truncate(t)

	_, err := newUserRepo().FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrUserNotFound)
}

// --- FindOne ---

func TestUserRepo_FindOne_ByEmail(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	u := testUser()
	id, err := r.Insert(ctx, u)
	require.NoError(t, err)

	found, err := r.FindOne(ctx, model.UserFilter{Email: &u.Email})

	require.NoError(t, err)
	assert.Equal(t, id, found.ID)
	assert.Equal(t, u.Email, found.Email)
}

func TestUserRepo_FindOne_ByPhone(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	u := testUser()
	id, err := r.Insert(ctx, u)
	require.NoError(t, err)

	found, err := r.FindOne(ctx, model.UserFilter{PhoneNumber: u.PhoneNumber})

	require.NoError(t, err)
	assert.Equal(t, id, found.ID)
}

func TestUserRepo_FindOne_NotFound(t *testing.T) {
	truncate(t)

	email := "ghost@example.com"
	_, err := newUserRepo().FindOne(context.Background(), model.UserFilter{Email: &email})

	assert.ErrorIs(t, err, model.ErrUserNotFound)
}

// --- Find ---

func TestUserRepo_Find_Empty(t *testing.T) {
	truncate(t)

	users, err := newUserRepo().Find(context.Background(), model.UserFilter{})

	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepo_Find_ReturnsAll(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	mustInsertUser(t, testUser())
	mustInsertUser(t, testUserWithEmail("other@example.com"))

	users, err := r.Find(ctx, model.UserFilter{})

	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepo_Find_FilterByIsEmailVerified(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	mustInsertUser(t, testUser())

	u2 := testUserWithEmail("verified@example.com")
	u2.IsEmailVerified = true
	mustInsertUser(t, u2)

	users, err := r.Find(ctx, model.UserFilter{IsEmailVerified: ptr(true)})

	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, u2.Email, users[0].Email)
}

func TestUserRepo_Find_FilterByIsSuspended(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	mustInsertUser(t, testUser())

	suspended := testUserWithEmail("suspended@example.com")
	suspended.IsSuspended = true
	mustInsertUser(t, suspended)

	users, err := r.Find(ctx, model.UserFilter{IsSuspended: ptr(true)})

	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, suspended.Email, users[0].Email)
}

func TestUserRepo_Find_Pagination(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	for i := range 4 {
		mustInsertUser(t, testUserWithEmail(t.Name()+string(rune('a'+i))+"@x.com"))
	}

	page1, err := r.Find(ctx, model.UserFilter{Pagination: &sharedmodel.Pagination{Limit: 2, Offset: 0}})
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := r.Find(ctx, model.UserFilter{Pagination: &sharedmodel.Pagination{Limit: 2, Offset: 2}})
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

// --- Update ---

func TestUserRepo_Update_ScalarFields(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	id := mustInsertUser(t, testUser())

	err := r.Update(ctx, id, model.UserUpdate{
		FirstName: ptr("Jane"),
		LastName:  ptr("Smith"),
		UpdatedAt: time.Now(),
	})
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Jane", found.FirstName)
	assert.Equal(t, "Smith", found.LastName)
}

func TestUserRepo_Update_Roles_Replaces(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	id := mustInsertUser(t, testUser()) // inserted with RoleUser

	err := r.Update(ctx, id, model.UserUpdate{
		Roles:     []sharedmodel.Role{sharedmodel.RoleAdmin},
		UpdatedAt: time.Now(),
	})
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	require.Len(t, found.Roles, 1)
	assert.Equal(t, sharedmodel.RoleAdmin, found.Roles[0])
}

func TestUserRepo_Update_IsEmailVerified(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	id := mustInsertUser(t, testUser())

	err := r.Update(ctx, id, model.UserUpdate{
		IsEmailVerified: ptr(true),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	found, err := r.FindByID(ctx, id)
	require.NoError(t, err)
	assert.True(t, found.IsEmailVerified)
}

func TestUserRepo_Update_DuplicateEmail(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	mustInsertUser(t, testUser())
	id2 := mustInsertUser(t, testUserWithEmail("other@example.com"))

	existing := "john@example.com"
	err := r.Update(ctx, id2, model.UserUpdate{Email: &existing, UpdatedAt: time.Now()})

	assert.ErrorIs(t, err, model.ErrDuplicateEmail)
}

func TestUserRepo_Update_NotFound(t *testing.T) {
	truncate(t)

	err := newUserRepo().Update(context.Background(), "00000000-0000-0000-0000-000000000000", model.UserUpdate{
		FirstName: ptr("Jane"),
		UpdatedAt: time.Now(),
	})

	assert.ErrorIs(t, err, model.ErrUserNotFound)
}

// --- Delete ---

func TestUserRepo_Delete_RemovesUser(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	id := mustInsertUser(t, testUser())

	err := r.Delete(ctx, id)
	require.NoError(t, err)

	_, err = r.FindByID(ctx, id)
	assert.ErrorIs(t, err, model.ErrUserNotFound)
}

func TestUserRepo_Delete_CascadesToRoles(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newUserRepo()

	id := mustInsertUser(t, testUser())
	require.NoError(t, r.Delete(ctx, id))

	// user_roles row should be gone — re-inserting same user must not hit a stale FK
	id2, err := r.Insert(ctx, testUser())
	require.NoError(t, err)
	assert.NotEmpty(t, id2)
}

func TestUserRepo_Delete_NotFound(t *testing.T) {
	truncate(t)

	err := newUserRepo().Delete(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrUserNotFound)
}
