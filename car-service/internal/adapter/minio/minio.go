package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
)

const presignedURLTTL = 15 * time.Minute

type ObjectStorage struct {
	client *minio.Client
	bucket string
}

func NewObjectStorage(client *minio.Client, bucket string) *ObjectStorage {
	return &ObjectStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *ObjectStorage) GetImageUploadData(ctx context.Context, prefix string) (model.ImageUploadData, error) {
	key := fmt.Sprintf("%s/uploads/%d", prefix, time.Now().UnixNano())

	u, err := s.client.PresignedPutObject(ctx, s.bucket, key, presignedURLTTL)
	if err != nil {
		return model.ImageUploadData{}, fmt.Errorf("presign put object: %w", err)
	}

	return model.ImageUploadData{
		ObjectKey: key,
		URL:       u.String(),
	}, nil
}

func (s *ObjectStorage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, presignedURLTTL, url.Values{})
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}

	return u.String(), nil
}
