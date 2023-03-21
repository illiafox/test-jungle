package monitoring

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"net/http"
	"time"
)

type HealthCheckHandler struct {
	pool        *pgxpool.Pool
	minioClient *minio.Client
}

func NewHealthCheckHandler(pool *pgxpool.Pool, minioClient *minio.Client) *HealthCheckHandler {
	return &HealthCheckHandler{pool: pool, minioClient: minioClient}
}

func (h HealthCheckHandler) HealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := h.pool.Ping(ctx)
	if err != nil {
		err = fmt.Errorf("ping pgx pool: %w", err)
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if !h.minioClient.IsOnline() {
		return fiber.NewError(http.StatusInternalServerError, "minio server is offline")
	}

	return nil
}
