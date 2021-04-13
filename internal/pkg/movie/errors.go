package movie

import (
	"errors"
	"net/http"
)

var (
	AlreadyExists      = errors.New("already exists")
	NotFoundError      = errors.New("movie not found")
	UnmarshallError    = errors.New("movie json unmarshalling error")
	InvalidFilterError = errors.New("invalid filter")
	InvalidVoteError   = errors.New("invalid vote")
	RatingUpdateError  = errors.New("can't update rating")
	ViewUpdateError    = errors.New("can't update rating")
	BadParamsError     = errors.New("invalid vote")
	InvalidViewCheck   = errors.New("invalid view check")
	InvalidViewAdd     = errors.New("invalid view add")
)

func CodeFromError(e error) (code int) {
	switch e {
	case NotFoundError:
		code = http.StatusNotFound
	case UnmarshallError:
		code = http.StatusBadRequest
	case InvalidFilterError:
		code = http.StatusBadRequest
	case InvalidVoteError:
		code = http.StatusInternalServerError
	case InvalidViewCheck:
		code = http.StatusInternalServerError
	case InvalidViewAdd:
		code = http.StatusInternalServerError
	case RatingUpdateError:
		code = http.StatusInternalServerError
	case ViewUpdateError:
		code = http.StatusInternalServerError
	default:
		code = http.StatusInternalServerError
	}
	return code
}
