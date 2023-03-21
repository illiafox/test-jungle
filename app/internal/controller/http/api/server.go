package api

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"jungle-test/internal/controller/http/api/jwt"
	"jungle-test/internal/controller/http/api/middleware"
	"jungle-test/internal/domain/services"
)

func NewServer(userService *services.UserService, imageService *services.ImageService, jc *jwt.JwtConfigurator) *fiber.App {
	app := fiber.New(fiber.Config{})
	app.Use(otelfiber.Middleware())
	app.Use(middleware.InternalErrorMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(jc)

	imageHandler := NewImageHandler(imageService)
	app.Post("/upload-picture", authMiddleware.AccessTokenMiddleware, imageHandler.UploadPhoto)

	userHandler := NewUserHandler(userService, imageService, jc)
	app.Post("/login", userHandler.Login)
	app.Get("/images", authMiddleware.AccessTokenMiddleware, userHandler.GetImages)

	return app
}
