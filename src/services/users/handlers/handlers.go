package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/dmitry-kostin/go-rest/src/application"
	"github.com/dmitry-kostin/go-rest/src/server"
	"github.com/dmitry-kostin/go-rest/src/services/users/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type CreateUserReq struct {
	IdentityId uuid.UUID `json:"identity_id" valid:"-"`
	Email      string    `json:"email" valid:"email"`
	FirstName  string    `json:"first_name" valid:"-"`
	LastName   string    `json:"last_name" valid:"-"`
}

func CreateUserHandler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	decoder := json.NewDecoder(r.Body)
	var data CreateUserReq
	err := decoder.Decode(&data)
	if err != nil {
		return &server.ErrorResponse{
			Message:      "Unable to decode the request body",
			ErrorMessage: err.Error(),
			Status:       http.StatusBadRequest,
		}
	}
	_, err = govalidator.ValidateStruct(&data)
	if err != nil {
		return err
	}
	user := &models.User{
		Id:         uuid.New(),
		IdentityId: data.IdentityId,
		Email:      data.Email,
		Role:       models.Customer,
		FirstName:  data.FirstName,
		LastName:   data.LastName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = app.UsersRepository.CreateUser(user)
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusCreated, &struct {
		Id uuid.UUID `json:"id"`
	}{user.Id})
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
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		return err
	}
	user, err := app.UsersRepository.GetUser(id)
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusOK, user)
}

func RemoveUserHandler(rw http.ResponseWriter, r *http.Request, app *application.Application) error {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		return err
	}
	err = app.UsersRepository.RemoveUser(id)
	if err != nil {
		return err
	}
	return app.Render.JSON(rw, http.StatusNoContent, nil)
}
