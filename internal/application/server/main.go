package server

import (
	"github.com/dmitry-kostin/go-rest/internal/application"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
	"net/http"
)

func StartServer(app application.Application) {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(secureMiddleware(app))
	n.UseHandler(routerHandler(app))
	startupMessage := "===> Starting app (v" + app.Version + ")"
	startupMessage = startupMessage + " on port " + app.Port
	startupMessage = startupMessage + " in " + app.Env + " mode."
	log.Println(startupMessage)
	if app.Env == "LOCAL" {
		n.Run("localhost:" + app.Port)
	} else {
		n.Run(":" + app.Port)
	}
}

func secureMiddleware(app application.Application) negroni.HandlerFunc {
	return secure.New(secure.Options{
		IsDevelopment:      app.Env == "LOCAL",
		AllowedHosts:       []string{},
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
	}).HandlerFuncWithNext
}

func routerHandler(app application.Application) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = MakeHandler(app, route.HandlerFunc)
		router.PathPrefix("/api").
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
