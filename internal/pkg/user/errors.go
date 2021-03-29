package user

import (
	"errors"
	"net/http"
)

var (
	NotFoundError      = errors.New("user not found")
	InvalidCredentials = errors.New("invalid credentials")
	AlreadyAddedError  = errors.New("user is added")
	InvalidUpdateError = errors.New("invalid user update")
	UnmarshallError    = errors.New("user json unmarshalling error")
	UnauthorizedError  = errors.New("user is not authorized")
)

func CodeFromError(e error) (code int) {
	switch e {
	case NotFoundError:
		code = http.StatusNotFound
	case InvalidCredentials:
		code = http.StatusForbidden
	case AlreadyAddedError:
		code = http.StatusForbidden
	case InvalidUpdateError:
		code = http.StatusNotAcceptable
	case UnmarshallError:
		code = http.StatusBadRequest
	case UnauthorizedError:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	return code
}
