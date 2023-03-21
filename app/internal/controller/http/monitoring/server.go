package monitoring

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"jungle-test/internal/controller/http/api/middleware"
)

func NewServer(postgresPool *pgxpool.Pool, minioClient *minio.Client) *fiber.App {
	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Use(middleware.InternalErrorMiddleware())

	// prometheus
	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")

	// pprof
	app.Use(pprof.New())

	// healthcheck
	hcHandler := NewHealthCheckHandler(postgresPool, minioClient)
	app.Get("/healtz", hcHandler.HealthCheck)

	return app
}
