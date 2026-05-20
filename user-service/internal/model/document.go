package model

import "time"

type DocumentStatus string

const (
	DocumentStatusPending   DocumentStatus = "pending"
	DocumentStatusProcessed DocumentStatus = "processed"
	DocumentStatusApproved  DocumentStatus = "approved"
	DocumentStatusRejected  DocumentStatus = "rejected"
)

var validDocumentStatuses = map[DocumentStatus]struct{}{
	DocumentStatusPending:   {},
	DocumentStatusProcessed: {},
	DocumentStatusApproved:  {},
	DocumentStatusRejected:  {},
}

func DocumentStatusFromString(s string) (DocumentStatus, error) {
	ds := DocumentStatus(s)
	if _, ok := validDocumentStatuses[ds]; !ok {
		return "", ErrInvalidDocumentStatus
	}
	return ds, nil
}

func (s DocumentStatus) String() string {
	return string(s)
}

type Document struct {
	ID        string
	UserID    string
	ImageType ImageType
	Status    DocumentStatus
	Error     *string
	Image     *Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentFilter struct {
	UserID        *string
	ExcludeStatus *DocumentStatus
	LatestPerType bool
}

type DocumentUpdate struct {
	Status    *DocumentStatus
	Error     *string
	UpdatedAt time.Time
}
