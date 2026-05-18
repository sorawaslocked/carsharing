package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "github.com/sorawaslocked/car-rental-protos/gen/base/car"
)

func ZoneFromProto(z *basecarpb.Zone) model.Zone {
	zone := model.Zone{
		ID:              z.GetId(),
		Name:            z.GetName(),
		Type:            z.GetType(),
		BoundaryGeoJSON: z.GetBoundaryGeoJson(),
		FeeAdjustment:   z.GetFeeAdjustment(),
		IsActive:        z.GetIsActive(),
	}
	if z.GetCreatedAt() != nil {
		zone.CreatedAt = z.GetCreatedAt().AsTime()
	}
	if z.GetUpdatedAt() != nil {
		zone.UpdatedAt = z.GetUpdatedAt().AsTime()
	}
	return zone
}
