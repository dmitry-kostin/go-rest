package main

import (
	"github.com/dmitry-kostin/go-rest/cmd/rest"
	"github.com/dmitry-kostin/go-rest/src/db"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"os"
)

func main() {
	stdConfig := pkg.NewConfig()
	stdLogger := initLogger(stdConfig)
	dbConnectionPool := db.InitPostgresConnectionPool(stdConfig, stdLogger)

	container := rest.InitContainer(stdConfig, stdLogger, dbConnectionPool)

	exitFn := func() { os.Exit(1) }
	service := rest.InitService(stdConfig, stdLogger, container, exitFn)

	go service.StartRESTServer()
	service.WaitForStopSignal()
}

func initLogger(appConfig *pkg.Config) *pkg.Logger {
	if appConfig.AppEnv == "LOCAL" {
		return pkg.NewPrettyLogger()
	}
	return pkg.NewLogger()
}
