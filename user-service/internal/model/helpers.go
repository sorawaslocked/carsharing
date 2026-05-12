package model

// Image holds the object storage key persisted in the DB and the
// pre-signed URL resolved on demand. URL is empty until resolved.
type Image struct {
	Key string
	URL string
}
