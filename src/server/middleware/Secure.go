package middleware

import (
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/unrolled/secure"
	"net/http"
)

type Secure struct {
	*secure.Secure
}

func NewSecure(config *pkg.Config) *Secure {
	opts := secure.Options{
		IsDevelopment:      config.AppEnv == "LOCAL",
		AllowedHosts:       []string{},
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
	}
	return &Secure{secure.New(opts)}
}

func (s *Secure) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s.HandlerFuncWithNext(w, r, next)
}
