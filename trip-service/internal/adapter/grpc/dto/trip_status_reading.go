package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	basetripmb "carsharing/protos/gen/base/trip"

	"carsharing/trip-service/internal/model"
)

func TripStatusReadingToProto(r model.TripStatusReading) *basetripmb.TripStatusReading {
	return &basetripmb.TripStatusReading{
		Id:         r.ID,
		TripId:     r.TripID,
		FromStatus: string(r.FromStatus),
		ToStatus:   string(r.ToStatus),
		ActorType:  string(r.ActorType),
		ActorId:    r.ActorID,
		Reason:     r.Reason,
		ChangedAt:  timestamppb.New(r.ChangedAt),
	}
}
