package lib

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error       bool        `json:"error"`
	FailedField string      `json:"failed_field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("notEmpty", notEmtpy)
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

func ValidateVar(value interface{}, tags string) error {

	err := validate.Var(value, tags)

	if err != nil {
		return err
	}
	return nil
}

func notEmtpy(fl validator.FieldLevel) bool {
	return fl.Field().Len() != 0
}
