package server

import (
	"fmt"
	"strconv"
)

type ErrorResponse struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	ErrorMessage string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return strconv.Itoa(e.Status) + ": " + e.Message + ": " + e.ErrorMessage
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (e *ErrorResponse) GoString() string {
	return fmt.Sprintf(`
{
	Status: %s,
	Message: %s,
	ErrorMessage: %s
}`,
		e.Status,
		e.Message,
		e.ErrorMessage,
	)
}
