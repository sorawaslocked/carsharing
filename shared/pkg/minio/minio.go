package minio

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string `yaml:"endpoint"         env:"MINIO_ENDPOINT"          env-required:"true"`
	AccessKeyID     string `yaml:"access_key_id"    env:"MINIO_ACCESS_KEY_ID"     env-required:"true"`
	SecretAccessKey string `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY" env-required:"true"`
	Bucket          string `yaml:"bucket"           env:"MINIO_BUCKET"            env-required:"true"`
	UseSSL          bool   `yaml:"use_ssl"          env:"MINIO_USE_SSL" env-default:"false"`
}

func NewClient(log *slog.Logger, cfg Config) (*minio.Client, error) {
	log = pkglog.WithMethod(log, "minio.NewClient")

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	}

	client, err := minio.New(cfg.Endpoint, opts)
	if err != nil {
		log.Error("connecting to minio", pkglog.Err(err), slog.String("endpoint", cfg.Endpoint))

		return nil, ErrFailedConnection
	}

	return client, nil
}

func Ping(ctx context.Context, log *slog.Logger, client *minio.Client) error {
	log = pkglog.WithMethod(log, "minio.Ping")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	_, err := client.ListBuckets(ctx)
	if err != nil {
		log.Error("pinging minio", pkglog.Err(err), slog.String("endpoint", client.EndpointURL().String()))

		return ErrFailedConnection
	}

	return nil
}
