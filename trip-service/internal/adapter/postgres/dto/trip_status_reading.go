package dto

import (
	"fmt"
	"time"

	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/model"
)

type statusReadingRow struct {
	ID         string
	TripID     string
	FromStatus string
	ToStatus   string
	ActorType  string
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

func ScanTripStatusReading(s scanner) (model.TripStatusReading, error) {
	var r statusReadingRow
	err := s.Scan(
		&r.ID, &r.TripID, &r.FromStatus, &r.ToStatus,
		&r.ActorType, &r.ActorID, &r.Reason, &r.ChangedAt,
	)
	if err != nil {
		return model.TripStatusReading{}, err
	}

	return model.TripStatusReading{
		ID:         r.ID,
		TripID:     r.TripID,
		FromStatus: model.TripStatus(r.FromStatus),
		ToStatus:   model.TripStatus(r.ToStatus),
		ActorType:  sharedmodel.ActorType(r.ActorType),
		ActorID:    r.ActorID,
		Reason:     r.Reason,
		ChangedAt:  r.ChangedAt,
	}, nil
}

func BuildStatusReadingWhereClauses(f model.TripStatusReadingFilter, b *ArgsBuilder) []string {
	clauses := []string{fmt.Sprintf("trip_id = %s", b.Add(f.TripID))}
	if f.TimeRange != nil {
		clauses = append(clauses, fmt.Sprintf("changed_at >= %s", b.Add(f.TimeRange.From)))
		clauses = append(clauses, fmt.Sprintf("changed_at <= %s", b.Add(f.TimeRange.To)))
	}
	return clauses
}
