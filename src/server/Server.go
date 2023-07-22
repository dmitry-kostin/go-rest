package server

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/server/middleware"
	"github.com/urfave/negroni"
	"net/http"
)

type Server struct {
	negroni    *negroni.Negroni
	logger     *pkg.Logger
	config     *pkg.Config
	httpServer *http.Server
}

func NewServer(config *pkg.Config, logger *pkg.Logger, router *Router) *Server {
	n := negroni.New()
	n.Use(middleware.NewLogger(logger))
	n.Use(middleware.NewSecure(config))
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.AppHostName, config.AppPort),
		Handler: n,
	}
	n.UseHandler(router.muxRouter)
	govalidator.SetFieldsRequiredByDefault(true)
	return &Server{negroni: n, logger: logger, config: config, httpServer: httpServer}
}

func (s *Server) ListenAndServe() error {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
