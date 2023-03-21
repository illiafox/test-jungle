package app

import (
	"context"
	"fmt"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"jungle-test/app/pkg/logger"
	"jungle-test/app/pkg/trace"
	"os"
	"os/signal"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	go func() {
		defer stop()

		err := app.deps.mainServer.Listen(fmt.Sprintf(":%d", app.cfg.Server.MainPort))
		if err != nil {
			logger.Get().Error(err, "mainServer.Listen")
		}
	}()

	go func() {
		defer stop()

		err := app.deps.metricsServer.Listen(fmt.Sprintf(":%d", app.cfg.Server.MetricsPort))
		if err != nil {
			logger.Get().Error(err, "metricsServer.Listen")
		}
	}()

	<-ctx.Done()
}
