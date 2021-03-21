package user

import "errors"

var (
	NotFoundError          = errors.New("user not found")
	AlreadyAddedError      = errors.New("user is added")
	InvalidUpdateError     = errors.New("invalid user update")
	UnmarshallError = errors.New("user json unmarshalling error")
	UnauthorizedError = errors.New("user is not authorized")
)
