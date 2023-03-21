package clients

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

// NewMinioClient creates a new client. Credentials must be provided through "MINIO_ROOT_USER" and "MINIO_ROOT_PASSWORD" envs
func NewMinioClient(endpoint string, useSSL bool, bucketName string, bucketLocation string) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewEnvMinio(),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("new minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("check whether bucket exists: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: bucketLocation})
		if err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return minioClient, err
}
