package dto

import (
	"database/sql"
	"fmt"
	"time"

	"carsharing/trip-service/internal/model"
)

type statusReadingRow struct {
	ID         string
	TripID     string
	FromStatus string
	ToStatus   string
	ActorType  string
	ActorID    sql.NullString
	Reason     sql.NullString
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

	reading := model.TripStatusReading{
		ID:         r.ID,
		TripID:     r.TripID,
		FromStatus: model.TripStatus(r.FromStatus),
		ToStatus:   model.TripStatus(r.ToStatus),
		ActorType:  model.ActorType(r.ActorType),
		ChangedAt:  r.ChangedAt,
	}
	if r.ActorID.Valid {
		reading.ActorID = &r.ActorID.String
	}
	if r.Reason.Valid {
		reading.Reason = &r.Reason.String
	}
	return reading, nil
}

func BuildStatusReadingWhereClauses(f model.TripStatusReadingFilter, b *ArgsBuilder) []string {
	clauses := []string{fmt.Sprintf("trip_id = %s", b.Add(f.TripID))}
	if f.From != nil {
		clauses = append(clauses, fmt.Sprintf("changed_at >= %s", b.Add(*f.From)))
	}
	if f.To != nil {
		clauses = append(clauses, fmt.Sprintf("changed_at <= %s", b.Add(*f.To)))
	}
	return clauses
}
