package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"carsharing/car-service/internal/model"
	"github.com/minio/minio-go/v7"
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

func (s *ObjectStorage) GetCarImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "cars")
}

func (s *ObjectStorage) GetCarModelImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "car-models")
}

func (s *ObjectStorage) GetInsuranceImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "insurance")
}

func (s *ObjectStorage) GetMaintenanceReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "maintenance/receipts")
}

func (s *ObjectStorage) getImageUploadData(ctx context.Context, prefix string) (model.ImageUploadData, error) {
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
