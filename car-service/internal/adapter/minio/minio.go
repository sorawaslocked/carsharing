package minio

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	"carsharing/shared/pkg/utils"

	"github.com/minio/minio-go/v7"
)

type ObjectStorage struct {
	log    *slog.Logger
	client *minio.Client
	cfg    pkgminio.Config
}

func NewObjectStorage(log *slog.Logger, client *minio.Client, cfg pkgminio.Config) (*ObjectStorage, error) {
	log = pkglog.WithComponent(log, "ObjectStorage")

	return &ObjectStorage{
		log:    log,
		client: client,
		cfg:    cfg,
	}, nil
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
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "getImageUploadData"), utils.MetadataFromCtx(ctx))

	objectKey := newObjectKey(prefix)

	presignedURL, err := s.client.PresignedPutObject(ctx, s.cfg.Bucket, objectKey, s.cfg.PresignedPutExpiry)
	if err != nil {
		log.Error("generating presigned put url", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, model.ErrObjectStorage
	}

	return sharedmodel.ImageUploadData{
		PresignedPutURL: presignedURL.String(),
		ObjectKey:       objectKey,
	}, nil
}

func (s *ObjectStorage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetPresignedURL"), utils.MetadataFromCtx(ctx))

	presignedURL, err := s.client.PresignedGetObject(ctx, s.cfg.Bucket, key, s.cfg.PresignedGetExpiry, url.Values{})
	if err != nil {
		log.Error("generating presigned get url", pkglog.Err(err))

		return "", model.ErrObjectStorage
	}

	return presignedURL.String(), nil
}

func newObjectKey(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}

	return fmt.Sprintf("%s/%d_%s", prefix, time.Now().UnixNano(), hex.EncodeToString(b))
}
