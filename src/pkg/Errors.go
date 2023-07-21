package pkg

import (
	"errors"
)

var (
	ErrBadInput = errors.New("provided input is invalid")
	ErrNotFound = errors.New("resource not found")
)
