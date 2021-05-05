package actor

import (
	"errors"
	"net/http"
)

var (
	NotFoundError = errors.New("actor not found")
)

func CodeFromError(e error) (code int) {
	switch e {
	case NotFoundError:
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}
	return code
}
