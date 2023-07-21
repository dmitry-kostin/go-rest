package server

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dmitry-kostin/go-rest/src/application"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/server/middleware"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"net/http"
)

type Server struct {
	negroni    *negroni.Negroni
	router     *mux.Router
	render     *render.Render
	logger     *pkg.Logger
	config     *pkg.Config
	httpServer *http.Server
}

func NewServer(config *pkg.Config, logger *pkg.Logger) *Server {
	n := negroni.New()
	router := mux.NewRouter().StrictSlash(true)
	renderer := render.New()
	n.Use(middleware.NewLogger(logger))
	n.Use(middleware.NewSecure(config))
	n.UseHandler(router)
	govalidator.SetFieldsRequiredByDefault(true)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.AppHostName, config.AppPort),
		Handler: n,
	}
	return &Server{n, router, renderer, logger, config, httpServer}
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

func (s *Server) AddRoute(route *application.Route) {
	s.router.PathPrefix("/api").
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		HandlerFunc(s.makeHandlerFunc(route.Handler))
}

func (s *Server) makeHandlerFunc(handlerFunc application.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		statusCode := rw.(negroni.ResponseWriter).Status()
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		rw.Header().Set("X-Powered-By", "PHP/5.4.0")
		v, err := handlerFunc(rw, r)
		if err != nil {
			s.onHandlerError(rw, r, err)
			return
		}
		_ = s.render.JSON(rw, statusCode, v)
	})
}

func (s *Server) onHandlerError(rw http.ResponseWriter, rq *http.Request, err error) {
	if err == nil {
		return
	}
	s.logger.Error().Err(err).Msg("[application error]")
	statusCode := rw.(negroni.ResponseWriter).Status()
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	var govalidatorErr govalidator.Errors
	errorResponse := ErrorResponse{statusCode, "Unknown application error", err.Error()}
	if errors.Cause(err) == pkg.ErrNotFound {
		errorResponse.Status = http.StatusNotFound
		errorResponse.Message = errors.Cause(err).Error()
		_ = s.render.JSON(rw, errorResponse.Status, errorResponse)
		return
	}
	if errors.Cause(err) == pkg.ErrBadInput {
		errorResponse.Status = http.StatusBadRequest
		errorResponse.Message = errors.Cause(err).Error()
		_ = s.render.JSON(rw, errorResponse.Status, errorResponse)
		return
	}
	if errors.As(err, &govalidatorErr) {
		errs := err.(govalidator.Errors).Errors()
		for _, e := range errs {
			s.logger.Error().Err(e).Msg("[govalidator error]")
		}
		errorResponse.Status = http.StatusBadRequest
		errorResponse.Message = errors.Cause(err).Error()
		_ = s.render.JSON(rw, errorResponse.Status, errorResponse)
		return
	}
	_ = s.render.JSON(rw, errorResponse.Status, errorResponse)
}
