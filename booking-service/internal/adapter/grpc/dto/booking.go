package dto

import (
	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/validation"
	basepb "carsharing/protos/gen/base"
	basebookingpb "carsharing/protos/gen/base/booking"
	servicebookingpb "carsharing/protos/gen/service/booking"
	sharedvalidation "carsharing/shared/validation"
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
		ActorType:  string(r.ActorType),
		ActorId:    r.ActorID,
		Reason:     r.Reason,
		ChangedAt:  timestamppb.New(r.ChangedAt),
	}
}

func BookingCreateFromProto(req *servicebookingpb.CreateBookingRequest) validation.BookingCreate {
	return validation.BookingCreate{
		UserID:           req.UserId,
		CarID:            req.CarId,
		PricingRuleID:    req.PricingRuleId,
		CommittedPeriods: req.CommittedPeriods,
	}
}

func BookingStatusUpdateFromProto(req *servicebookingpb.UpdateBookingStatusRequest) validation.BookingStatusUpdate {
	return validation.BookingStatusUpdate{
		Status: req.Status,
		Reason: req.Reason,
	}
}

func BookingListFilterFromProto(req *servicebookingpb.ListBookingsRequest) validation.BookingListFilter {
	filter := validation.BookingListFilter{
		UserID:        req.UserId,
		CarID:         req.CarId,
		Status:        req.Status,
		PricingRuleID: req.PricingRuleId,
	}
	if req.Pagination != nil {
		p := sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
		filter.Pagination = &p
	}
	return filter
}

func BookingStatusHistoryFilterFromProto(req *servicebookingpb.GetBookingStatusHistoryRequest) validation.BookingStatusHistoryFilter {
	filter := validation.BookingStatusHistoryFilter{
		BookingID: req.Id,
	}
	if req.GetTimeRange() != nil {
		tr := &sharedvalidation.TimeRange{}
		if req.GetTimeRange().GetFrom() != nil {
			t := req.GetTimeRange().GetFrom().AsTime()
			tr.From = &t
		}
		if req.GetTimeRange().GetTo() != nil {
			t := req.GetTimeRange().GetTo().AsTime()
			tr.To = &t
		}
		filter.TimeRange = tr
	}
	if req.Pagination != nil {
		p := sharedvalidation.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
		filter.Pagination = &p
	}
	return filter
}
