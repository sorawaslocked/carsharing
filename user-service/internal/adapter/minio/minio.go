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
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage struct {
	log           *slog.Logger
	client        *minio.Client
	presignClient *minio.Client // signs URLs with PublicEndpoint so the host in the signature matches the URL the caller uses
	cfg           pkgminio.Config
}

func NewObjectStorage(log *slog.Logger, client *minio.Client, cfg pkgminio.Config) (*ObjectStorage, error) {
	log = pkglog.WithComponent(log, "ObjectStorage")

	var presignClient *minio.Client
	if cfg.PublicEndpoint != "" {
		c, err := minio.New(cfg.PublicEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
			Secure: cfg.UseSSL,
		})
		if err != nil {
			log.Error("creating presign client for public endpoint", pkglog.Err(err), slog.String("public_endpoint", cfg.PublicEndpoint))
			return nil, err
		}
		presignClient = c
	}

	return &ObjectStorage{
		log:           log,
		client:        client,
		presignClient: presignClient,
		cfg:           cfg,
	}, nil
}

// presigner returns the client to use for generating presigned URLs.
// When PublicEndpoint is configured, this is a separate client whose endpoint
// matches the host clients will actually reach, so the HMAC signature is valid.
func (s *ObjectStorage) presigner() *minio.Client {
	if s.presignClient != nil {
		return s.presignClient
	}
	return s.client
}

func (s *ObjectStorage) GetDocumentImageUploadData(ctx context.Context, imageType string) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetDocumentImageUploadData"), utils.MetadataFromCtx(ctx))

	objectKey := newObjectKey("documents/" + imageType)

	presignedURL, err := s.presigner().PresignedPutObject(ctx, s.cfg.Bucket, objectKey, s.cfg.PresignedPutExpiry)
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

	presignedURL, err := s.presigner().PresignedPutObject(ctx, s.cfg.Bucket, objectKey, s.cfg.PresignedPutExpiry)
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

	presignedURL, err := s.presigner().PresignedGetObject(ctx, s.cfg.Bucket, objectKey, s.cfg.PresignedGetExpiry, url.Values{})
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
