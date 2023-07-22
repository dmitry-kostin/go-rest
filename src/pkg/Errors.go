package pkg

import "github.com/cockroachdb/errors"

var (
	ErrBadInput      = errors.New("provided input is invalid")
	ErrNotFound      = errors.New("resource not found")
	ErrDuplicate     = errors.New("resource is not unique")
	ErrDatabaseError = errors.New("database error")
)

func AnnotateError(original, markAs error, wrapWith string) error {
	return errors.Mark(errors.Wrap(original, wrapWith), markAs)
}

func AnnotateErrorWithDetail(original, markAs error, wrapWith, humanMessage string) error {
	return AnnotateError(errors.WithDetail(original, humanMessage), markAs, wrapWith)
}
