package application

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/dmitry-kostin/go-rest/internal/application/users/models"
	"github.com/dmitry-kostin/go-rest/internal/application/users/models/adapters"
	"github.com/dmitry-kostin/go-rest/pkg/config"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

type Application struct {
	Render          *render.Render
	Version         string
	Env             string
	Port            string
	Pool            *pgxpool.Pool
	UsersRepository models.UserRepository
}

func NewApplication() Application {
	pool := getConnectionPool()
	govalidator.SetFieldsRequiredByDefault(true)
	return Application{
		Render:          render.New(),
		Version:         config.AppConfig.Version,
		Env:             config.AppConfig.AppEnv,
		Port:            config.AppConfig.Port,
		Pool:            pool,
		UsersRepository: adapters.NewPgxRepository(pool),
	}
}

func getConnectionPool() *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), config.AppConfig.DatabaseURI)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database not available: %v\n", err)
	}
	log.Infof("===> Database contected")
	return conn
}
