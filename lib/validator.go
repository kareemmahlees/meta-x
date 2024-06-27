package lib

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Value       interface{} `json:"value"`
	FailedField string      `json:"failed_field"`
	Tag         string      `json:"tag"`
	Error       bool        `json:"error"`
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true
			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
