package movie

import (
	"errors"
	"net/http"
)

var (
	NotFoundError   = errors.New("user not found")
	UnmarshallError = errors.New("user json unmarshalling error")
)

func CodeFromError(e error) (code int) {
	switch e {
	case NotFoundError:
		code = http.StatusNotFound
	case UnmarshallError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	return code
}
