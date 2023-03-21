package main

import (
	"log"

	"jungle-test/internal/app"
	"jungle-test/pkg/logger"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		logger.Get().Info("shutting down server")
		a.CloseDeps()
		logger.Get().Info("done. exiting")
	}()

	if err = a.ReadConfig(); err != nil {
		logger.Get().Error(err, "read config")
		return
	}

	if err = a.InitTracer(); err != nil {
		logger.Get().Error(err, "init tracer")
		return
	}

	if err = a.SetupDeps(); err != nil {
		logger.Get().Error(err, "setup dependencies")
		return
	}

	a.Start()
}
