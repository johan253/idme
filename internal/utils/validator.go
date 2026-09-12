package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasNumber    = regexp.MustCompile(`[0-9]`)
	hasSymbol    = regexp.MustCompile(`[!@#~$%^&*()+|_]{1}`)
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("password_secure", validatePasswordSecure)
	return validate
}

func validatePasswordSecure(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return hasUppercase.MatchString(password) &&
		hasNumber.MatchString(password) &&
		hasSymbol.MatchString(password)
}
