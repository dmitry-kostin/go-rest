package ping

import (
	"context"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/services/ping/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type Service struct {
	pg     *pgxpool.Pool
	config *pkg.Config
}

func NewService(pg *pgxpool.Pool, config *pkg.Config) *Service {
	return &Service{pg, config}
}

func (s Service) Ping(http.ResponseWriter, *http.Request) (interface{}, error) {
	err := s.pg.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	logger := pkg.NewLogger()
	logger.Info().Msg("Ping")
	return &models.Ping{
		Pong:    "You reached the destination. Pong.",
		Version: s.config.Version,
		Env:     s.config.AppEnv,
	}, nil
}
