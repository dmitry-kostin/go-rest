package application

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, error)
