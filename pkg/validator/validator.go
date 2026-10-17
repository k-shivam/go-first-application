package validator

import (
	"fmt"
	"go-first-application/pkg/response"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) ([]response.FieldError, bool) {
	err := validate.Struct(s)
	if err == nil {
		return nil, true
	}

	var fieldErrors []response.FieldError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldErrors = append(fieldErrors, response.FieldError{
				Field:   e.Field(),
				Message: formatErrorMessage(e),
			})
		}
	}
	return fieldErrors, false
}

func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	default:
		return "invalid value"
	}
}
