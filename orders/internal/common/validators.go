package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func IsInMask(field string, mask []string) bool {
	for _, f := range mask {
		if f == field {
			return true
		}
	}
	return false
}

func MarshalToBytes(input interface{}) (driver.Value, error) {
	return json.Marshal(input)
}

func NewValidator() *validator.Validate {
	// Create a new validator for a Transaction model.
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return false
		}
		return true
	})

	return validate
}

// ValidatorErrors shows validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := make(map[string]string)

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = generateCustomErrorMessage(err)
	}

	return fields
}

func generateCustomErrorMessage(err validator.FieldError) string {
	fieldName := err.Field()
	tag := err.Tag()

	// You can add custom error messages based on the validation tag
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required.", fieldName)
	case "min":
		minLength := err.Param()
		return fmt.Sprintf("%s must be at least %s characters long.", fieldName, minLength)
	case "max":
		maxLength := err.Param()
		return fmt.Sprintf("%s must not exceed %s characters.", fieldName, maxLength)
	case "email":
		return fmt.Sprintf("Invalid %s input. Please enter a valid email address.", fieldName)
	case "gt":
		param := err.Param()
		return fmt.Sprintf("%s must be greater than %s.", fieldName, param)
	case "e164":
		return fmt.Sprintf("Invalid %s input. Please enter a valid phone number.", fieldName)
	case "http_url":
		return fmt.Sprintf("Invalid %s input. Please enter a valid http url.", fieldName)
	default:
		return fmt.Sprintf("Invalid %s input. Please enter a valid  '%s'.", fieldName, tag)
	}
}
