package main

import (
	"jungle-test/app/internal/app"
	"jungle-test/app/pkg/logger"
	"log"
)

func main() {

	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer a.CloseDeps()

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
