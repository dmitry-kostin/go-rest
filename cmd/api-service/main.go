package main

import (
	"github.com/dmitry-kostin/go-rest/internal/application"
	"github.com/dmitry-kostin/go-rest/internal/application/server"
	"github.com/dmitry-kostin/go-rest/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	initLogger()
	config.LoadConfig()
	app := application.NewApplication()
	server.StartServer(app)
}

func initLogger() {
	if config.AppConfig.AppEnv == "LOCAL" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}
}
