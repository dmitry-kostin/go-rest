package user

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/services/user/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Service struct {
	repository models.UserRepository
	config     *pkg.Config
}

func NewService(repository models.UserRepository, config *pkg.Config) *Service {
	return &Service{repository, config}
}

type CreateUserRequestBody struct {
	IdentityId uuid.UUID `json:"identity_id" valid:"-"`
	Email      string    `json:"email" valid:"email"`
	FirstName  string    `json:"first_name" valid:"-"`
	LastName   string    `json:"last_name" valid:"-"`
}

type ResponseDto struct {
	Id models.EntityId `json:"id"`
}

func (s *Service) CreateUser(rw http.ResponseWriter, rq *http.Request) (interface{}, error) {
	wrapWith := "[Service.CreateUser]"
	var data CreateUserRequestBody
	err := json.NewDecoder(rq.Body).Decode(&data)
	if err != nil {
		return nil, pkg.AnnotateErrorWithDetail(err, pkg.ErrBadInput, wrapWith, "Provided input is invalid")
	}
	_, err = govalidator.ValidateStruct(&data)
	if err != nil {
		return nil, pkg.AnnotateErrorWithDetail(err, pkg.ErrBadInput, wrapWith, "Input validation failed, please recheck your data")
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
	err = s.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &ResponseDto{user.Id}, nil
}

func (s *Service) ListUsers(rw http.ResponseWriter, rq *http.Request) (interface{}, error) {
	wrapWith := "[Service.ListUsers]"
	users, err := s.repository.ListUsers()
	if err != nil {
		return nil, errors.Wrap(err, wrapWith)
	}
	return users, nil
}

func (s *Service) GetUser(rw http.ResponseWriter, rq *http.Request) (interface{}, error) {
	wrapWith := "[Service.GetUser]"
	id, err := s.getUserIdFromRequest(rq)
	if err != nil {
		return nil, pkg.AnnotateErrorWithDetail(err, pkg.ErrBadInput, wrapWith, "Provided input is invalid")
	}
	user, err := s.repository.GetUser(id)
	if err != nil {
		return nil, errors.Wrap(err, wrapWith)
	}
	return user, nil
}

func (s *Service) RemoveUser(rw http.ResponseWriter, rq *http.Request) (interface{}, error) {
	wrapWith := "[Service.RemoveUser]"
	id, err := s.getUserIdFromRequest(rq)
	if err != nil {
		return nil, pkg.AnnotateErrorWithDetail(err, pkg.ErrBadInput, wrapWith, "Provided input is invalid")
	}
	err = s.repository.RemoveUser(id)
	if err != nil {
		return nil, errors.Wrap(err, wrapWith)
	}
	return &ResponseDto{id}, nil
}

func (s *Service) getUserIdFromRequest(rq *http.Request) (models.EntityId, error) {
	var id uuid.UUID
	id, err := uuid.Parse(mux.Vars(rq)["id"])
	if err != nil {
		return id, err
	}
	if id.Version() != 4 {
		return id, errors.New("invalid uuid version")
	}
	return id, nil
}
