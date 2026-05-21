package model

type Image struct {
	Key string
	URL string
}

type ImageUploadData struct {
	PresignedPutURL string
	ObjectKey       string
}
