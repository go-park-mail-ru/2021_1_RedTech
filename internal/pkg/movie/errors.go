package movie

import (
	"errors"
	"net/http"
)

var (
	AlreadyExists   = errors.New("already exists")
	NotFoundError   = errors.New("movie not found")
	UnmarshallError = errors.New("movie json unmarshalling error")
	InvalidFilterError = errors.New("invalid filter")
)

func CodeFromError(e error) (code int) {
	switch e {
	case NotFoundError:
		code = http.StatusNotFound
	case UnmarshallError:
		code = http.StatusBadRequest
	case InvalidFilterError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	return code
}
