package storages

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"io"
	"jungle-test/app/internal/domain/entity"
	"jungle-test/app/internal/domain/services"
	"jungle-test/app/pkg/apperrors"
	"net/url"
	"path/filepath"
)

var _ = services.UploadStorage(ImagesStorage{})

type ImagesStorage struct {
	client     *minio.Client
	bucketName string
}

func NewImagesStorage(client *minio.Client, bucketName string) *ImagesStorage {
	return &ImagesStorage{client: client, bucketName: bucketName}
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
		s.bucketName, name, reader, image.Size, minio.PutObjectOptions{ContentType: image.ContentType},
	)
	if err != nil {
		return nil, apperrors.NewInternal("upload image to minio", err)
	}

	url = s.client.EndpointURL()
	url.Path = filepath.Join(url.Path, s.bucketName, name)
	return url, nil
}
