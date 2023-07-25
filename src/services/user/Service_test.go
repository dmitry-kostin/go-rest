package user_test

import (
	"bytes"
	"context"
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/db"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/services/user"
	"github.com/dmitry-kostin/go-rest/src/services/user/models"
	"github.com/dmitry-kostin/go-rest/src/services/user/models/adapters"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type DI struct {
	repoWithSuccess models.UserRepository
	repoWithFail    models.UserRepository
	createUserInDB  func(email string) *models.User
	cleanup         func()
}

type FailedUsersRepository struct {
	pool *pgxpool.Pool
}

func (f FailedUsersRepository) CreateUser(*models.User) error {
	return errors.Mark(errors.New("failed"), pkg.ErrDatabaseError)
}

func (f FailedUsersRepository) ListUsers() ([]*models.User, error) {
	return nil, errors.Mark(errors.New("failed"), pkg.ErrDatabaseError)
}

func (f FailedUsersRepository) GetUser(models.EntityId) (*models.User, error) {
	return nil, errors.Mark(errors.New("failed"), pkg.ErrDatabaseError)
}

func (f FailedUsersRepository) RemoveUser(models.EntityId) error {
	return errors.Mark(errors.New("failed"), pkg.ErrDatabaseError)
}

func NewDI(t *testing.T) *DI {
	logger := pkg.NewEmptyLogger()
	dbConfig := &pkg.Config{
		DatabaseName:     "tests",
		DatabaseUser:     "admin_user",
		DatabasePassword: "admin_password",
		DatabaseHostname: "localhost",
	}
	conn := db.InitPostgresConnectionPool(dbConfig, logger)
	repoWithSuccess := adapters.NewPgxRepository(conn)
	repoWithFail := FailedUsersRepository{conn}
	createUserInDB := func(email string) *models.User {
		u := models.NewUser(uuid.New(), email, "Hello", "World", models.Customer)
		err := repoWithSuccess.CreateUser(u)
		if err != nil {
			t.Fatal(err)
		}
		return u
	}
	cleanup := func() {
		_, err := conn.Exec(context.Background(), "TRUNCATE TABLE users")
		if err != nil {
			t.Fatal(err)
		}
	}
	return &DI{repoWithSuccess, repoWithFail, createUserInDB, cleanup}
}

func TestService_CreateUser(t *testing.T) {
	di := NewDI(t)
	di.createUserInDB("hello1@world.com")
	mockRequest := func(payload string) *http.Request {
		rq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(payload))
		rq.Header.Set("Content-Type", "application/json")
		return rq
	}

	type args struct {
		rw http.ResponseWriter
		rq *http.Request
	}
	tests := []struct {
		name    string
		args    args
		repo    models.UserRepository
		wantErr error
	}{
		{
			name: "Should create user",
			args: args{nil, mockRequest(`{
				"email":    	"user@world.com",
				"first_name":  	"Hello",
				"last_name":   	"World",
				"identity_id": 	"9ba2ee37-7a43-47c8-b23e-e542f06a1ccd"
			}`)},
			repo:    di.repoWithSuccess,
			wantErr: nil,
		},
		{
			name:    "Should throw an error when payload is invalid",
			args:    args{nil, mockRequest("invalid payload")},
			repo:    di.repoWithFail,
			wantErr: pkg.ErrBadInput,
		},
		{
			name: "Should throw an error when email invalid",
			args: args{nil, mockRequest(`{
				"email": "Hello"
			}`)},
			repo:    di.repoWithSuccess,
			wantErr: pkg.ErrBadInput,
		},
		{
			name: "Should throw an error when identity invalid",
			args: args{nil, mockRequest(`{
				"identity_id": "Hello"
			}`)},
			repo:    di.repoWithSuccess,
			wantErr: pkg.ErrBadInput,
		},
		{
			name: "Should throw an error when failed to save",
			args: args{nil, mockRequest(`{
				"email":    	"hello@world.com",
				"first_name":  	"Hello",
				"last_name":   	"World",
				"identity_id": 	"9ba2ee37-7a43-47c8-b23e-e542f06a1ccd"
			}`)},
			repo:    di.repoWithFail,
			wantErr: pkg.ErrDatabaseError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user.NewService(tt.repo)
			got, err := s.CreateUser(tt.args.rw, tt.args.rq)
			if tt.wantErr != nil {
				if errors.Is(err, tt.wantErr) {
					return
				}
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, ok := got.(*user.ResponseDto); !ok {
				t.Errorf("CreateUser() got = %v, want valid dto interface", got)
				return
			}
			if got.(*user.ResponseDto).Id == uuid.Nil {
				t.Errorf("CreateUser() got = %v, want valid dto", got)
				return
			}
		})
		di.cleanup()
	}
}

func TestService_GetUser(t *testing.T) {
	di := NewDI(t)
	user1 := di.createUserInDB("hello1@world.com")
	user2 := di.createUserInDB("hello2@world.com")
	mockRequest := func(id uuid.UUID) *http.Request {
		rq := httptest.NewRequest(http.MethodGet, "/", nil)
		rq = mux.SetURLVars(rq, map[string]string{"id": id.String()})
		rq.Header.Set("Content-Type", "application/json")
		return rq
	}

	type args struct {
		rw http.ResponseWriter
		rq *http.Request
	}
	tests := []struct {
		name    string
		args    args
		repo    models.UserRepository
		want    interface{}
		wantErr error
	}{
		{
			name:    "Should get user",
			args:    args{nil, mockRequest(user1.Id)},
			repo:    di.repoWithSuccess,
			want:    user1,
			wantErr: nil,
		},
		{
			name:    "Should throw an error when id format invalid",
			args:    args{nil, mockRequest(uuid.Nil)},
			repo:    di.repoWithSuccess,
			want:    user1,
			wantErr: pkg.ErrBadInput,
		},
		{
			name:    "Should throw an error when failed to get user",
			args:    args{nil, mockRequest(user2.Id)},
			repo:    di.repoWithFail,
			want:    user2,
			wantErr: pkg.ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user.NewService(tt.repo)
			got, err := s.GetUser(tt.args.rw, tt.args.rq)
			if tt.wantErr != nil {
				if errors.Is(err, tt.wantErr) {
					return
				}
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
	di.cleanup()
}

func TestService_ListUsers(t *testing.T) {
	di := NewDI(t)
	user1 := di.createUserInDB("hello1@world.com")
	user2 := di.createUserInDB("hello2@world.com")

	type args struct {
		rw http.ResponseWriter
		rq *http.Request
	}
	tests := []struct {
		name    string
		args    args
		repo    models.UserRepository
		want    interface{}
		wantErr error
	}{
		{
			name:    "Should get users",
			args:    args{nil, nil},
			repo:    di.repoWithSuccess,
			want:    []*models.User{user1, user2},
			wantErr: nil,
		},
		{
			name:    "Should throw an error when failed to list users",
			args:    args{nil, nil},
			repo:    di.repoWithFail,
			want:    nil,
			wantErr: pkg.ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user.NewService(tt.repo)
			got, err := s.ListUsers(tt.args.rw, tt.args.rq)
			if tt.wantErr != nil {
				if errors.Is(err, tt.wantErr) {
					return
				}
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
	di.cleanup()
}

func TestService_RemoveUser(t *testing.T) {
	di := NewDI(t)
	user1 := di.createUserInDB("hello1@world.com")
	user2 := di.createUserInDB("hello2@world.com")
	mockRequest := func(id uuid.UUID) *http.Request {
		rq := httptest.NewRequest(http.MethodGet, "/", nil)
		rq = mux.SetURLVars(rq, map[string]string{"id": id.String()})
		rq.Header.Set("Content-Type", "application/json")
		return rq
	}

	type args struct {
		rw http.ResponseWriter
		rq *http.Request
	}
	tests := []struct {
		name    string
		args    args
		repo    models.UserRepository
		want    interface{}
		wantErr error
	}{
		{
			name:    "Should remove user",
			args:    args{nil, mockRequest(user1.Id)},
			repo:    di.repoWithSuccess,
			want:    &user.ResponseDto{Id: user1.Id},
			wantErr: nil,
		},
		{
			name:    "Should throw an error when id format invalid",
			args:    args{nil, mockRequest(uuid.Nil)},
			repo:    di.repoWithSuccess,
			want:    nil,
			wantErr: pkg.ErrBadInput,
		},
		{
			name:    "Should throw an error when failed to remove user",
			args:    args{nil, mockRequest(user2.Id)},
			repo:    di.repoWithFail,
			want:    nil,
			wantErr: pkg.ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user.NewService(tt.repo)
			got, err := s.RemoveUser(tt.args.rw, tt.args.rq)
			if tt.wantErr != nil {
				if errors.Is(err, tt.wantErr) {
					return
				}
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
	di.cleanup()
}
