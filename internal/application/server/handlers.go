package server

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/dmitry-kostin/go-rest/internal/application"
	"github.com/dmitry-kostin/go-rest/pkg/status"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, *application.Application) error

func MakeHandler(app application.Application, fn HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Terry Pratchett tribute
		rw.Header().Set("X-Powered-By", "PHP/5.4.0")
		// return function with AppEnv
		err := fn(rw, r, &app)
		if err != nil {
			responseStatus := http.StatusInternalServerError
			message := "Internal Server Error. Please try again later"
			errorMessage := err.Error()
			// check for pg errors
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Errorf("PGX error: %v, %v", pgErr.Message, pgErr.Code)
				message = "Internal Server Error. Invalid db operation"
				errorMessage = pgErr.Message
			}
			// check for 404
			if errors.Is(err, pgx.ErrNoRows) {
				responseStatus = http.StatusNotFound
				message = "Requested data not found"
			}
			// check for govalidator errors
			var govalidatorErr govalidator.Errors
			if errors.As(err, &govalidatorErr) {
				errs := err.(govalidator.Errors).Errors()
				for _, e := range errs {
					log.Errorf("Govalidator: %v", e)
				}
				responseStatus = http.StatusBadRequest
				message = "Provided data is invalid"
				errorMessage = err.Error()
			}
			// check for our own errors
			var customErr *status.ErrorResponse
			if errors.As(err, &customErr) {
				log.Errorf("Error: %v", customErr)
				responseStatus = customErr.Status
				message = customErr.Message
				errorMessage = customErr.ErrorMessage
			}
			// default
			log.Errorf("Error: %v", err)
			app.Render.JSON(rw, responseStatus, status.ErrorResponse{
				Status:       responseStatus,
				Message:      message,
				ErrorMessage: errorMessage,
			})
		}
	}
}
