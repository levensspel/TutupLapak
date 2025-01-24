package helper

import (
	response "github.com/TimDebug/FitByte/src/model/web"
	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	Validator *validator.Validate
}

func (v XValidator) Validate(data interface{}) []response.ErrorResponse {
	validationErrors := []response.ErrorResponse{}

	errs := v.Validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem response.ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
