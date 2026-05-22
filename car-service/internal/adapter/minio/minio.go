package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	sharedmodel "carsharing/shared/model"
	pkgminio "carsharing/shared/pkg/minio"
	"github.com/minio/minio-go/v7"
)

const presignedURLTTL = 15 * time.Minute

type ObjectStorage struct {
	client *minio.Client
	cfg    pkgminio.Config
}

func NewObjectStorage(client *minio.Client, cfg pkgminio.Config) *ObjectStorage {
	return &ObjectStorage{
		client: client,
		cfg:    cfg,
	}
}

func (s *ObjectStorage) GetCarImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "cars")
}

func (s *ObjectStorage) GetCarModelImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "car-models")
}

func (s *ObjectStorage) GetInsuranceImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "insurance")
}

func (s *ObjectStorage) GetMaintenanceReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	return s.getImageUploadData(ctx, "maintenance/receipts")
}

func (s *ObjectStorage) getImageUploadData(ctx context.Context, prefix string) (sharedmodel.ImageUploadData, error) {
	key := fmt.Sprintf("%s/uploads/%d", prefix, time.Now().UnixNano())

	u, err := s.client.PresignedPutObject(ctx, s.cfg.Bucket, key, presignedURLTTL)
	if err != nil {
		return sharedmodel.ImageUploadData{}, fmt.Errorf("presign put object: %w", err)
	}

	return sharedmodel.ImageUploadData{
		ObjectKey:       key,
		PresignedPutURL: s.publicURL(u).String(),
	}, nil
}

func (s *ObjectStorage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.cfg.Bucket, key, presignedURLTTL, url.Values{})
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}

	return s.publicURL(u).String(), nil
}

func (s *ObjectStorage) publicURL(u *url.URL) *url.URL {
	if s.cfg.PublicEndpoint == "" {
		return u
	}
	rewritten := *u
	rewritten.Host = s.cfg.PublicEndpoint
	return &rewritten
}
