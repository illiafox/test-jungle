package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"io"
	"jungle-test/app/internal/controller/http/api"
	"jungle-test/app/internal/controller/http/api/jwt"
	"jungle-test/app/internal/controller/http/monitoring"
	"jungle-test/app/internal/domain/services"
	"jungle-test/app/internal/storages"
	"jungle-test/app/pkg/clients"
	"jungle-test/app/pkg/logger"
)

type Deps struct {
	postgresPool    *pgxpool.Pool
	minioClient     *minio.Client
	jwtConfigurator *jwt.JwtConfigurator
	closeTracer     io.Closer
	zapLogger       *zap.Logger
	//
	imageService *services.ImageService
	userService  *services.UserService
	//
	mainServer    *fiber.App
	metricsServer *fiber.App
}

func (deps *Deps) Setup(config Config) (err error) {

	// Connections
	deps.postgresPool, err = clients.NewPostgresClient(config.Postgres.URL)
	if err != nil {
		return fmt.Errorf("create postgres client: %w")
	}
	deps.minioClient, err = clients.NewMinioClient(
		config.Minio.Endpoint,
		config.Minio.SSL,
		config.Minio.BucketName,
		config.Minio.BucketLocation,
	)
	if err != nil {
		return fmt.Errorf("create minio client: %w")
	}

	// JWT
	deps.jwtConfigurator, err = jwt.NewJwtConfigurator(config.Jwt.AccessTokenDuration, config.Jwt.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("create jwt configurator: %w")
	}

	// Storages, services etc.
	userStorage := storages.NewUserStorage(deps.postgresPool)
	imageListStorage := storages.NewImageListStorage(deps.postgresPool)
	uploadStorage := storages.NewImagesStorage(deps.minioClient, config.Minio.BucketName)

	deps.userService = services.NewUserService(userStorage)
	deps.imageService = services.NewImageService(uploadStorage, imageListStorage)

	// Servers
	deps.mainServer = api.NewServer(deps.userService, deps.imageService, deps.jwtConfigurator)
	deps.metricsServer = monitoring.NewServer(deps.postgresPool, deps.minioClient)

	return nil
}

func (deps *Deps) Close() {

	if deps.mainServer != nil {
		if err := deps.mainServer.Shutdown(); err != nil {
			logger.Get().Error(err, "shutdown main server")
		}
	}

	if deps.metricsServer != nil {
		if err := deps.metricsServer.Shutdown(); err != nil {
			logger.Get().Error(err, "shutdown metrics server")
		}
	}

	if deps.postgresPool != nil {
		deps.postgresPool.Close()
	}

	// deps.minioClient does not need .Close (v7 version)

	if deps.closeTracer != nil {
		if err := deps.closeTracer.Close(); err != nil {
			logger.Get().Error(err, "close tracer")
		}
	}

	if deps.zapLogger != nil {
		if err := deps.zapLogger.Sync(); err != nil {
			logger.Get().Error(err, "sync zap logger")
		}
	}
}
