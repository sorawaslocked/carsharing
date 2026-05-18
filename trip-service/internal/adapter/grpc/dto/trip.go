package dto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	basetripmb "github.com/sorawaslocked/car-rental-protos/gen/base/trip"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"

	"carsharing/trip-service/internal/model"
)

func TripToProto(t model.Trip) *basetripmb.Trip {
	proto := &basetripmb.Trip{
		Id:                 t.ID,
		BookingId:          t.BookingID,
		UserId:             t.UserID,
		CarId:              t.CarID,
		Status:             string(t.Status),
		StartedAt:          timestamppb.New(t.StartedAt),
		StartLocation:      &basepb.Location{Latitude: t.StartLocation.Latitude, Longitude: t.StartLocation.Longitude},
		StartMileageKm:     t.StartMileageKM,
		StartFuelLevel:     t.StartFuelLevel,
		EndedAt:            timeToTimestamp(t.EndedAt),
		EndMileageKm:       t.EndMileageKM,
		EndFuelLevel:       t.EndFuelLevel,
		DistanceTraveledKm: t.DistanceTraveledKM,
		DurationSeconds:    t.DurationSeconds,
		FinalCostTenge:     t.FinalCostTenge,
		CancelReason:       t.CancelReason,
		CreatedAt:          timestamppb.New(t.CreatedAt),
		UpdatedAt:          timestamppb.New(t.UpdatedAt),
	}
	if t.EndLocation != nil {
		proto.EndLocation = &basepb.Location{
			Latitude:  t.EndLocation.Latitude,
			Longitude: t.EndLocation.Longitude,
		}
	}
	return proto
}

func FilterFromProto(req *tripsvc.ListTripsRequest) model.TripFilter {
	f := model.TripFilter{
		UserID: req.UserId,
		CarID:  req.CarId,
	}
	if req.Status != nil {
		s := model.TripStatus(*req.Status)
		f.Status = &s
	}
	if req.StartedAfter != nil {
		t := req.StartedAfter.AsTime()
		f.StartedAfter = &t
	}
	if req.StartedBefore != nil {
		t := req.StartedBefore.AsTime()
		f.StartedBefore = &t
	}
	if req.Pagination != nil {
		f.Pagination = &model.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return f
}

func StatusHistoryFilterFromProto(req *tripsvc.GetTripStatusHistoryRequest) model.TripStatusReadingFilter {
	f := model.TripStatusReadingFilter{TripID: req.Id}
	if req.From != nil {
		t := req.From.AsTime()
		f.From = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		f.To = &t
	}
	if req.Pagination != nil {
		f.Pagination = &model.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return f
}

func timeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
