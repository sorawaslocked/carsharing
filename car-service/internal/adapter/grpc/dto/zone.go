package dto

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/model"

	basecar "github.com/sorawaslocked/car-rental-protos/gen/base/car"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateZoneRequest(req *carsvc.CreateZoneRequest) model.ZoneCreateInput {
	return model.ZoneCreateInput{
		Name:            req.Name,
		Type:            req.Type,
		BoundaryGeoJSON: req.BoundaryGeoJson,
		FeeAdjustment:   req.FeeAdjustment,
	}
}

func FromListZonesRequest(req *carsvc.ListZonesRequest) model.ZoneFilterInput {
	filter := model.ZoneFilterInput{
		Type:     req.Type,
		IsActive: req.IsActive,
	}
	if req.Pagination != nil {
		limit := req.Pagination.Limit
		offset := req.Pagination.Offset
		filter.PaginationInput = model.PaginationInput{
			Limit:  &limit,
			Offset: &offset,
		}
	}
	return filter
}

func FromUpdateZoneRequest(req *carsvc.UpdateZoneRequest) model.ZoneUpdateInput {
	return model.ZoneUpdateInput{
		Name:            req.Name,
		Type:            req.Type,
		BoundaryGeoJSON: req.BoundaryGeoJson,
		FeeAdjustment:   req.FeeAdjustment,
		IsActive:        req.IsActive,
	}
}

func ToZoneProto(z model.Zone) *basecar.Zone {
	return &basecar.Zone{
		Id:              z.ID,
		Name:            z.Name,
		Type:            string(z.Type),
		BoundaryGeoJson: z.BoundaryGeoJSON,
		FeeAdjustment:   z.FeeAdjustment,
		IsActive:        z.IsActive,
		CreatedAt:       timestamppb.New(z.CreatedAt),
		UpdatedAt:       timestamppb.New(z.UpdatedAt),
	}
}

func ToZoneProtos(zones []model.Zone) []*basecar.Zone {
	protos := make([]*basecar.Zone, len(zones))
	for i, z := range zones {
		protos[i] = ToZoneProto(z)
	}
	return protos
}
