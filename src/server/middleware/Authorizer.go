package middleware

import (
	"github.com/unrolled/render"
	"net/http"
)

type Authorizer struct {
	renderer *render.Render
}

func NewAuthorizer(renderer *render.Render) *Authorizer {
	return &Authorizer{renderer}
}

func (s *Authorizer) ServeHTTP(rw http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
	if !s.isAuthorized(rq) {
		_ = s.renderer.JSON(rw, http.StatusUnauthorized, struct {
			Message string `json:"message"`
		}{"Permission denied"})
		return
	}
	next(rw, rq)
}

func (s *Authorizer) isAuthorized(rq *http.Request) bool {
	return false
}
