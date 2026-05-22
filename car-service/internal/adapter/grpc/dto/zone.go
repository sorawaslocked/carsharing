package dto

import (
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"

	basecar "github.com/sorawaslocked/car-rental-protos/gen/base/car"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateZoneRequest(req *carsvc.CreateZoneRequest) validation.ZoneCreate {
	return validation.ZoneCreate{
		Name:            req.Name,
		Type:            req.Type,
		BoundaryGeoJSON: req.BoundaryGeoJson,
		FeeAdjustment:   req.FeeAdjustment,
	}
}

func FromListZonesRequest(req *carsvc.ListZonesRequest) validation.ZoneFilter {
	filter := validation.ZoneFilter{
		Type:     req.Type,
		IsActive: req.IsActive,
	}
	if req.Pagination != nil {
		filter.Pagination = &sharedmodel.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func FromUpdateZoneRequest(req *carsvc.UpdateZoneRequest) validation.ZoneUpdate {
	return validation.ZoneUpdate{
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
