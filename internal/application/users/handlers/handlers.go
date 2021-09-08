package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/dmitry-kostin/go-rest/internal/application"
	"github.com/dmitry-kostin/go-rest/internal/application/users/models"
	"github.com/dmitry-kostin/go-rest/pkg/status"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateUserHandler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	decoder := json.NewDecoder(r.Body)
	var data models.CreateUserReq
	err := decoder.Decode(&data)
	if err != nil {
		return &status.ErrorResponse{
			Message:      "Unable to decode the request body",
			ErrorMessage: err.Error(),
			Status:       http.StatusBadRequest,
		}
	}
	data.Role = models.Customer
	_, err = govalidator.ValidateStruct(&data)
	if err != nil {
		return err
	}
	user, err := app.UsersRepository.CreateUser(&data)
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusCreated, user)
}

func ListUsersHandler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	users, err := app.UsersRepository.ListUsers()
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusOK, users)
}

func GetUserHandler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	user, err := app.UsersRepository.GetUser(id)
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusOK, user)
}

func RemoveUserHandler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	err := app.UsersRepository.RemoveUser(id)
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusNoContent, nil)
}
