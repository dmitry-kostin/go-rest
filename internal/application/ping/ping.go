package ping

import (
	"github.com/dmitry-kostin/go-rest/internal/application"
	"net/http"
)

func Handler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	response := Ping{
		Pong:    "You reached the destination. Pong.",
		Version: app.Version,
		Env:     app.Env,
	}
	return app.Render.JSON(rw, http.StatusBadRequest, response)
}
