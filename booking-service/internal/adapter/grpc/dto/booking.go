package dto

import (
	"carsharing/booking-service/internal/model"
	basepb "carsharing/protos/gen/base"
	basebookingpb "carsharing/protos/gen/base/booking"
	servicebookingpb "carsharing/protos/gen/service/booking"
	sharedmodel "carsharing/shared/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BookingToProto(b model.Booking) *basebookingpb.Booking {
	pb := &basebookingpb.Booking{
		Id:              b.ID,
		UserId:          b.UserID,
		CarId:           b.CarID,
		Status:          string(b.Status),
		PricingRuleId:   b.PricingRuleID,
		PricingSnapshot: PricingSnapshotToProto(b.PricingSnapshot),
		ExpiresAt:       timestamppb.New(b.ExpiresAt),
		CreatedAt:       timestamppb.New(b.CreatedAt),
		UpdatedAt:       timestamppb.New(b.UpdatedAt),
	}
	if b.CommittedPeriods != nil {
		pb.CommittedPeriods = b.CommittedPeriods
	}
	return pb
}

func PricingSnapshotToProto(s model.PricingSnapshot) *basepb.PricingSnapshot {
	return &basepb.PricingSnapshot{
		RateTenge:         s.RateTenge,
		RatePerKmTenge:    s.RatePerKMTenge,
		FreeMinutes:       s.FreeMinutes,
		MinChargeTenge:    s.MinChargeTenge,
		OvertimePolicy:    s.OvertimePolicy,
		OvertimeRateTenge: s.OvertimeRateTenge,
	}
}

func BookingStatusReadingToProto(r model.BookingStatusReading) *basebookingpb.BookingStatusReading {
	return &basebookingpb.BookingStatusReading{
		Id:         r.ID,
		BookingId:  r.BookingID,
		FromStatus: r.FromStatus,
		ToStatus:   r.ToStatus,
		ActorType:  r.ActorType,
		ActorId:    r.ActorID,
		Reason:     r.Reason,
		ChangedAt:  timestamppb.New(r.ChangedAt),
	}
}

func BookingListFilterFromProto(req *servicebookingpb.ListBookingsRequest) model.BookingListFilter {
	filter := model.BookingListFilter{
		UserID:        req.UserId,
		CarID:         req.CarId,
		Status:        req.Status,
		PricingRuleID: req.PricingRuleId,
	}
	if req.Pagination != nil {
		filter.Pagination = sharedmodel.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}

func BookingStatusHistoryFilterFromProto(req *servicebookingpb.GetBookingStatusHistoryRequest) model.BookingStatusHistoryFilter {
	filter := model.BookingStatusHistoryFilter{
		BookingID: req.Id,
	}
	if req.From != nil {
		t := req.From.AsTime()
		filter.From = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		filter.To = &t
	}
	if req.Pagination != nil {
		filter.Pagination = sharedmodel.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}
	return filter
}
