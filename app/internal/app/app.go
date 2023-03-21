package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"jungle-test/pkg/logger"
	"jungle-test/pkg/trace"
)

type App struct {
	cfg  Config
	deps Deps
}

func New() (*App, error) {
	app := new(App)

	if err := app.SetupLogger(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) SetupDeps() error {
	return a.deps.Setup(a.cfg)
}

func (a *App) CloseDeps() {
	a.deps.Close()
}

func (app *App) SetupLogger() error {
	l, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init zap: %w", err)
	}

	logger.Set(zapr.NewLogger(l))
	app.deps.zapLogger = l

	return nil
}

func (app *App) InitTracer() error {
	c, err := trace.InitTracer(logger.Get(), app.cfg.Telemetry.JaegerURL, app.cfg.Telemetry.ServiceName)
	if err != nil {
		return fmt.Errorf("init tracer: %w", err)
	}

	app.deps.closeTracer = c
	return nil
}

func (app *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		defer stop()

		addr := fmt.Sprintf(":%d", app.cfg.Server.MainPort)
		logger.Get().Info("api server started on", "addr", addr)
		err := app.deps.mainServer.Listen(addr)
		if err != nil {
			logger.Get().Error(err, "mainServer.Listen")
		}
	}()

	go func() {
		defer stop()

		addr := fmt.Sprintf(":%d", app.cfg.Server.MetricsPort)
		logger.Get().Info("metrics server started on", "addr", addr)
		err := app.deps.metricsServer.Listen(addr)
		if err != nil {
			logger.Get().Error(err, "metricsServer.Listen")
		}
	}()

	<-ctx.Done()
}
