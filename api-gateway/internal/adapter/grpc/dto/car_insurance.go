package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "github.com/sorawaslocked/car-rental-protos/gen/base/car"
)

func CarInsuranceFromProto(i *basecarpb.CarInsurance) model.CarInsurance {
	ins := model.CarInsurance{
		ID:        i.GetId(),
		CarID:     i.GetCarId(),
		Type:      i.GetType(),
		Provider:  i.GetProvider(),
		PolicyNum: i.GetPolicyNum(),
		CostTenge: i.GetCostTenge(),
		Status:    i.GetStatus(),
		ImageURLs: i.GetImageUrls(),
		Notes:     i.Notes,
	}
	if i.GetStartsAt() != nil {
		ins.StartsAt = i.GetStartsAt().AsTime()
	}
	if i.GetExpiresAt() != nil {
		ins.ExpiresAt = i.GetExpiresAt().AsTime()
	}
	if i.GetCreatedAt() != nil {
		ins.CreatedAt = i.GetCreatedAt().AsTime()
	}
	if i.GetUpdatedAt() != nil {
		ins.UpdatedAt = i.GetUpdatedAt().AsTime()
	}
	return ins
}
