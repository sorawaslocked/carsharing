package model

type Pagination struct {
	Limit  int64
	Offset int64
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type ImageUploadData struct {
	PresignedUrl string
	ObjectKey    string
}
