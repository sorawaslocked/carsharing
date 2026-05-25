package model

import (
	sharedmodel "carsharing/shared/model"
	"time"
)

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

func DocumentStatusFromString(s string) (DocumentStatus, bool) {
	ds := DocumentStatus(s)
	if _, ok := validDocumentStatuses[ds]; !ok {
		return "", false
	}
	return ds, true
}

func (s DocumentStatus) String() string {
	return string(s)
}

type DocumentImageType string

const (
	DocumentImageTypeIDFront             DocumentImageType = "id_front"
	DocumentImageTypeIDBack              DocumentImageType = "id_back"
	DocumentImageTypeDrivingLicenseFront DocumentImageType = "driving_license_front"
	DocumentImageTypeDrivingLicenseBack  DocumentImageType = "driving_license_back"
)

var validDocumentImageTypes = map[DocumentImageType]struct{}{
	DocumentImageTypeIDFront:             {},
	DocumentImageTypeIDBack:              {},
	DocumentImageTypeDrivingLicenseFront: {},
	DocumentImageTypeDrivingLicenseBack:  {},
}

func AllDocumentImageTypes() []DocumentImageType {
	return []DocumentImageType{
		DocumentImageTypeIDFront,
		DocumentImageTypeIDBack,
		DocumentImageTypeDrivingLicenseFront,
		DocumentImageTypeDrivingLicenseBack,
	}
}

func DocumentImageTypeFromString(s string) (DocumentImageType, bool) {
	it := DocumentImageType(s)
	if _, ok := validDocumentImageTypes[it]; ok {
		return it, true
	}
	return "", false
}

func (t DocumentImageType) String() string {
	return string(t)
}

type Document struct {
	ID        string
	UserID    string
	ImageType DocumentImageType
	Status    DocumentStatus
	Error     *string
	Image     sharedmodel.Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentSort string

const (
	DocumentSortCreatedAtAsc  DocumentSort = "+createdAt"
	DocumentSortCreatedAtDesc DocumentSort = "-createdAt"
)

var validDocumentSorts = map[DocumentSort]struct{}{
	DocumentSortCreatedAtAsc:  {},
	DocumentSortCreatedAtDesc: {},
}

func DocumentSortFromString(s string) (DocumentSort, bool) {
	ds := DocumentSort(s)
	if _, ok := validDocumentSorts[ds]; !ok {
		return "", false
	}
	return ds, true
}

func (s DocumentSort) String() string {
	return string(s)
}

type DocumentFilter struct {
	UserID     string
	Status     *DocumentStatus
	ImageType  *DocumentImageType
	Sort       *DocumentSort
	Pagination *sharedmodel.Pagination

	LatestPerType bool
}

type DocumentUpdate struct {
	Status    *DocumentStatus
	Error     *string
	UpdatedAt time.Time
}
