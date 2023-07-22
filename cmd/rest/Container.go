package rest

import (
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/server"
	"github.com/dmitry-kostin/go-rest/src/services/ping"
	"github.com/dmitry-kostin/go-rest/src/services/user"
	"github.com/dmitry-kostin/go-rest/src/services/user/models"
	"github.com/dmitry-kostin/go-rest/src/services/user/models/adapters"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	config *pkg.Config
	logger *pkg.Logger
	pgpool *pgxpool.Pool

	services struct {
		pingService *ping.Service
		userService *user.Service
	}

	repositories struct {
		userRepository models.UserRepository
	}

	server struct {
		REST   *server.Server
		router *server.Router
	}
}

func InitContainer(config *pkg.Config, logger *pkg.Logger, pgpool *pgxpool.Pool) *Container {
	container := &Container{config: config, logger: logger, pgpool: pgpool}
	container.init()
	return container
}

func (s *Container) init() {
	_ = s.GetPingService()
	_ = s.GetUserService()
	_ = s.GetRestServer()
}

func (s *Container) GetPingService() *ping.Service {
	if s.services.pingService == nil {
		s.services.pingService = ping.NewService(s.pgpool, s.config)
	}
	return s.services.pingService
}

func (s *Container) GetUserRepository() models.UserRepository {
	if s.repositories.userRepository == nil {
		s.repositories.userRepository = adapters.NewPgxRepository(s.pgpool)
	}
	return s.repositories.userRepository
}

func (s *Container) GetUserService() *user.Service {
	if s.services.userService == nil {
		s.services.userService = user.NewService(s.GetUserRepository(), s.config)
	}
	return s.services.userService
}

func (s *Container) GetRestServer() *server.Server {
	if s.server.REST == nil {
		s.server.REST = server.NewServer(s.config, s.logger, s.GetRestServerRouter())
	}
	return s.server.REST
}

func (s *Container) GetRestServerRouter() *server.Router {
	if s.server.router == nil {
		s.server.router = server.NewRouter(s.logger, s.config)
	}
	return s.server.router
}
