package server

import (
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/application"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/server/middleware"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
)

type Router struct {
	muxRouter            *mux.Router
	renderer             *render.Render
	logger               *pkg.Logger
	authorizerMiddleware *middleware.Authorizer
}

func NewRouter(logger *pkg.Logger, config *pkg.Config) *Router {
	muxRouter := mux.NewRouter().StrictSlash(true)
	muxRouter.Use(middleware.NewLogger(logger).Middleware)
	muxRouter.Use(middleware.NewSecure(config).Middleware)
	renderer := render.New()
	return &Router{
		muxRouter:            muxRouter,
		renderer:             renderer,
		logger:               logger,
		authorizerMiddleware: middleware.NewAuthorizer(renderer, config, logger),
	}
}

func (s *Router) AddRoutes(prefix string, protected bool, route []application.Route) {
	gr := s.muxRouter.PathPrefix(prefix).Subrouter().StrictSlash(true)
	if protected {
		gr.Use(s.authorizerMiddleware.Middleware)
	}
	for _, r := range route {
		gr.Methods(r.Method).Path(r.Pattern).Name(r.Name).
			HandlerFunc(s.makeHandlerFunc(r.Handler))
	}
}

func (s *Router) makeHandlerFunc(handlerFunc application.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		statusCode := rw.(*middleware.StatusCodeResponseWriter).Status()
		v, err := handlerFunc(rw, r)
		if err != nil {
			s.onHandlerError(rw, r, err)
			return
		}
		_ = s.renderer.JSON(rw, statusCode, v)
	}
}

func (s *Router) onHandlerError(rw http.ResponseWriter, rq *http.Request, err error) {
	if err == nil {
		return
	}
	s.logger.Error().Err(err).Send()
	errorResponse := ErrorResponse{
		Status:  rw.(*middleware.StatusCodeResponseWriter).Status(),
		Message: "Something went wrong. Please try again later ...",
		Errors:  []string{errors.Cause(err).Error()},
	}
	var govalidatorErrs govalidator.Errors
	if errors.As(err, &govalidatorErrs) {
		errorResponse.Status = http.StatusBadRequest
		errorResponse.Errors = func() []string {
			var errs []string
			for _, err := range govalidatorErrs.Errors() {
				errs = append(errs, err.Error())
			}
			return errs
		}()
	}
	if errors.Is(err, pkg.ErrNotFound) {
		errorResponse.Status = http.StatusNotFound
		errorResponse.Message = "Resource not found"
	}
	if errors.Is(err, pkg.ErrDuplicate) {
		errorResponse.Status = http.StatusConflict
		errorResponse.Message = "Duplicate conflict"
	}
	if errors.Is(err, pkg.ErrBadInput) {
		errorResponse.Status = http.StatusBadRequest
		errorResponse.Message = "Bad request"
	}
	errorDetails := errors.FlattenDetails(err)
	if len(errorDetails) > 0 {
		errorResponse.Message = errorDetails
	}
	if errorResponse.Status == http.StatusOK {
		errorResponse.Status = http.StatusInternalServerError
	}
	_ = s.renderer.JSON(rw, errorResponse.Status, errorResponse)
}
