package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint           string        `yaml:"endpoint" env:"MINIO_ENDPOINT" env-required:"true"`
	AccessKeyID        string        `yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey    string        `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY" env-required:"true"`
	BucketName         string        `yaml:"bucket_name" env:"MINIO_BUCKET_NAME" env-required:"true"`
	UseSSL             bool          `yaml:"use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
	PresignedPutExpiry time.Duration `yaml:"presigned_put_expiry" env:"MINIO_PRESIGNED_PUT_EXPIRY" env-default:"15m"`
	PresignedGetExpiry time.Duration `yaml:"presigned_get_expiry" env:"MINIO_PRESIGNED_GET_EXPIRY" env-default:"1h"`
}

func NewClient(cfg Config) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
}

// EnsureBucket creates the bucket if it does not already exist.
func EnsureBucket(ctx context.Context, client *minio.Client, cfg Config) error {
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return fmt.Errorf("checking bucket existence: %w", err)
	}
	if exists {
		return nil
	}
	if err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}
	return nil
}
