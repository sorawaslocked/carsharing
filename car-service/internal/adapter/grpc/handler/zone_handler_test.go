package handler

import (
	"context"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	carsvc "carsharing/protos/gen/service/car"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestZoneHandlerCreateZone(t *testing.T) {
	ctx := context.Background()

	t.Run("returns id from service", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.MatchedBy(func(in validation.ZoneCreate) bool {
			return in.Name == "Downtown" && in.Type == "operating"
		})).Return("zone-123", nil)

		resp, err := h.CreateZone(ctx, &carsvc.CreateZoneRequest{
			Name: "Downtown", Type: "operating", BoundaryGeoJson: `{}`,
		})
		assert.NoError(t, err)
		assert.Equal(t, "zone-123", resp.Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Create(ctx, mock.Anything).Return("", errInternal)

		_, err := h.CreateZone(ctx, &carsvc.CreateZoneRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestZoneHandlerGetZone(t *testing.T) {
	ctx := context.Background()
	zoneID := "zone-123"

	t.Run("returns populated zone proto", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Get(ctx, zoneID).Return(model.Zone{
			ID: zoneID, Name: "Uptown", Type: model.ZoneTypeOperating, IsActive: true,
		}, nil)

		resp, err := h.GetZone(ctx, &carsvc.GetZoneRequest{Id: zoneID})
		assert.NoError(t, err)
		assert.Equal(t, zoneID, resp.Zone.Id)
		assert.Equal(t, "Uptown", resp.Zone.Name)
		assert.True(t, resp.Zone.IsActive)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Get(ctx, zoneID).Return(model.Zone{}, model.ErrZoneNotFound)

		_, err := h.GetZone(ctx, &carsvc.GetZoneRequest{Id: zoneID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestZoneHandlerListZones(t *testing.T) {
	ctx := context.Background()

	t.Run("returns zone list", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().List(ctx, mock.Anything).Return([]model.Zone{
			{ID: "z-1", Name: "North"},
			{ID: "z-2", Name: "South"},
		}, nil)

		resp, err := h.ListZones(ctx, &carsvc.ListZonesRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Zones, 2)
		assert.Equal(t, "z-1", resp.Zones[0].Id)
	})

	t.Run("service error maps to gRPC Internal", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().List(ctx, mock.Anything).Return(nil, errInternal)

		_, err := h.ListZones(ctx, &carsvc.ListZonesRequest{})
		assert.Equal(t, codes.Internal, grpcCode(err))
	})
}

func TestZoneHandlerUpdateZone(t *testing.T) {
	ctx := context.Background()
	zoneID := "zone-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Update(ctx, zoneID, mock.Anything).Return(nil)

		resp, err := h.UpdateZone(ctx, &carsvc.UpdateZoneRequest{Id: zoneID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Update(ctx, zoneID, mock.Anything).Return(model.ErrZoneNotFound)

		_, err := h.UpdateZone(ctx, &carsvc.UpdateZoneRequest{Id: zoneID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}

func TestZoneHandlerDeleteZone(t *testing.T) {
	ctx := context.Background()
	zoneID := "zone-123"

	t.Run("returns empty on success", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Delete(ctx, zoneID).Return(nil)

		resp, err := h.DeleteZone(ctx, &carsvc.DeleteZoneRequest{Id: zoneID})
		assert.NoError(t, err)
		assert.IsType(t, &emptypb.Empty{}, resp)
	})

	t.Run("not found maps to gRPC NotFound", func(t *testing.T) {
		svc := mocks.NewMockZoneService(t)
		h := NewZoneHandler(discardLogger(), svc)

		svc.EXPECT().Delete(ctx, zoneID).Return(model.ErrZoneNotFound)

		_, err := h.DeleteZone(ctx, &carsvc.DeleteZoneRequest{Id: zoneID})
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})
}
