package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type Document struct {
	ID        string
	UserID    string
	ImageType string
	Status    string
	Reason    *string
	Image     sharedmodel.Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentFilter struct {
	UserID     string
	Status     *string
	ImageType  *string
	Sort       *string
	Pagination *sharedmodel.Pagination
}

type DocumentAnalyzedEvent struct {
	DocumentID string
	UserID     string
	Passed     bool
	Defects    []DocumentDefect
}

type DocumentDefect struct {
	Type        string
	Description string
}
