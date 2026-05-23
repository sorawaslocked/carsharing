package dto

import (
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedvalidation "carsharing/shared/validation"

	basecar "carsharing/protos/gen/base/car"
	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateCarInsuranceRequest(req *carsvc.CreateCarInsuranceRequest) validation.CarInsuranceCreate {
	return validation.CarInsuranceCreate{
		CarID:     req.CarId,
		Type:      req.Type,
		Provider:  req.Provider,
		PolicyNum: req.PolicyNum,
		StartsAt:  req.StartsAt.AsTime(),
		ExpiresAt: req.ExpiresAt.AsTime(),
		CostTenge: req.CostTenge,
		Notes:     req.Notes,
	}
}

func FromListCarInsurancesRequest(req *carsvc.ListCarInsurancesRequest) validation.CarInsuranceFilter {
	filter := validation.CarInsuranceFilter{
		CarID:              req.CarId,
		Type:               req.Type,
		Status:             req.Status,
		ExpiringWithinDays: req.ExpiringWithinDays,
	}
	if req.Pagination != nil {
		filter.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func FromUpdateCarInsuranceRequest(req *carsvc.UpdateCarInsuranceRequest) validation.CarInsuranceUpdate {
	update := validation.CarInsuranceUpdate{
		Provider:  req.Provider,
		PolicyNum: req.PolicyNum,
		CostTenge: req.CostTenge,
		Status:    req.Status,
		Notes:     req.Notes,
		ImageKeys: req.ImageKeys,
	}
	if req.StartsAt != nil {
		t := req.StartsAt.AsTime()
		update.StartsAt = &t
	}
	if req.ExpiresAt != nil {
		t := req.ExpiresAt.AsTime()
		update.ExpiresAt = &t
	}
	return update
}

func ToCarInsuranceProto(ins model.CarInsurance) *basecar.CarInsurance {
	proto := &basecar.CarInsurance{
		Id:        ins.ID,
		CarId:     ins.CarID,
		Type:      string(ins.Type),
		Provider:  ins.Provider,
		PolicyNum: ins.PolicyNum,
		StartsAt:  timestamppb.New(ins.StartsAt),
		ExpiresAt: timestamppb.New(ins.ExpiresAt),
		CostTenge: ins.CostTenge,
		Status:    string(ins.Status),
		ImageUrls: imageURLsFromImages(ins.Images),
		CreatedAt: timestamppb.New(ins.CreatedAt),
		UpdatedAt: timestamppb.New(ins.UpdatedAt),
	}
	proto.Notes = ins.Notes
	return proto
}

func ToCarInsuranceProtos(insurances []model.CarInsurance) []*basecar.CarInsurance {
	protos := make([]*basecar.CarInsurance, len(insurances))
	for i, ins := range insurances {
		protos[i] = ToCarInsuranceProto(ins)
	}
	return protos
}
