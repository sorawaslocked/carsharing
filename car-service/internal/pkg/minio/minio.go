package minio

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string `yaml:"endpoint"         env:"MINIO_ENDPOINT"          env-required:"true"`
	AccessKeyID     string `yaml:"access_key_id"    env:"MINIO_ACCESS_KEY_ID"     env-required:"true"`
	SecretAccessKey string `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY" env-required:"true"`
	Bucket          string `yaml:"bucket"           env:"MINIO_BUCKET"            env-required:"true"`
	UseSSL          bool   `yaml:"use_ssl"          env:"MINIO_USE_SSL"`
}

func NewClient(cfg Config) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return client, nil
}
