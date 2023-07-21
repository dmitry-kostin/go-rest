package rest

import (
	"context"
	"github.com/dmitry-kostin/go-rest/src/application"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/pkg/errors"
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
	s.container.services.restServer.AddRoute(
		&application.Route{Name: "Ping", Method: "GET", Pattern: "/ping", Handler: s.container.services.pingService.Ping},
	)
}

func (s *Service) StartRESTServer() {
	s.logger.Info().Msgf("[service] starting REST server on %s:%s ...", s.config.AppHostName, s.config.AppPort)
	restServer := s.container.services.restServer
	if err := restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	restServer := s.container.services.restServer
	if restServer != nil {
		s.logger.Info().Msg("[service] stopping server gracefully ...")
		if err := restServer.Shutdown(context.Background()); err != nil {
			s.logger.Warn().Msgf("shutdown: failed to stop the REST server: %s", err)
		}
	}
	s.logger.Info().Msg("shutdown: all services stopped - Hasta la vista, baby!")
	s.exitFn()
}
