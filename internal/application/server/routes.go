package server

import (
	"github.com/dmitry-kostin/go-rest/internal/application/ping"
	"github.com/dmitry-kostin/go-rest/internal/application/users/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"Ping", "GET", "/ping", ping.Handler},
	//
	Route{"CreateUser", "POST", "/users", handlers.CreateUserHandler},
	Route{"ListUsers", "GET", "/users", handlers.ListUsersHandler},
	Route{"ListUsers", "GET", "/users/{id:[0-9]+}", handlers.GetUserHandler},
	Route{"ListUsers", "DELETE", "/users/{id:[0-9]+}", handlers.RemoveUserHandler},
}
