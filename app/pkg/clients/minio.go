package clients

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	return minioClient, err
}
