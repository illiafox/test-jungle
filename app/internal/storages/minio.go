package storages

import (
	"context"
	"io"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"jungle-test/internal/domain/entity"
	"jungle-test/internal/domain/services"
	"jungle-test/pkg/apperrors"
)

var _ = services.UploadStorage(ImagesStorage{})

type ImagesStorage struct {
	client     *minio.Client
	bucketName string
	publicHost string
}

func NewImagesStorage(client *minio.Client, bucketName string, publicHost string) *ImagesStorage {
	return &ImagesStorage{client: client, bucketName: bucketName, publicHost: publicHost}
}

func formName(userID uuid.UUID, name string) string {
	return userID.String() + name
}

func (s ImagesStorage) UploadPhoto(
	ctx context.Context,
	userID uuid.UUID,
	image entity.Image,
	reader io.Reader,
) (url *url.URL, err error) {
	name := formName(userID, image.Name)

	_, err = s.client.PutObject(ctx,
		s.bucketName, name, reader, image.Size, minio.PutObjectOptions{
			ContentType: image.ContentType,
		},
	)
	if err != nil {
		return nil, apperrors.NewInternal("upload image to minio", err)
	}

	url = s.client.EndpointURL()
	url.Host = s.publicHost
	url.Path = filepath.Join(url.Path, s.bucketName, name)
	return url, nil
}
