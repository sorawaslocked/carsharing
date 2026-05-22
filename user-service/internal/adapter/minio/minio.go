package minio

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"

	"github.com/minio/minio-go/v7"
)

type ObjectStorage struct {
	log    *slog.Logger
	client *minio.Client
	cfg    pkgminio.Config
}

func NewObjectStorage(log *slog.Logger, client *minio.Client, cfg pkgminio.Config) *ObjectStorage {
	return &ObjectStorage{
		log:    pkglog.WithComponent(log, "ObjectStorage"),
		client: client,
		cfg:    cfg,
	}
}

func (s *ObjectStorage) GetDocumentImageUploadData(ctx context.Context, imageType string) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetDocumentImageUploadData"), utils.MetadataFromCtx(ctx))

	objectKey := newObjectKey("documents/" + imageType)

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

func (s *ObjectStorage) GetUserProfileImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetUserProfileImageUploadData"), utils.MetadataFromCtx(ctx))

	objectKey := newObjectKey("users")

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

func (s *ObjectStorage) GetImageURL(ctx context.Context, objectKey string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageURL"), utils.MetadataFromCtx(ctx))

	presignedURL, err := s.client.PresignedGetObject(ctx, s.cfg.Bucket, objectKey, s.cfg.PresignedGetExpiry, url.Values{})
	if err != nil {
		log.Error("generating presigned get url", pkglog.Err(err))

		return "", model.ErrObjectStorage
	}

	return presignedURL.String(), nil
}

// newObjectKey returns "{prefix}/{unix_timestamp}_{random_hex}".
func newObjectKey(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}

	return fmt.Sprintf("%s/%d_%s", prefix, time.Now().UnixNano(), hex.EncodeToString(b))
}
