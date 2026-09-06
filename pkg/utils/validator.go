package utils

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string

		for _, e := range validationErrors {
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("Field '%s' is required", e.Field())
			case "email":
				msg = fmt.Sprintf("Field '%s' must be a valid email address", e.Field())
			case "min":
				msg = fmt.Sprintf("Field '%s' must be at least %s characters long", e.Field(), e.Param())
			case "oneof":
				msg = fmt.Sprintf("Field '%s' must be one of the following: [%s]", e.Field(), e.Param())
			default:
				msg = fmt.Sprintf("Field '%s' failed on validation rule '%s'", e.Field(), e.Tag())
			}
			errorMessages = append(errorMessages, msg)
		}

		// Joins all validation messages separated by a comma and space
		return strings.Join(errorMessages, ", ")
	}

	return err.Error()
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}
