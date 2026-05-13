package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type Checker struct {
	client     *minio.Client
	bucketName string
}

func NewChecker(client *minio.Client, bucketName string) *Checker {
	return &Checker{client: client, bucketName: bucketName}
}

func (c *Checker) Ping(ctx context.Context) error {
	_, err := c.client.BucketExists(ctx, c.bucketName)
	return err
}
