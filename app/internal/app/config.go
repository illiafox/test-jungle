package app

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Postgres struct {
	URL string `env:"POSTGRES_URL"`
}

type Telemetry struct {
	JaegerURL   string `env:"JAEGER_URL"          env-required:""`
	ServiceName string `env:"JAEGER_SERVICE_NAME" env-default:"User Service"`
}

type Minio struct {
	Endpoint       string `env:"MINIO_ENDPOINT"          env-required:""`
	SSL            bool   `env:"MINIO_SSL_MODE"          env-required:""`
	BucketName     string `env:"MINIO_BUCKET_NAME" env-required:""`
	BucketLocation string `env:"MINIO_BUCKET_LOCATION" env-required:""`
}

type Jwt struct {
	PrivateKeyPath      string        `env:"JWT_PUBLIC_KEY_PATH"          env-required:""`
	AccessTokenDuration time.Duration `env:"JWT_ACCESS_TOKEN_DURATION" env-required:""`
}

type Server struct {
	MainPort    int `env:"HTTP_PORT" env-default:"8080"`
	MetricsPort int `env:"HTTP_METRICS_PORT" env-default:"8082"`
}

type Config struct {
	Minio     Minio
	Postgres  Postgres
	Telemetry Telemetry
	Server    Server
	Jwt       Jwt
}

func (app *App) ReadConfig() error {
	return cleanenv.ReadEnv(&app.cfg)
}
