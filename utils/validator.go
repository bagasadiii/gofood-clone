package utils

import (
	"fmt"
	"strings"

	"github.com/bagasadiii/gofood-clone/model"
	"github.com/go-playground/validator/v10"
)

var validation = validator.New()

func ValidateUser(user *model.RegisterReq) error {
	err := validation.Struct(user)
	if err != nil {
		var errMsg []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("%v: %s", ErrValidation, strings.Join(errMsg, "\n"))
	}
	return nil
}

func ValidateLogin(data *model.LoginReq) error {
	err := validation.Struct(data)
	if err != nil {
		var errMsg []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("%v: %s", ErrValidation, strings.Join(errMsg, "\n"))
	}
	return nil
}

func ValidateMerchant(merchant *model.Merchant) error {
	err := validation.Struct(merchant)
	if err != nil {
		var errMsg []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("%v: \n%s", ErrValidation, strings.Join(errMsg, "\n"))
	}
	return nil
}

func ValidateDriver(driver *model.Driver) error {
	err := validation.Struct(driver)
	if err != nil {
		var errMsg []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("%v: %s", ErrValidation, strings.Join(errMsg, "\n"))
	}
	return nil
}
