package service

import (
	"context"
	"testing"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/service/mocks"
	"carsharing/car-service/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestZoneService(t *testing.T, zoneRepo ZoneRepository) *ZoneService {
	t.Helper()
	return NewZoneService(discardLogger(), newTestValidator(t), zoneRepo)
}

func TestZoneServiceCreate(t *testing.T) {
	ctx := context.Background()

	validInput := validation.ZoneCreate{
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

		repo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrSql)

		_, err := svc.Create(ctx, validInput)
		assert.Error(t, err)
	})

	t.Run("validation rejects missing name", func(t *testing.T) {
		svc := newTestZoneService(t, nil)

		_, err := svc.Create(ctx, validation.ZoneCreate{
			Type:            string(model.ZoneTypeOperating),
			BoundaryGeoJSON: `{}`,
		})
		assert.Error(t, err)
	})
}

func TestZoneServiceGet(t *testing.T) {
	ctx := context.Background()
	zoneID := "b0000000-0000-4000-8000-000000000001"

	t.Run("returns zone", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().FindByID(ctx, zoneID).Return(model.Zone{ID: zoneID, Name: "Uptown"}, nil)

		got, err := svc.Get(ctx, zoneID)
		assert.NoError(t, err)
		assert.Equal(t, zoneID, got.ID)
	})

	t.Run("not found returns ErrZoneNotFound", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().FindByID(ctx, zoneID).Return(model.Zone{}, model.ErrZoneNotFound)

		_, err := svc.Get(ctx, zoneID)
		assert.ErrorIs(t, err, model.ErrZoneNotFound)
	})
}

func TestZoneServiceGetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, nil)

		got, err := svc.List(ctx, validation.ZoneFilter{})
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

		got, err := svc.List(ctx, validation.ZoneFilter{IsActive: &active})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrSql)

		_, err := svc.List(ctx, validation.ZoneFilter{})
		assert.Error(t, err)
	})
}

func TestZoneServiceUpdate(t *testing.T) {
	ctx := context.Background()
	zoneID := "b0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		name := "New Name"
		repo.EXPECT().Update(ctx, zoneID, mock.MatchedBy(func(u model.ZoneUpdate) bool {
			return u.Name != nil && *u.Name == "New Name"
		})).Return(nil)

		assert.NoError(t, svc.Update(ctx, zoneID, validation.ZoneUpdate{Name: &name}))
	})

	t.Run("type string is parsed to enum", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		zoneType := string(model.ZoneTypeNoDrop)
		repo.EXPECT().Update(ctx, zoneID, mock.MatchedBy(func(u model.ZoneUpdate) bool {
			return u.Type != nil && *u.Type == model.ZoneTypeNoDrop
		})).Return(nil)

		assert.NoError(t, svc.Update(ctx, zoneID, validation.ZoneUpdate{Type: &zoneType}))
	})

	t.Run("not found returns ErrZoneNotFound", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Update(ctx, zoneID, mock.Anything).Return(model.ErrZoneNotFound)

		assert.ErrorIs(t, svc.Update(ctx, zoneID, validation.ZoneUpdate{}), model.ErrZoneNotFound)
	})
}

func TestZoneServiceGetZonePricing(t *testing.T) {
	ctx := context.Background()
	lat, lng := 51.1801, 71.4460

	t.Run("no zone at location returns zero fee", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().FindByLocation(ctx, lat, lng).Return(nil, nil)

		fee, err := svc.GetZonePricing(ctx, validation.ZoneGetPricing{Latitude: lat, Longitude: lng})
		assert.NoError(t, err)
		assert.Equal(t, int32(0), fee)
	})

	t.Run("operating zone returns fee adjustment", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		zone := model.Zone{Type: model.ZoneTypeOperating, FeeAdjustment: 500}
		repo.EXPECT().FindByLocation(ctx, lat, lng).Return(&zone, nil)

		fee, err := svc.GetZonePricing(ctx, validation.ZoneGetPricing{Latitude: lat, Longitude: lng})
		assert.NoError(t, err)
		assert.Equal(t, int32(500), fee)
	})

	t.Run("parking hub returns negative fee adjustment", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		zone := model.Zone{Type: model.ZoneParkingHub, FeeAdjustment: -200}
		repo.EXPECT().FindByLocation(ctx, lat, lng).Return(&zone, nil)

		fee, err := svc.GetZonePricing(ctx, validation.ZoneGetPricing{Latitude: lat, Longitude: lng})
		assert.NoError(t, err)
		assert.Equal(t, int32(-200), fee)
	})

	t.Run("no_drop zone returns ErrLocationInNoDropZone", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		zone := model.Zone{Type: model.ZoneTypeNoDrop}
		repo.EXPECT().FindByLocation(ctx, lat, lng).Return(&zone, nil)

		_, err := svc.GetZonePricing(ctx, validation.ZoneGetPricing{Latitude: lat, Longitude: lng})
		assert.ErrorIs(t, err, model.ErrLocationInNoDropZone)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().FindByLocation(ctx, lat, lng).Return(nil, model.ErrSql)

		_, err := svc.GetZonePricing(ctx, validation.ZoneGetPricing{Latitude: lat, Longitude: lng})
		assert.ErrorIs(t, err, model.ErrSql)
	})
}

func TestZoneServiceDelete(t *testing.T) {
	ctx := context.Background()
	zoneID := "b0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Delete(ctx, zoneID).Return(nil)

		assert.NoError(t, svc.Delete(ctx, zoneID))
	})

	t.Run("not found returns ErrZoneNotFound", func(t *testing.T) {
		repo := mocks.NewMockZoneRepository(t)
		svc := newTestZoneService(t, repo)

		repo.EXPECT().Delete(ctx, zoneID).Return(model.ErrZoneNotFound)

		assert.ErrorIs(t, svc.Delete(ctx, zoneID), model.ErrZoneNotFound)
	})
}
