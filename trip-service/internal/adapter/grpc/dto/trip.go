package dto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	basepb "carsharing/protos/gen/base"
	basetripmb "carsharing/protos/gen/base/trip"
	tripsvc "carsharing/protos/gen/service/trip"

	sharedvalidation "carsharing/shared/validation"
	"carsharing/trip-service/internal/model"
	"carsharing/trip-service/internal/validation"
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

func FilterFromProto(req *tripsvc.ListTripsRequest) validation.TripFilter {
	f := validation.TripFilter{
		UserID: req.UserId,
		CarID:  req.CarId,
	}
	if req.Status != nil {
		s := *req.Status
		f.Status = &s
	}
	if req.StartedAfter != nil || req.StartedBefore != nil {
		tr := &sharedvalidation.TimeRange{}
		if req.StartedAfter != nil {
			t := req.StartedAfter.AsTime()
			tr.From = &t
		}
		if req.StartedBefore != nil {
			t := req.StartedBefore.AsTime()
			tr.To = &t
		}
		f.TimeRange = tr
	}
	if req.Pagination != nil {
		f.Pagination = &sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return f
}

func StatusHistoryFilterFromProto(req *tripsvc.GetTripStatusHistoryRequest) validation.TripStatusHistoryFilter {
	f := validation.TripStatusHistoryFilter{TripID: req.Id}
	if req.From != nil || req.To != nil {
		tr := &sharedvalidation.TimeRange{}
		if req.From != nil {
			t := req.From.AsTime()
			tr.From = &t
		}
		if req.To != nil {
			t := req.To.AsTime()
			tr.To = &t
		}
		f.TimeRange = tr
	}
	if req.Pagination != nil {
		f.Pagination = &sharedvalidation.Pagination{
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
