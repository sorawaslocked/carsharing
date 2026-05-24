package dto

import (
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	sharedmodel "carsharing/shared/model"
	"google.golang.org/protobuf/types/known/structpb"
)

func PricingSnapshotFromProto(s *basepb.PricingSnapshot) model.PricingSnapshot {
	return model.PricingSnapshot{
		RateTenge:         s.GetRateTenge(),
		RatePerKMTenge:    s.RatePerKmTenge,
		FreeMinutes:       s.FreeMinutes,
		MinChargeTenge:    s.MinChargeTenge,
		OvertimePolicy:    s.OvertimePolicy,
		OvertimeRateTenge: s.OvertimeRateTenge,
	}
}

func LocationFromProto(l *basepb.Location) sharedmodel.Location {
	if l == nil {
		return sharedmodel.Location{}
	}
	return sharedmodel.Location{
		Latitude:  l.GetLatitude(),
		Longitude: l.GetLongitude(),
	}
}

func LocationToProto(l sharedmodel.Location) *basepb.Location {
	return &basepb.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	}
}

func structToMap(s *structpb.Struct) map[string]any {
	if s == nil {
		return nil
	}

	return s.AsMap()
}

func ImageUploadDataFromProto(u *basepb.ImageUploadData) sharedmodel.ImageUploadData {
	if u == nil {
		return sharedmodel.ImageUploadData{}
	}
	return sharedmodel.ImageUploadData{
		PresignedPutURL: u.GetPresignedPutUrl(),
		ObjectKey:       u.GetObjectKey(),
	}
}
