package validation

import sharedvalidation "carsharing/shared/validation"

type DocumentCreate struct {
	ObjectKey string `validate:"required"`
	ImageType string `validate:"required,document_image_type"`
}

type DocumentUpdate struct {
	Status string `validate:"required,document_status"`
	Error  *string
}

type DocumentFilter struct {
	UserID     string  `validate:"required,uuid4"`
	Status     *string `validate:"omitempty,document_status"`
	ImageType  *string `validate:"omitempty,document_image_type"`
	Sort       *string `validate:"omitempty,document_sort"`
	Pagination *sharedvalidation.Pagination
}
