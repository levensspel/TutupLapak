package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IsCategoryProduct(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		if value == "Food" {
			return true
		}
		if value == "Beverage" {
			return true
		}
		if value == "Clothes" {
			return true
		}
		if value == "Furniture" {
			return true
		}
		if value == "Tools" {
			return true
		}
	}
	return false
}

func IsSearchCategoryProduct(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	rgx := regexp.MustCompile(`^sold-\d+$`)
	if ok {
		if value == "newest" {
			return true
		}
		if value == "cheapest" {
			return true
		}
		if rgx.MatchString(value) {
			return true
		}
	}
	return false
}
