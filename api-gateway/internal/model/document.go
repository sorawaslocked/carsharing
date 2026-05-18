package model

import "time"

type Document struct {
	ID        string
	UserID    string
	ImageType string
	Status    string
	Reason    *string
	ImageURL  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentAnalyzedEvent struct {
	DocumentID string
	Passed     bool
	Defects    []DocumentDefect
}

type DocumentDefect struct {
	Type        string
	Description string
}
