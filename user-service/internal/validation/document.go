package validation

type DocumentCreate struct {
	ObjectKey string `validate:"required"`
	ImageType string `validate:"required,imagetype"`
}

type DocumentUpdate struct {
	Status string `validate:"required,documentstatus"`
	Error  *string
}
