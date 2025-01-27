package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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

func ValidationErr(input interface{})error{
	validate := validator.New()

	var errMsg []string
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
			case "required":
				errMsg = append(errMsg, fmt.Sprintf("Field '%s' is required", err.Field()))
			case "email":
				errMsg = append(errMsg, fmt.Sprintf("Field '%s' must be a valid email address", err.Field()))
			case "min":
				errMsg = append(errMsg, fmt.Sprintf("Field '%s' should be at least %s characters long", err.Field(), err.Param()))
			case "max":
				errMsg = append(errMsg, fmt.Sprintf("Field '%s' can be at most %s characters long", err.Field(), err.Param()))
			case "oneof":
				errMsg = append(errMsg, fmt.Sprintf("Field '%s' must be one of the following values: %s", err.Field(), err.Param()))
			case "phone":
				errMsg = append(errMsg, fmt.Sprintf("Field '%s' must be valid phone number", err.Field()))
			}
		}
		return fmt.Errorf("%s", strings.Join(errMsg, "\n"))
	}
	return nil
}