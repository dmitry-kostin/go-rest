package rest

import (
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/server"
	"github.com/dmitry-kostin/go-rest/src/services/ping"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Container struct {
	config *pkg.Config
	logger *pkg.Logger
	pgpool *pgxpool.Pool

	services struct {
		pingService *ping.Service
		restServer  *server.Server
	}
}

func InitContainer(config *pkg.Config, logger *pkg.Logger, pgpool *pgxpool.Pool) *Container {
	container := &Container{config: config, logger: logger, pgpool: pgpool}
	container.init()
	return container
}

func (s *Container) init() {
	_ = s.GetPingService()
	_ = s.GetRestServer()
}

func (s *Container) GetPingService() *ping.Service {
	if s.services.pingService == nil {
		s.logger.Info().Msg("configuring ping service ...")
		s.services.pingService = ping.NewService(s.pgpool, s.config)
	}
	return s.services.pingService
}

func (s *Container) GetRestServer() *server.Server {
	if s.services.restServer == nil {
		s.services.restServer = server.NewServer(s.config, s.logger)
	}
	return s.services.restServer
}
