package config

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type config struct {
	AppEnv      string `env:"APP_ENV" envDefault:"LOCAL"`
	Port        string `env:"APP_PORT" envDefault:"3001"`
	DatabaseURI string `env:"DB_URI" envDefault:"postgres://username:password@localhost:5432/database_name"`
	Version     string `envDefault:"-"`
}

var AppConfig = &config{}

func LoadConfig() *config {
	version, err := ParseVersionFile("cmd/api-service/VERSION")
	if err != nil {
		log.Fatal(err)
	}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.Parse(AppConfig, opts); err != nil {
		log.Fatal(err)
	}
	AppConfig.Version = version
	return AppConfig
}
