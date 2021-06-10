package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func DescriptiveError(verr validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errors[f.Field()] = getMessage(f.Field(), err)
	}

	return errors
}

func getMessage(fieldName string, errorMessage string) string {
	switch errorMessage {
	case "required":
		return fieldName + "is required"
	case "email":
		return "Please provide the valid email id"
	}
	return errorMessage
}
