package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	dto "jungle-test/app/internal/controller/http/api/dto"
	"jungle-test/app/internal/controller/http/api/jwt"
	"jungle-test/app/internal/controller/http/api/middleware"
	"jungle-test/app/internal/domain/services"
	"jungle-test/app/pkg/apperrors"
)

type UserHandler struct {
	userService  *services.UserService
	imageService *services.ImageService
	jc           *jwt.JwtConfigurator
}

func NewUserHandler(userService *services.UserService, imageService *services.ImageService, jc *jwt.JwtConfigurator) *UserHandler {
	return &UserHandler{userService: userService, imageService: imageService, jc: jc}
}

func (h UserHandler) Login(c *fiber.Ctx) error {

	var req dto.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "No image uploaded")
	}

	if req.Username == "" {
		return fiber.NewError(fiber.StatusBadRequest, "'username' field is not provided")
	}
	if req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "'password' field is not provided")
	}

	ctx := c.UserContext()
	user, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		if err == apperrors.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		if err == apperrors.ErrWrongPassword {
			return fiber.NewError(fiber.StatusUnauthorized, "wrong password")
		}
		return fmt.Errorf("login: %w", err)
	}

	accessToken, err := h.jc.GenerateAccessToken(user)
	if err != nil {
		return fmt.Errorf("generate access token: %w", err)
	}

	return c.JSON(dto.LoginUserResponse{
		AccessToken: accessToken,
	})
}

func (h UserHandler) GetImages(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)

	ctx := c.UserContext()
	images, err := h.imageService.GetImages(ctx, claims.UserID)
	if err != nil {
		if err == apperrors.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fmt.Errorf("get images: %w", err)
	}

	return c.JSON(dto.GetImagesResponse{
		Images: dto.EntityImagesToDTO(images),
	})
}
