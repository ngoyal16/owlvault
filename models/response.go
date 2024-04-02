package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ValidationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numberic", e.Field())
	case "unique":
		return fmt.Sprintf("%s must contains unique values", e.Field())
	}
	return fmt.Sprintf("%s is not valid", e.Field())

}

func FormatErrors(err error) []string {
	var errors []string

	switch err.(type) {
	case *json.UnmarshalTypeError:
		return append(errors, "Data type failed Please provide valid values.")
	case validator.ValidationErrors:
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errors = append(errors, fieldErr.StructField()+": "+ValidationErrorToText(fieldErr))
		}
	}

	return errors
}
