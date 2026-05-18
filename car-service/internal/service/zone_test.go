package service

import (
	"context"
	"testing"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestZoneService(t *testing.T, zoneRepo ZoneRepository) *ZoneService {
	t.Helper()
	return NewZoneService(zoneRepo, newTestValidator(t), discardLogger())
}

func TestZoneServiceCreate(t *testing.T) {
	ctx := context.Background()

	validInput := model.ZoneCreateInput{
		Name:            "Downtown",
		Type:            string(model.ZoneTypeOperating),
		BoundaryGeoJSON: `{"type":"Polygon","coordinates":[]}`,
	}

	t.Run("happy path returns inserted id", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(z model.Zone) bool {
			return z.Name == "Downtown" && z.IsActive
		})).Return("zone-123", nil)

		id, err := svc.Create(ctx, validInput)
		assert.NoError(t, err)
		assert.Equal(t, "zone-123", id)
	})

	t.Run("new zone is always active", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(z model.Zone) bool {
			return z.IsActive
		})).Return("zone-456", nil)

		_, err := svc.Create(ctx, validInput)
		assert.NoError(t, err)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrInternalServerError)

		_, err := svc.Create(ctx, validInput)
		assert.Error(t, err)
	})

	t.Run("validation rejects missing name", func(t *testing.T) {
		svc := newTestZoneService(t, nil)

		_, err := svc.Create(ctx, model.ZoneCreateInput{
			Type:            string(model.ZoneTypeOperating),
			BoundaryGeoJSON: `{}`,
		})
		assert.Error(t, err)
	})
}

func TestZoneServiceGet(t *testing.T) {
	ctx := context.Background()
	zoneID := "zone-123"

	t.Run("returns zone", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().FindByID(ctx, zoneID).Return(model.Zone{ID: zoneID, Name: "Uptown"}, nil)

		got, err := svc.Get(ctx, zoneID)
		assert.NoError(t, err)
		assert.Equal(t, zoneID, got.ID)
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().FindByID(ctx, zoneID).Return(model.Zone{}, model.ErrNotFound)

		_, err := svc.Get(ctx, zoneID)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestZoneServiceGetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, nil)

		got, err := svc.GetAll(ctx, model.ZoneFilterInput{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("forwards filter to repo", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		active := true
		repo.EXPECT().Find(ctx, mock.MatchedBy(func(f model.ZoneFilter) bool {
			return f.IsActive != nil && *f.IsActive
		})).Return([]model.Zone{{ID: "zone-1"}}, nil)

		got, err := svc.GetAll(ctx, model.ZoneFilterInput{IsActive: &active})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrInternalServerError)

		_, err := svc.GetAll(ctx, model.ZoneFilterInput{})
		assert.Error(t, err)
	})
}

func TestZoneServiceUpdate(t *testing.T) {
	ctx := context.Background()
	zoneID := "zone-123"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		name := "New Name"
		repo.EXPECT().Update(ctx, zoneID, mock.MatchedBy(func(u model.ZoneUpdate) bool {
			return u.Name != nil && *u.Name == "New Name"
		})).Return(nil)

		assert.NoError(t, svc.Update(ctx, zoneID, model.ZoneUpdateInput{Name: &name}))
	})

	t.Run("type string is parsed to enum", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		zoneType := string(model.ZoneTypeNoDrop)
		repo.EXPECT().Update(ctx, zoneID, mock.MatchedBy(func(u model.ZoneUpdate) bool {
			return u.Type != nil && *u.Type == model.ZoneTypeNoDrop
		})).Return(nil)

		assert.NoError(t, svc.Update(ctx, zoneID, model.ZoneUpdateInput{Type: &zoneType}))
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Update(ctx, zoneID, mock.Anything).Return(model.ErrNotFound)

		assert.ErrorIs(t, svc.Update(ctx, zoneID, model.ZoneUpdateInput{}), model.ErrNotFound)
	})
}

func TestZoneServiceDelete(t *testing.T) {
	ctx := context.Background()
	zoneID := "zone-123"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Delete(ctx, zoneID).Return(nil)

		assert.NoError(t, svc.Delete(ctx, zoneID))
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Delete(ctx, zoneID).Return(model.ErrNotFound)

		assert.ErrorIs(t, svc.Delete(ctx, zoneID), model.ErrNotFound)
	})
}
