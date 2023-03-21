package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"jungle-test/app/internal/domain/entity"
	"net/url"
	"time"
)

type UploadStorage interface {
	UploadPhoto(
		ctx context.Context,
		userID uuid.UUID,
		image entity.Image,
		reader io.Reader,
	) (url *url.URL, err error)
}

type ImageListStorage interface {
	AddImage(ctx context.Context, userID uuid.UUID, image entity.Image) error
	GetImages(ctx context.Context, userID uuid.UUID) (images []entity.Image, err error)
}

type ImageService struct {
	uploadStorage    UploadStorage
	imageListStorage ImageListStorage
}

func NewImageService(uploadStorage UploadStorage, imageListStorage ImageListStorage) *ImageService {
	return &ImageService{uploadStorage: uploadStorage, imageListStorage: imageListStorage}
}

func (s ImageService) UploadPhoto(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	contentType string,
	size int64,
	reader io.Reader,
) (url *url.URL, err error) {

	image := entity.Image{
		Name:        name,
		ContentType: contentType,
		URL:         "",
		Size:        size,
		Created:     time.Now(),
	}

	url, err = s.uploadStorage.UploadPhoto(ctx, userID, image, reader)
	if err != nil {
		return nil, fmt.Errorf("upload photo: %w", err)
	}

	err = s.imageListStorage.AddImage(ctx, userID, image)
	if err != nil {
		return nil, fmt.Errorf("add image to list: %w", err)
	}

	return url, nil
}

func (s ImageService) GetImages(ctx context.Context, userID uuid.UUID) ([]entity.Image, error) {
	return s.imageListStorage.GetImages(ctx, userID)
}
