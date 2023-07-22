package rest

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/application"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	config    *pkg.Config
	logger    *pkg.Logger
	container *Container
	exitFn    func()
}

func InitService(
	config *pkg.Config,
	logger *pkg.Logger,
	container *Container,
	exitFn func(),
) *Service {
	s := &Service{
		config:    config,
		logger:    logger,
		container: container,
		exitFn:    exitFn,
	}
	s.logger.Info().Msgf("[service] staring service %s in %s mode ...", s.config.Version, s.config.AppEnv)
	s.buildRESTSService()
	return s
}

func (s *Service) buildRESTSService() {
	router := s.container.server.router
	router.AddRoutes("/api", false, []application.Route{
		{Name: "Ping", Method: "GET", Pattern: "/ping", Handler: s.container.services.pingService.Ping},
	})
	//router.AddRoutes("/api", true, []application.Route{
	//	{Name: "CreateUser", Method: "POST", Pattern: "/users", Handler: s.container.services.userService.CreateUser},
	//	{Name: "ListUsers", Method: "GET", Pattern: "/users", Handler: s.container.services.userService.ListUsers},
	//	{Name: "GetUser", Method: "GET", Pattern: "/users/{id}", Handler: s.container.services.userService.GetUser},
	//	{Name: "RemoveUser", Method: "DELETE", Pattern: "/users/{id}", Handler: s.container.services.userService.RemoveUser},
	//})
}

func (s *Service) StartRESTServer() {
	s.logger.Info().Msgf("[service] starting REST server on %s:%s ...", s.config.AppHostName, s.config.AppPort)
	server := s.container.server.REST
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error().Msgf("[service] server failed to start: %s", err)
		s.shutdown()
	}
}

func (s *Service) WaitForStopSignal() {
	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt, syscall.SIGTERM)
	sig := <-stopSignalChannel
	if _, ok := sig.(os.Signal); ok {
		s.logger.Info().Msgf("[service] received stop signal: %s", sig)
		close(stopSignalChannel)
		s.shutdown()
	}
}

func (s *Service) shutdown() {
	s.logger.Info().Msg("[service] stopping services ...")
	if s.exitFn != nil {
		s.logger.Info().Msg("[service] canceling context ...")
		s.exitFn()
	}
	server := s.container.server.REST
	if server != nil {
		s.logger.Info().Msg("[service] stopping server gracefully ...")
		if err := server.Shutdown(context.Background()); err != nil {
			s.logger.Warn().Msgf("shutdown: failed to stop the REST server: %s", err)
		}
	}
	s.logger.Info().Msg("shutdown: all services stopped - Hasta la vista, baby!")
	s.exitFn()
}
