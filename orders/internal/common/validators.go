package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// To do: handle proper validation
func ValidateGeneric(param any) error {
	if err := validator.New().Struct(param); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		validationErrors := err.(validator.ValidationErrors) // nolint: errorlint
		for _, ve := range validationErrors {
			if ve.Tag() == "oneof" {
				return fmt.Errorf(fmt.Sprintf("%s field should be one of %s", strings.ToLower(ve.Field()), ve.Param()))
			}
			return fmt.Errorf("%s field is required", strings.ToLower(ve.Field())) //nolint: staticcheck

		}
	}

	return nil
}

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
