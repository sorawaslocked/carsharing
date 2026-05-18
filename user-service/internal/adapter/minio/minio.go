package minio

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"carsharing/user-service/internal/model"
	pkglog "carsharing/user-service/internal/pkg/log"
	miniocfg "carsharing/user-service/internal/pkg/minio"
	"carsharing/user-service/internal/pkg/utils"
	"github.com/minio/minio-go/v7"
)

type MinioObjectStorage struct {
	log    *slog.Logger
	client *minio.Client
	cfg    miniocfg.Config
}

func NewMinioObjectStorage(log *slog.Logger, client *minio.Client, cfg miniocfg.Config) *MinioObjectStorage {
	return &MinioObjectStorage{
		log:    pkglog.WithComponent(log, "MinioObjectStorage"),
		client: client,
		cfg:    cfg,
	}
}

func (s *MinioObjectStorage) GetDocumentImageUploadData(ctx context.Context, imageType string) (model.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetDocumentImageUploadData"), utils.MetadataFromCtx(ctx))

	objectKey := newObjectKey("documents/" + imageType)

	presignedURL, err := s.client.PresignedPutObject(ctx, s.cfg.BucketName, objectKey, s.cfg.PresignedPutExpiry)
	if err != nil {
		log.Error("generating presigned put url", pkglog.Err(err))
		return model.ImageUploadData{}, model.ErrObjectStorage
	}

	return model.ImageUploadData{
		PresignedPutURL: presignedURL.String(),
		ObjectKey:       objectKey,
	}, nil
}

func (s *MinioObjectStorage) GetUserProfileImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetUserProfileImageUploadData"), utils.MetadataFromCtx(ctx))

	objectKey := newObjectKey("users")

	presignedURL, err := s.client.PresignedPutObject(ctx, s.cfg.BucketName, objectKey, s.cfg.PresignedPutExpiry)
	if err != nil {
		log.Error("generating presigned put url", pkglog.Err(err))
		return model.ImageUploadData{}, model.ErrObjectStorage
	}

	return model.ImageUploadData{
		PresignedPutURL: presignedURL.String(),
		ObjectKey:       objectKey,
	}, nil
}

func (s *MinioObjectStorage) GetImageURL(ctx context.Context, objectKey string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageURL"), utils.MetadataFromCtx(ctx))

	presignedURL, err := s.client.PresignedGetObject(ctx, s.cfg.BucketName, objectKey, s.cfg.PresignedGetExpiry, url.Values{})
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
