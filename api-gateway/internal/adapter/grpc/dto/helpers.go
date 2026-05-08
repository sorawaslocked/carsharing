package dto

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
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

func LocationFromProto(l *basepb.Location) model.Location {
	if l == nil {
		return model.Location{}
	}
	return model.Location{
		Latitude:  l.GetLatitude(),
		Longitude: l.GetLongitude(),
	}
}

func LocationToProto(l model.Location) *basepb.Location {
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

func ImageUploadDataFromProto(u *basepb.ImageUploadData) model.ImageUploadData {
	if u == nil {
		return model.ImageUploadData{}
	}
	return model.ImageUploadData{
		PresignedPutURL: u.GetPresignedPutUrl(),
		ObjectKey:       u.GetObjectKey(),
	}
}
