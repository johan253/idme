package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasNumber    = regexp.MustCompile(`[0-9]`)
	hasSymbol    = regexp.MustCompile(`[!@#~$%^&*()+|_]{1}`)
)

// Validator is a process-wide validator instance. Constructing a validator
// registers tag rules and builds a struct cache, so it's meant to be reused.
var Validator = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		if name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]; name != "" && name != "-" {
			return name
		}
		return f.Name
	})
	v.RegisterValidation("password_secure", validatePasswordSecure)
	return v
}

func validatePasswordSecure(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return hasUppercase.MatchString(password) &&
		hasNumber.MatchString(password) &&
		hasSymbol.MatchString(password)
}

// FriendlyValidationError converts a validator error into a human-readable
// message safe to return to API clients.
func FriendlyValidationError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return "Invalid request"
	}
	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, fieldMessage(fe))
	}
	return strings.Join(msgs, "; ")
}

func fieldMessage(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "password_secure":
		return fmt.Sprintf("%s must contain an uppercase letter, a number, and a symbol", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
