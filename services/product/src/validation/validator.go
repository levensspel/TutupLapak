package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/samber/do/v2"
)

func NewValidatorInject(i do.Injector) (*validator.Validate, error) {
	validator := validator.New()
	validator.RegisterValidation("category_product", IsCategoryProduct)
	validator.RegisterValidation("category_search", IsSearchCategoryProduct)

	return validator, nil
}
