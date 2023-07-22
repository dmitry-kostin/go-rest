package pkg

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"os"
)

// Version can now be provided on build using -ldflags "-X=back/pkg/config.Versi"
var Version string

type Config struct {
	AppHostName      string `env:"APP_HOST_NAME" envDefault:"127.0.0.1"`
	AppEnv           string `env:"APP_ENV" envDefault:"LOCAL"`
	AppPort          string `env:"APP_PORT" envDefault:"3001"`
	DatabaseName     string `env:"DB_APP_NAME"`
	DatabaseUser     string `env:"DB_APP_USER"`
	DatabasePassword string `env:"DB_APP_PASSWORD"`
	DatabaseHostname string `env:"DB_APP_HOSTNAME" envDefault:"127.0.0.1"`
	Version          string `envDefault:"-"`
}

func NewConfig() *Config {
	cnf := &Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cnf, opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cnf.Version = Version
	return cnf
}
