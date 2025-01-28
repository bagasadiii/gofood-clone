package utils

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound         = errors.New("no data found")
	ErrDatabase         = errors.New("database error")
	ErrUniqueConstraint = errors.New("username or email already exists")
	ErrInternal         = errors.New("internal error")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUnexpected       = errors.New("unexpected err")
	ErrBadRequest       = errors.New("bad request")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrValidation       = errors.New("validation error")
)

func ErrCheck(err error) (int, error) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, err
	case errors.Is(err, ErrDatabase):
		return http.StatusInternalServerError, err
	case errors.Is(err, ErrUniqueConstraint):
		return http.StatusBadRequest, err
	case errors.Is(err, ErrInvalidPassword):
		return http.StatusBadRequest, err
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError, err
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, err
	default:
		return http.StatusInternalServerError, errors.New("unexpected error: " + err.Error())
	}
}