package api

import (
	"github.com/gofiber/fiber/v2"
	"jungle-test/app/internal/controller/http/api/jwt"
	"jungle-test/app/internal/controller/http/api/middleware"
	"jungle-test/app/internal/domain/services"
)

func NewServer(userService *services.UserService, imageService *services.ImageService, jc *jwt.JwtConfigurator) *fiber.App {
	app := fiber.New()
	app.Use(middleware.InternalErrorMiddleware())

	imageHandler := NewImageHandler(imageService)
	app.Post("/upload-picture", imageHandler.UploadPhoto)

	userHandler := NewUserHandler(userService, imageService, jc)
	app.Post("/login", userHandler.Login)
	app.Get("/images", userHandler.GetImages)

	return app
}
