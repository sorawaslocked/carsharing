package validation

type DocumentCreate struct {
	ObjectKey string `validate:"required"`
	ImageType string `validate:"required,document_image_type"`
}

type DocumentUpdate struct {
	Status string `validate:"required,document_status"`
	Error  *string
}
