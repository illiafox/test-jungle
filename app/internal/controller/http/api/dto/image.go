package dto

import (
	"jungle-test/app/internal/domain/entity"
	"time"
)

type UploadPhotoResponse struct {
	URL string `json:"url"`
}

type Image struct {
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	Created     time.Time `json:"created_at"`
}

func EntityImagesToDTO(images []entity.Image) []Image {
	conv := make([]Image, len(images))
	for i, m := range images {
		conv[i] = Image(m)
	}
	return conv
}

type GetImagesResponse struct {
	Images []Image
}
