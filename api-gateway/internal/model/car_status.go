package model

import "time"

type CarStatusLogEntry struct {
	ID         string
	CarID      string
	FromStatus string
	ToStatus   string
	ActorType  string
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	ChangedAt  time.Time
}

type CarStatusLogFilter struct {
	CarID      *string
	Pagination *Pagination
}
