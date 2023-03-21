package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"jungle-test/app/internal/controller/http/api/dto"
	"jungle-test/app/internal/controller/http/api/middleware"
	"jungle-test/app/internal/domain/services"
	"jungle-test/app/pkg/apperrors"
)

type ImageHandler struct {
	imageService *services.ImageService
}

func NewImageHandler(imageService *services.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

func (h ImageHandler) UploadPhoto(c *fiber.Ctx) (err error) {
	claims := middleware.GetUserClaims(c)

	// Parse the multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// Get the "image" field from the form data
	fileHeaders, _ := form.File["image"]
	if len(fileHeaders) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No image uploaded")
	}
	fileHeader := fileHeaders[0]

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, apperrors.NewInternal("close file", file.Close()))
	}()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "no 'Content-Type' header")
	}

	ctx := c.UserContext()
	url, err := h.imageService.UploadPhoto(ctx,
		claims.UserID, fileHeader.Filename, contentType, fileHeader.Size,
		file,
	)
	if err != nil {
		return fmt.Errorf("upload photo: %w", err)
	}

	return c.JSON(dto.UploadPhotoResponse{URL: url.String()})
}
